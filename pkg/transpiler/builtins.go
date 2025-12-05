package transpiler

import (
	"go/ast"
	"go/token"

	"github.com/rryanrussell/goscript/pkg/js"
)

func handleLen(ctx *Context, c *ast.CallExpr) js.Expr {
	return &js.Parens{
		X: &js.Binary{
			X: &js.Selector{
				X:             expr(ctx, c.Args[0]),
				Sel:           js.IdentP("length"),
				OptionalChain: true,
			},
			Op: js.HuhHuh,
			Y:  &js.BasicLit{Value: "0"},
		},
	}
}

func handleAppend(ctx *Context, c *ast.CallExpr) js.Expr {
	elts := []js.Expr{
		&js.Spread{
			X: &js.Parens{
				X: &js.Binary{
					X:  expr(ctx, c.Args[0]),
					Op: js.HuhHuh,
					Y:  &js.ArrayLit{},
				},
			},
		},
	}

	for i, arg := range c.Args[1:] {
		e := expr(ctx, arg)
		// Check if this is the last argument and it has an ellipsis (spread)
		if c.Ellipsis != token.NoPos && i == len(c.Args)-2 {
			// Safety check for nil slice: ...(arg ?? [])
			e = &js.Spread{
				X: &js.Parens{
					X: &js.Binary{
						X:  e,
						Op: js.HuhHuh,
						Y:  &js.ArrayLit{},
					},
				},
			}
		}
		elts = append(elts, e)
	}

	return &js.ArrayLit{
		Elts: elts,
	}
}

func handleCopy(ctx *Context, c *ast.CallExpr) js.Expr {
	return &js.Call{
		Func: &js.Selector{
			X:   js.Ident("runtime"),
			Sel: js.IdentP("copy"),
		},
		Args: exprs(ctx, c.Args),
	}
}

func handleMake(ctx *Context, c *ast.CallExpr) js.Expr {
	// check first argument
	arg0 := c.Args[0]
	switch arg0.(type) {
	case *ast.ArrayType:
		// make([]T, len, cap)
		args := []js.Expr{}
		if len(c.Args) > 1 {
			args = append(args, expr(ctx, c.Args[1])) // len
		}
		if len(c.Args) > 2 {
			args = append(args, expr(ctx, c.Args[2])) // cap
		} else if len(c.Args) == 2 {
			// pass undefined/null for cap
		}

		// Add zero value
		// args = append(args, zeroValue(t.Elt))

		return &js.Call{
			Func: &js.Selector{
				X:   js.Ident("runtime"),
				Sel: js.IdentP("makeSlice"),
			},
			Args: args,
		}
	case *ast.MapType:
		// make(map[K]V)
		return &js.Call{
			Func: &js.Selector{
				X:   js.Ident("runtime"),
				Sel: js.IdentP("makeMap"),
			},
			Args: []js.Expr{},
		}
	case *ast.ChanType:
		// make(chan T, [cap])
		args := []js.Expr{}
		if len(c.Args) > 1 {
			args = append(args, expr(ctx, c.Args[1])) // cap
		}
		return &js.Call{
			Func: &js.Selector{
				X:   js.Ident("runtime"),
				Sel: js.IdentP("makeChan"),
			},
			Args: args,
		}
	}
	return unknown(ctx, c)
}

func handlePanic(ctx *Context, c *ast.CallExpr) js.Expr {
	return &js.Call{
		Func: &js.Selector{
			X:   js.Ident("runtime"),
			Sel: js.IdentP("panic"),
		},
		Args: exprs(ctx, c.Args),
	}
}

func handleRecover(ctx *Context, c *ast.CallExpr) js.Expr {
	return &js.Call{
		Func: &js.Selector{
			X:   js.Ident("runtime"),
			Sel: js.IdentP("recover"),
		},
		Args: []js.Expr{},
	}
}

func handleFmt(ctx *Context, c *ast.CallExpr, name string) js.Expr {
	switch name {
	case "Printf":
		return &js.Call{
			Func: &js.Selector{
				X: &js.Selector{
					X:   js.Ident("runtime"),
					Sel: js.IdentP("fmt"),
				},
				Sel: js.IdentP("printf"),
			},
			Args: exprs(ctx, c.Args),
		}
	case "Sprintf":
		return &js.Call{
			Func: &js.Selector{
				X: &js.Selector{
					X:   js.Ident("runtime"),
					Sel: js.IdentP("fmt"),
				},
				Sel: js.IdentP("sprintf"),
			},
			Args: exprs(ctx, c.Args),
		}
	case "Println":
		return &js.Call{
			Func: &js.Selector{
				X: &js.Selector{
					X:   js.Ident("runtime"),
					Sel: js.IdentP("fmt"),
				},
				Sel: js.IdentP("println"),
			},
			Args: exprs(ctx, c.Args),
		}
	}
	return nil
}

func handleReflect(ctx *Context, c *ast.CallExpr, name string) js.Expr {
	return &js.Call{
		Func: &js.Selector{
			X: &js.Selector{
				X:   js.Ident("runtime"),
				Sel: js.IdentP("reflect"),
			},
			Sel: js.IdentP(name),
		},
		Args: exprs(ctx, c.Args),
	}
}

func handleStrings(ctx *Context, c *ast.CallExpr, name string) js.Expr {
	// Map strings.X to runtime.strings.X
	return &js.Call{
		Func: &js.Selector{
			X: &js.Selector{
				X:   js.Ident("runtime"),
				Sel: js.IdentP("strings"),
			},
			Sel: js.IdentP(name),
		},
		Args: exprs(ctx, c.Args),
	}
}

func handleErrors(ctx *Context, c *ast.CallExpr, name string) js.Expr {
	return &js.Call{
		Func: &js.Selector{
			X: &js.Selector{
				X:   js.Ident("runtime"),
				Sel: js.IdentP("errors"),
			},
			Sel: js.IdentP(name),
		},
		Args: exprs(ctx, c.Args),
	}
}

func handleStrconv(ctx *Context, c *ast.CallExpr, name string) js.Expr {
	return &js.Call{
		Func: &js.Selector{
			X: &js.Selector{
				X:   js.Ident("runtime"),
				Sel: js.IdentP("strconv"),
			},
			Sel: js.IdentP(name),
		},
		Args: exprs(ctx, c.Args),
	}
}

func handleMath(ctx *Context, c *ast.CallExpr, name string) js.Expr {
	return &js.Call{
		Func: &js.Selector{
			X: &js.Selector{
				X:   js.Ident("runtime"),
				Sel: js.IdentP("math"),
			},
			Sel: js.IdentP(name),
		},
		Args: exprs(ctx, c.Args),
	}
}

func handleTime(ctx *Context, c *ast.CallExpr, name string) js.Expr {
	call := &js.Call{
		Func: &js.Selector{
			X: &js.Selector{
				X:   js.Ident("runtime"),
				Sel: js.IdentP("time"),
			},
			Sel: js.IdentP(name),
		},
		Args: exprs(ctx, c.Args),
	}

	if name == "Sleep" {
		return &js.AwaitExpr{X: call}
	}
	return call
}

func handleSort(ctx *Context, c *ast.CallExpr, name string) js.Expr {
	return &js.Call{
		Func: &js.Selector{
			X: &js.Selector{
				X:   js.Ident("runtime"),
				Sel: js.IdentP("sort"),
			},
			Sel: js.IdentP(name),
		},
		Args: exprs(ctx, c.Args),
	}
}
