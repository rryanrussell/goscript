package link

import (
	"go/ast"
	"go/token"
	"regexp"

	x "github.com/rryanrussell/goscript/pkg/transpiler"
)

func ensureMap[T comparable, U any](m *map[T]U) {
	if *m == nil {
		*m = make(map[T]U)
	}
}

func kind(e ast.Expr) Kind {
	switch e.(type) {
	case *ast.Ident:
		return Enum
	case *ast.StructType:
		return Struct
	}

	return Unknown
}

func rename(cg *ast.CommentGroup) string {
	if cg == nil {
		return ""
	}

	re := regexp.MustCompile(`rename:((new )?\S+)`)
	match := re.FindStringSubmatch(cg.Text())

	if len(match) < 2 {
		return ""
	}

	return match[1]
}

func vars(c *Context, node *ast.GenDecl) {
	for _, spec := range node.Specs {
		v, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}

		symbol := &Symbol{
			Sub: rename(node.Doc),
		}

		symbol.Type = x.TypeExpr(c, v.Type)

		switch {
		case len(v.Values) > 1:
			panic("unexpected tuple")
		case len(v.Values) == 0:
		case symbol.Sub != "":
		default:
			symbol.Sub = x.TypeExpr(c, v.Values[0])
		}

		name := v.Names[0].Name

		t, ok := c.Symbols[symbol.Type]

		if ok {
			ensureMap(&t.Members)
			t.Members[name] = symbol
		}

		ensureMap(&c.Symbols)
		c.Symbols[name] = symbol
	}
}

func fields(c *Context, t *Symbol, node *ast.TypeSpec) {
	s, ok := node.Type.(*ast.StructType)
	if !ok {
		return
	}

	for _, field := range s.Fields.List {
		if len(field.Names) == 0 {
			t.Extends = append(t.Extends, x.TypeExpr(c, field.Type))
			continue
		}
		if len(field.Names) > 1 {
			panic("unexpected field name count")
		}

		ensureMap(&t.Members)
		t.Members[field.Names[0].Name] = &Symbol{
			Type: x.TypeExpr(c, field.Type),
			Sub:  rename(field.Doc),
			Kind: Ref,
		}
	}
}

func types(c *Context, node *ast.GenDecl) {
	for _, spec := range node.Specs {
		v, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}

		ensureMap(&c.Symbols)
		var t Symbol

		if v.Assign.IsValid() {
			t.Kind = Alias
		} else {
			t.Kind = kind(v.Type)
		}

		t.Type = v.Name.Name

		fields(c, &t, v)
		c.Symbols[v.Name.Name] = &t
	}
}

func fun(c *Context, node *ast.FuncDecl) {
	t := ""

	if node.Type.Results != nil && len(node.Type.Results.List) > 0 {
		t = x.TypeExpr(c, node.Type.Results.List[0].Type)
	}

	symbol := &Symbol{
		Sub:  rename(node.Doc),
		Type: t,
	}

	name := node.Name.Name

	if node.Recv != nil && len(node.Recv.List) > 0 {
		recv := x.TypeExpr(c, node.Recv.List[0].Type)
		ensureMap(&c.Symbols[recv].Members)
		c.Symbols[recv].Members[name] = symbol
	} else {
		ensureMap(&c.Symbols)
		c.Symbols[name] = symbol
	}
}

type Context struct {
	x.Complainer
	Symbols map[string]*Symbol
}

func Symbolize(complainer x.Complainer, file *ast.File) map[string]*Symbol {
	c := &Context{Complainer: complainer}
	for _, node := range file.Decls {
		switch decl := node.(type) {
		case *ast.GenDecl:
			if decl.Tok == token.TYPE {
				types(c, decl)
			}
		}
	}

	for _, node := range file.Decls {
		switch decl := node.(type) {
		case *ast.GenDecl:
			if decl.Tok == token.VAR {
				vars(c, decl)
			}
		case *ast.FuncDecl:
			fun(c, decl)
		}
	}

	return c.Symbols
}
