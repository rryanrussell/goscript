package transpiler

import (
	"fmt"
	"go/ast"

	"github.com/rryanrussell/goscript/js"
)

func typeSwitch(ctx *Context, s *ast.TypeSwitchStmt) js.Stmt {
	block := &js.Block{}

	// Handle init statement if present
	if s.Init != nil {
		block.Lines = append(block.Lines, statement(ctx, s.Init))
	}

	// Identify the expression being switched on
	var initialTagExpr js.Expr
	var variableName string
	var isShortDecl bool

	switch assign := s.Assign.(type) {
	case *ast.ExprStmt:
		// x.(type)
		typeAssert := assign.X.(*ast.TypeAssertExpr)
		initialTagExpr = expr(ctx, typeAssert.X)
	case *ast.AssignStmt:
		// v := x.(type)
		typeAssert := assign.Rhs[0].(*ast.TypeAssertExpr)
		initialTagExpr = expr(ctx, typeAssert.X)
		// lhs is the variable
		if ident, ok := assign.Lhs[0].(*ast.Ident); ok {
			variableName = ident.Name
			isShortDecl = true
		}
	}

	// Capture tagExpr into a temporary variable to avoid multiple evaluations
	// and to ensure safety if initialTagExpr has side effects.
	// But if initialTagExpr is already an identifier, we can maybe reuse it?
	// To be safe and consistent with feedback, we always capture it or check if it's ident.
	// If it is an identifier, it's safe to reuse unless it's modified in the loop (unlikely).
	// However, creating a temp var `_typeSwitchTag` is safer.

	// But wait, if we have `switch v := x.(type)`, we have `variableName` ("v").
	// But `v` is only assigned inside the case blocks in Go semantics (with narrow type).
	// So we still need a holding variable for the value of `x`.

	tagVarName := "_tag"
	// Ensure unique name if nested? transpiler doesn't track scopes well yet for uniqueness.
	// But JS `let` is block scoped. So `let _tag` inside the `block` we created is safe.

	block.Lines = append(block.Lines, &js.Assign{
		Define: true,
		Lhs:    js.IdentP(tagVarName),
		Rhs:    initialTagExpr,
	})
	tagExpr := js.IdentP(tagVarName)

	// We need to construct an if-else chain because JS switch doesn't support complex type checks easily
	// (unless we rely solely on __goType matching a string, but primitives need typeof)

	// We'll traverse the cases and build the if-else chain bottom-up or top-down.
	// Actually, top-down is easier for IfStmt structure.

	var firstIf *js.IfStmt
	var currentIf *js.IfStmt

	for _, stmt := range s.Body.List {
		cc, ok := stmt.(*ast.CaseClause)
		if !ok {
			continue
		}

		// Each case clause has a list of types (or nil for default)
		// and a body.

		// Build the condition for this case
		var condition js.Expr

		if cc.List == nil {
			// Default case.
			// It should be the final else, but in Go default can be anywhere.
			// However, for the if-else chain, it must be at the end.
			// We'll handle default separately after the loop maybe?
			// Or we just check if it's default and attach it to the end of the chain.
			continue
		}

		for _, typeNode := range cc.List {
			cond := generateTypeCheck(ctx, tagExpr, typeNode)
			if condition == nil {
				condition = cond
			} else {
				condition = &js.Binary{
					X:  condition,
					Op: js.Or,
					Y:  cond,
				}
			}
		}

		// Body of the case
		caseBody := &js.Block{}
		if variableName != "" {
			// declare v = x
			// In Go, v has the type of the case. In JS, it's just the value.
			caseBody.Lines = append(caseBody.Lines, &js.Assign{
				Define: isShortDecl, // actually always true for type switch assign
				Lhs:    js.IdentP(variableName),
				Rhs:    tagExpr,
			})
		}

		for _, st := range cc.Body {
			caseBody.Lines = append(caseBody.Lines, statement(ctx, st))
		}

		if firstIf == nil {
			firstIf = &js.IfStmt{Cond: condition, Body: caseBody}
			currentIf = firstIf
		} else {
			nextIf := &js.IfStmt{Cond: condition, Body: caseBody}
			currentIf.Else = nextIf
			currentIf = nextIf
		}
	}

	// Handle default case
	for _, stmt := range s.Body.List {
		cc, ok := stmt.(*ast.CaseClause)
		if !ok || cc.List != nil {
			continue
		}

		// Default case body
		defaultBody := &js.Block{}
		if variableName != "" {
			defaultBody.Lines = append(defaultBody.Lines, &js.Assign{
				Define: isShortDecl,
				Lhs:    js.IdentP(variableName),
				Rhs:    tagExpr,
			})
		}
		for _, st := range cc.Body {
			defaultBody.Lines = append(defaultBody.Lines, statement(ctx, st))
		}

		if firstIf == nil {
			// Only default case
			// Just execute body
			for _, line := range defaultBody.Lines {
				block.Lines = append(block.Lines, line)
			}
			return block
		}

		// Attach to the last Else
		currentIf.Else = defaultBody
	}

	if firstIf != nil {
		block.Lines = append(block.Lines, firstIf)
	}

	return block
}

