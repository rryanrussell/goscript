package transpiler

import (
	"fmt"
	"go/ast"

	"github.com/rryanrussell/goscript/js"
)

func typeAssert(ctx *Context, x *ast.TypeAssertExpr) js.Expr {
	targetType := x.Type
	tagExpr := expr(ctx, x.X)

	// IIFE approach:
	// (() => { let _v = x; if (check(_v)) return _v; runtime.panic(...) })()

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
	check := generateTypeCheck(ctx, vIdent, targetType)

	// if (check) return _v
	fn.Body.Lines = append(fn.Body.Lines, &js.IfStmt{
		Cond: check,
		Body: &js.Block{
			Lines: []js.Stmt{
				&js.Return{X: vIdent},
			},
		},
	})

	// runtime.panic(...)
	typeName := TypeExpr(ctx, targetType)
	fn.Body.Lines = append(fn.Body.Lines, js.Panic(fmt.Sprintf("interface conversion: interface is not %s", typeName)))

	return &js.Call{
		Func: &js.Parens{X: &js.GuestExpr{Node: fn}},
		Args: []js.Expr{},
	}
}
