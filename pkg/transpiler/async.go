package transpiler

import (
	"go/ast"
	"go/token"

	"github.com/rryanrussell/goscript/pkg/js"
)

func goStmt(ctx *Context, s *ast.GoStmt) js.Stmt {
	// go f() -> runtime.go(() => f())

	// We wrap the call in an arrow function
	call := expr(ctx, s.Call)

	return &js.ExprStmt{
		X: &js.Call{
			Func: &js.Selector{
				X:   js.IdentP("runtime"),
				Sel: js.IdentP("go"),
			},
			Args: []js.Expr{
				&js.ArrowFunc{
					Body: &js.Block{
						Lines: []js.Stmt{
							&js.ExprStmt{X: call},
						},
					},
				},
			},
		},
	}
}

func sendStmt(ctx *Context, s *ast.SendStmt) js.Stmt {
	// ch <- x -> await ch.send(x)

	return &js.ExprStmt{
		X: &js.AwaitExpr{
			X: &js.Call{
				Func: &js.Selector{
					X:   expr(ctx, s.Chan),
					Sel: js.IdentP("send"),
				},
				Args: []js.Expr{expr(ctx, s.Value)},
			},
		},
	}
}

func selectStmt(ctx *Context, s *ast.SelectStmt) js.Stmt {
	// select { ... } -> await runtime.select([...])

	cases := []js.Expr{}

	for _, stmt := range s.Body.List {
		cc, ok := stmt.(*ast.CommClause)
		if !ok {
			continue
		}

		// { op: 'recv'|'send'|'default', chan: ..., val: ..., handler: ... }

		obj := &js.ObjectLit{}

		var handlerBody *js.Block
		handlerBody = block(ctx, &ast.BlockStmt{List: cc.Body})

		if cc.Comm == nil {
			// default case
			obj.Props = append(obj.Props, &js.Prop{Key: js.IdentP("op"), Value: &js.BasicLit{Value: "'default'"}})

			// handler: () => { body }
			obj.Props = append(obj.Props, &js.Prop{
				Key: js.IdentP("handler"),
				Value: &js.ArrowFunc{
					Body: handlerBody,
				},
			})
		} else {
			switch comm := cc.Comm.(type) {
			case *ast.SendStmt:
				// case ch <- v:
				obj.Props = append(obj.Props, &js.Prop{Key: js.IdentP("op"), Value: &js.BasicLit{Value: "'send'"}})
				obj.Props = append(obj.Props, &js.Prop{Key: js.IdentP("chan"), Value: expr(ctx, comm.Chan)})
				obj.Props = append(obj.Props, &js.Prop{Key: js.IdentP("val"), Value: expr(ctx, comm.Value)})

				obj.Props = append(obj.Props, &js.Prop{
					Key: js.IdentP("handler"),
					Value: &js.ArrowFunc{
						Body: handlerBody,
					},
				})

			case *ast.ExprStmt:
				// case <-ch:
				if unary, ok := comm.X.(*ast.UnaryExpr); ok && unary.Op == token.ARROW {
					obj.Props = append(obj.Props, &js.Prop{Key: js.IdentP("op"), Value: &js.BasicLit{Value: "'recv'"}})
					obj.Props = append(obj.Props, &js.Prop{Key: js.IdentP("chan"), Value: expr(ctx, unary.X)})

					// handler: () => { body }
					// But wait, the value is discarded here.
					obj.Props = append(obj.Props, &js.Prop{
						Key: js.IdentP("handler"),
						Value: &js.ArrowFunc{
							Body: handlerBody,
						},
					})
				}

			case *ast.AssignStmt:
				// case x := <-ch:
				// case x, ok := <-ch: (not supported yet in runtime select)

				if len(comm.Rhs) == 1 {
					if unary, ok := comm.Rhs[0].(*ast.UnaryExpr); ok && unary.Op == token.ARROW {
						obj.Props = append(obj.Props, &js.Prop{Key: js.IdentP("op"), Value: &js.BasicLit{Value: "'recv'"}})
						obj.Props = append(obj.Props, &js.Prop{Key: js.IdentP("chan"), Value: expr(ctx, unary.X)})

						// handler: (val) => { x = val; body }
						// We need to inject the assignment into the body

						// Create a new block for the handler
						newBody := &js.Block{}

						// Assign val to x
						// x := val
						// We can use the existing assign logic but we need to map `val` (arg) to `x`

						// comm.Lhs[0] is `x`
						// We create `x = val`

						// If it's `x := ...` (DEFINE), we need `let x = val`

						assign := &js.Assign{
							Define: comm.Tok == token.DEFINE,
							Lhs:    single(ctx, comm.Lhs),
							Rhs:    js.IdentP("_val"),
						}

						newBody.Lines = append(newBody.Lines, assign)
						newBody.Lines = append(newBody.Lines, handlerBody.Lines...)

						obj.Props = append(obj.Props, &js.Prop{
							Key: js.IdentP("handler"),
							Value: &js.ArrowFunc{
								Args: []*js.Var{{Named: js.Named{Name: js.IdentP("_val")}}},
								Body: newBody,
							},
						})
					}
				}
			}
		}

		cases = append(cases, obj)
	}

	return &js.ExprStmt{
		X: &js.AwaitExpr{
			X: &js.Call{
				Func: &js.Selector{
					X:   js.IdentP("runtime"),
					Sel: js.IdentP("select"),
				},
				Args: []js.Expr{
					&js.ArrayLit{Elts: cases},
				},
			},
		},
	}
}

func unaryExpr(ctx *Context, u *ast.UnaryExpr) js.Expr {
	if u.Op == token.ARROW {
		// <-ch -> await ch.recv()
		return &js.AwaitExpr{
			X: &js.Call{
				Func: &js.Selector{
					X:   expr(ctx, u.X),
					Sel: js.IdentP("recv"),
				},
				Args: []js.Expr{},
			},
		}
	}

	// Delegate back to normal unary
	return &js.Unary{
		X:  expr(ctx, u.X),
		Op: op(ctx, u.Op),
	}
}

func hasAsync(body *ast.BlockStmt) bool {
	has := false
	ast.Inspect(body, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.SendStmt:
			has = true
			return false
		case *ast.SelectStmt:
			has = true
			return false
		case *ast.UnaryExpr:
			if x.Op == token.ARROW {
				has = true
				return false
			}
		case *ast.CallExpr:
			if sel, ok := x.Fun.(*ast.SelectorExpr); ok {
				if sel.Sel.Name == "Lock" || sel.Sel.Name == "Wait" || sel.Sel.Name == "Do" {
					has = true
					return false
				}
				if id, ok := sel.X.(*ast.Ident); ok && id.Name == "time" && sel.Sel.Name == "Sleep" {
					has = true
					return false
				}
			}
		case *ast.FuncLit:
			// Don't inspect nested functions
			return false
		}
		return true
	})
	return has
}
