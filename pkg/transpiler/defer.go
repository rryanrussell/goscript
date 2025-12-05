package transpiler

import (
	"go/ast"

	"github.com/rryanrussell/goscript/pkg/js"
)

func deferStmt(ctx *Context, s *ast.DeferStmt) js.Stmt {
	// defer expr -> __gs_deferred.push(() => expr)
	// Note: Arguments are currently captured by reference (closure), not by value.
	// This is a known limitation for Phase 1.

	call := expr(ctx, s.Call)

	return &js.ExprStmt{
		X: &js.Call{
			Func: &js.Selector{
				X:   js.IdentP("__gs_deferred"),
				Sel: js.IdentP("push"),
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

// hasEffects checks if the function body contains defer or recover
func hasEffects(body *ast.BlockStmt) (hasDefer bool, hasRecover bool) {
	ast.Inspect(body, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.DeferStmt:
			hasDefer = true
		case *ast.CallExpr:
			// Check for recover()
			if call, ok := n.(*ast.CallExpr); ok {
				if id, ok := call.Fun.(*ast.Ident); ok && id.Name == "recover" {
					hasRecover = true
				}
			}
		case *ast.FuncLit:
			// Don't inspect nested functions
			return false
		}
		return true
	})
	return
}

// frameBody wraps the function body in a Universal Frame if needed
func frameBody(ctx *Context, body *js.Block, hasDefer, hasRecover bool) *js.Block {
	if !hasDefer && !hasRecover {
		return body
	}

	newBody := &js.Block{}

	// Prologue
	if hasDefer {
		newBody.Lines = append(newBody.Lines, &js.Assign{
			Define: true,
			Lhs:    js.IdentP("__gs_deferred"),
			Rhs:    &js.ArrayLit{},
		})
	}

	if hasRecover {
		newBody.Lines = append(newBody.Lines, &js.Assign{
			Define: true,
			Lhs:    js.IdentP("__gs_panicValue"),
			Rhs:    &js.BasicLit{Value: "null"},
		})
	}

	// Main logic
	tryStmt := &js.TryStmt{
		Body: body,
	}

	if hasRecover {
		tryStmt.Catch = &js.CatchClause{
			Param: js.IdentP("__e"),
			Body: &js.Block{
				Lines: []js.Stmt{
					&js.Assign{
						Lhs: &js.Selector{X: js.IdentP("runtime"), Sel: js.IdentP("__gs_panicValue")},
						Rhs: js.IdentP("__e"),
					},
				},
			},
		}
	}

	if hasDefer {
		tryStmt.Finally = &js.Block{
			Lines: []js.Stmt{
				&js.ExprStmt{
					X: &js.Call{
						Func: &js.Selector{X: js.IdentP("runtime"), Sel: js.IdentP("runDefers")},
						Args: []js.Expr{js.IdentP("__gs_deferred")},
					},
				},
			},
		}
	}

	newBody.Lines = append(newBody.Lines, tryStmt)
	return newBody
}