func generateTypeCheck(ctx *Context, tagExpr js.Expr, typeNode ast.Expr) js.Expr {
	// typeNode is the AST for the type (e.g. *ast.Ident "int", *ast.StructType, etc)
	// We need to convert it to a check expression.

	switch t := typeNode.(type) {
	case *ast.Ident:
		switch t.Name {
		case "int", "int8", "int16", "int32", "int64",
			"uint", "uint8", "uint16", "uint32", "uint64",
			"float32", "float64", "complex64", "complex128", "uintptr", "byte", "rune":
			// number
			return &js.Binary{
				X: &js.Unary{Op: js.TypeOf, X: tagExpr},
				Op: js.EqEq, // Use EqEq to match strict equality if possible, but Eq is "==". js.EqEq is "==="?
				Y: &js.BasicLit{Value: "'number'"},
			}
		case "string":
			return &js.Binary{
				X: &js.Unary{Op: js.TypeOf, X: tagExpr},
				Op: js.EqEq,
				Y: &js.BasicLit{Value: "'string'"},
			}
		case "bool":
			return &js.Binary{
				X: &js.Unary{Op: js.TypeOf, X: tagExpr},
				Op: js.EqEq,
				Y: &js.BasicLit{Value: "'boolean'"},
			}
		case "nil":
			// x == null
			return &js.Binary{
				X: tagExpr,
				Op: js.EqEq,
				Y: &js.BasicLit{Value: "null"},
			}
		default:
			// Assume it's a struct name or interface name
			// Check __goType
			// We must prepend package name if it's a local struct, because constructor does that.
			// Currently ctx.Module.Name has the package name.
			// But we need to know if t.Name is defined in this package.
			// For now, if it's a simple Ident (not a primitive), we assume it's a type in the current package
			// OR an interface name (which we don't fully support yet).
			// If it matches a struct in the current file/package, we should prefix it.

			// We can iterate ctx.Module.Decls to see if t.Name is a class.
			// But ctx.Module.Decls might not be fully populated yet if we are in the middle of processing?
			// Actually Decl processing happens before we print, but here we are in statement processing.
			// transpiler.go processes Decls sequentially?
			// Decl function calls `clsdecl`.

			// Let's just blindly prepend pkg name if ctx.Module.Name is set.
			// This matches `clsdecl` behavior which unconditionally sets PkgName.

			typeName := t.Name
			if ctx.Module.Name != nil {
				typeName = string(*ctx.Module.Name) + "." + typeName
			}

			return &js.Binary{
				X: &js.Selector{X: tagExpr, Sel: js.IdentP("__goType"), OptionalChain: true},
				Op: js.EqEq,
				Y: &js.BasicLit{Value: fmt.Sprintf("\"%s\"", typeName)},
			}
		}
	case *ast.StarExpr:
		// Pointer type.
		// In GoScript, pointers are often just values (except when taking address explicitly).
		// But if we have *Person, the value is just the Person object.
		// So we check the underlying type.
		// Wait, if it's *Person, the `__goType` on the object should still be "Person".
		// But `__goType` is attached to the object.
		// If I have `p := &Person{}`, `p` is the Person object in JS (passed by reference).
		// So `p.__goType` is "Person".
		// So `case *Person:` matches if `__goType == "Person"`.
		// What about `case Person:`? In Go, `Person` and `*Person` are different.
		// In JS transpilation, do we distinguish?
		// "Structs get __goType: pkg.StructName"
		// If I pass by value, it's a copy?
		// Current transpilation of struct copy is ... unknown. `assign` just references.
		// So `p1 = p2` is reference copy in JS.
		// So effectively everything is a pointer.
		// So `case Person` and `case *Person` might be ambiguous in this JS model.
		// For now, I'll ignore the star and check the type name.
		return generateTypeCheck(ctx, tagExpr, t.X)

	case *ast.SelectorExpr:
		// pkg.Type
		// TypeExpr returns ident name. For SelectorExpr it recurses.
		// TypeExpr implementation in transpiler.go:
		// case *ast.SelectorExpr: returns... wait, TypeExpr doesn't handle SelectorExpr in the snippet I saw!
		// Let's check TypeExpr in transpiler.go again.

		// It handles *ast.StarExpr, *ast.ArrayType, *ast.StructType, *ast.BasicLit, *ast.Ident.
		// It returns "" for others and warns.

		fullType := TypeExpr(ctx, t)
		return &js.Binary{
			X: &js.Selector{X: tagExpr, Sel: js.IdentP("__goType"), OptionalChain: true},
			Op: js.EqEq,
			Y: &js.BasicLit{Value: fmt.Sprintf("\"%s\"", fullType)},
		}

	default:
		ctx.Warn(typeNode, fmt.Sprintf("unsupported type in type switch: %T", typeNode))
		return &js.BasicLit{Value: "false"}
	}
}
