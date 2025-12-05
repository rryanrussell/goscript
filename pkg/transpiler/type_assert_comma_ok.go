package transpiler

import (
	"go/ast"

	"github.com/rryanrussell/goscript/pkg/js"
)

func typeAssertCommaOk(ctx *Context, lhs1, lhs2 ast.Expr, rhs *ast.TypeAssertExpr, define bool) *js.Assign {
	// v, ok := x.(T)
	// We need to generate code that returns [v, ok]
	// Since js.Assign expects single Lhs/Rhs, and we are assigning to 2 variables.
	// We can use destructuring assignment in JS if `js` package supports it.
	// `js.Assign` has `Lhs` which is `js.Expr`. `js.ArrayLit` is `js.Expr`.
	// So `[v, ok] = ...` can be represented as `Lhs: &js.ArrayLit{Elts: [v, ok]}`.

	// The RHS should be an expression that evaluates to `[val, true]` or `[zero, false]`.
	// We can use an IIFE.

	targetType := rhs.Type
	tagExpr := expr(ctx, rhs.X)
	check := generateTypeCheck(ctx, tagExpr, targetType)

	// IIFE:
	// (() => {
	//   let _v = x;
	//   if (check(_v)) return [_v, true];
	//   return [null, false]; // or zero value for T
	// })()

	fn := &js.FuncDecl{
		Body: &js.Block{},
		Named: js.Named{Name: js.IdentP("")},
	}
	fn.Body.Inline = true

	// let _v = x
	vName := "_v"
	fn.Body.Lines = append(fn.Body.Lines, &js.Assign{
		Define: true,
		Lhs:    js.IdentP(vName),
		Rhs:    tagExpr,
	})

	vIdent := js.IdentP(vName)
	check = generateTypeCheck(ctx, vIdent, targetType)

	// if (check) return [_v, true]
	fn.Body.Lines = append(fn.Body.Lines, &js.IfStmt{
		Cond: check,
		Body: &js.Block{
			Lines: []js.Stmt{
				&js.Return{
					X: &js.ArrayLit{
						Elts: []js.Expr{
							vIdent,
							&js.BasicLit{Value: "true"},
						},
					},
				},
			},
		},
	})

	// return [null, false]
	// Ideally we should return zero value of T.
	// For now, null/undefined is probably okay for structs/pointers.
	// For primitives, it might be 0, "", false.
	// `TypeExpr` gives string repr.
	// We don't have a `ZeroValue(T)` helper yet.
	// I'll return `null` for now.

	fn.Body.Lines = append(fn.Body.Lines, &js.Return{
		X: &js.ArrayLit{
			Elts: []js.Expr{
				&js.BasicLit{Value: "null"},
				&js.BasicLit{Value: "false"},
			},
		},
	})

	rhsExpr := &js.Call{
		Func: &js.Parens{X: &js.GuestExpr{Node: fn}},
		Args: []js.Expr{},
	}

	// LHS
	lhsExpr1 := expr(ctx, lhs1)
	lhsExpr2 := expr(ctx, lhs2)

	// Verify if lhs1/lhs2 are valid for assignment (idents, selectors, etc)

	return &js.Assign{
		Define: define, // This might be tricky if one is new and other is old, but in Go `:=` implies both new (or at least one new)
		// Actually `v, ok := ...` defines both if `:=`.
		// If js.Assign `Define` is true, it puts `let` before Lhs.
		// `let [v, ok] = ...` works in JS.
		Lhs: &js.ArrayLit{
			Elts: []js.Expr{lhsExpr1, lhsExpr2},
		},
		Rhs: rhsExpr,
	}
}
