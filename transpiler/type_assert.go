package transpiler

import (
	"fmt"
	"go/ast"

	"github.com/rryanrussell/goscript/js"
)

func typeAssert(ctx *Context, x *ast.TypeAssertExpr) js.Expr {
	// Generate IIFE for type assertion:
	// (() => {
	//   if (check(x)) return x;
	//   throw new Error("panic: interface conversion: interface is ... not T");
	// })()

	// However, IIFE adds overhead and complexity for reading.
	// Can we use a ternary?
	// (check(x) ? x : (()=>{throw "panic"})())

	// We'll use a ternary if possible, or a helper if it gets complex.

	targetType := x.Type
	tagExpr := expr(ctx, x.X)
	check := generateTypeCheck(ctx, tagExpr, targetType)

	// Since we don't have a robust runtime panic yet, and we want to be readable:
	// (x?.__goType == "T" ? x : <panic>)

	// Construct panic expression
	// Since we can't easily put `throw` in an expression in JS (without IIFE or helper),
	// we'll use an IIFE for the failure case or the whole thing.

	// Also, if x is evaluated multiple times (in check and in return), it might be bad if x has side effects.
	// But x is usually a variable or a simple expression in type assertions.
	// If x is a function call `f().(T)`, evaluating `f()` twice is bad.
	// So we should strictly use IIFE or assign to temp var.

	// IIFE approach:
	// (() => { let _v = x; if (check(_v)) return _v; throw ... })()

	// Given "Educational clarity", IIFE is okay-ish.

	// Let's reuse generateTypeCheck from type_switch.go (need to make sure it's accessible or move it)
	// I'll assume it's in the same package `transpiler`.

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

	// Update check to use vIdent instead of tagExpr
	// Re-generate check using vIdent
	check = generateTypeCheck(ctx, vIdent, targetType)

	// if (check) return _v
	fn.Body.Lines = append(fn.Body.Lines, &js.IfStmt{
		Cond: check,
		Body: &js.Block{
			Lines: []js.Stmt{
				&js.Return{X: vIdent},
			},
		},
	})

	// throw error
	typeName := TypeExpr(ctx, targetType)
	fn.Body.Lines = append(fn.Body.Lines, &js.TokenStmt{
		Token: js.Token(fmt.Sprintf("throw new Error(\"interface conversion: interface is not %s\")", typeName)),
	})

	return &js.Call{
		Func: &js.Parens{X: &js.GuestExpr{Node: fn}},
		Args: []js.Expr{},
	}
}
