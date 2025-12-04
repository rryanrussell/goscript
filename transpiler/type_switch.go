package transpiler

import (
	"fmt"
	"go/ast"

	"github.com/rryanrussell/goscript/js"
)

func typeSwitch(ctx *Context, s *ast.TypeSwitchStmt) js.Stmt {
	block := &js.Block{}

	if s.Init != nil {
		block.Lines = append(block.Lines, statement(ctx, s.Init))
	}

	var initialTagExpr js.Expr
	var variableName string
	var isShortDecl bool

	switch assign := s.Assign.(type) {
	case *ast.ExprStmt:
		typeAssert := assign.X.(*ast.TypeAssertExpr)
		initialTagExpr = expr(ctx, typeAssert.X)
	case *ast.AssignStmt:
		typeAssert := assign.Rhs[0].(*ast.TypeAssertExpr)
		initialTagExpr = expr(ctx, typeAssert.X)
		if ident, ok := assign.Lhs[0].(*ast.Ident); ok {
			variableName = ident.Name
			isShortDecl = true
		}
	}

	tagVarName := "_tag"
	block.Lines = append(block.Lines, &js.Assign{
		Define: true,
		Lhs:    js.IdentP(tagVarName),
		Rhs:    initialTagExpr,
	})
	tagExpr := js.IdentP(tagVarName)

	var firstIf *js.IfStmt
	var currentIf *js.IfStmt

	for _, stmt := range s.Body.List {
		cc, ok := stmt.(*ast.CaseClause)
		if !ok {
			continue
		}

		var condition js.Expr

		if cc.List == nil {
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

		caseBody := &js.Block{}
		if variableName != "" {
			caseBody.Lines = append(caseBody.Lines, &js.Assign{
				Define: isShortDecl,
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

	for _, stmt := range s.Body.List {
		cc, ok := stmt.(*ast.CaseClause)
		if !ok || cc.List != nil {
			continue
		}

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
			for _, line := range defaultBody.Lines {
				block.Lines = append(block.Lines, line)
			}
			return block
		}

		currentIf.Else = defaultBody
	}

	if firstIf != nil {
		block.Lines = append(block.Lines, firstIf)
	}

	return block
}

func generateTypeCheck(ctx *Context, tagExpr js.Expr, typeNode ast.Expr) js.Expr {
	switch t := typeNode.(type) {
	case *ast.Ident:
		switch t.Name {
		case "int", "int8", "int16", "int32", "int64",
			"uint", "uint8", "uint16", "uint32", "uint64",
			"float32", "float64", "complex64", "complex128", "uintptr", "byte", "rune":
			return &js.Binary{
				X: &js.Unary{Op: js.TypeOf, X: tagExpr},
				Op: js.EqEq,
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
			return &js.Binary{
				X: tagExpr,
				Op: js.EqEq,
				Y: &js.BasicLit{Value: "null"},
			}
		default:
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
		return generateTypeCheck(ctx, tagExpr, t.X)

	case *ast.SelectorExpr:
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
