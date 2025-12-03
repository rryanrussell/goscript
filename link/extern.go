package link

import (
	"fmt"
	"strings"

	"github.com/rryanrussell/goscript/js"
	x "github.com/rryanrussell/goscript/transpiler"
)

type LinkContext struct {
	Extern    map[string]*Symbol
	Semantics map[js.Node]x.Semantics
	Warns     []error
}

func (c *LinkContext) Warn(err error) {
	c.Warns = append(c.Warns, err)
}

type Scope struct {
	*LinkContext
	Symbols map[string]*Symbol
	Scopes  []*Scope
	Parent  *Scope
}

func LinkExternals(c *LinkContext, program *js.Module) {
	vis := &Scope{
		LinkContext: c,
		Symbols:     make(map[string]*Symbol),
	}
	vis.Visit(program)
}

func (parent *Scope) NewChildScope() *Scope {
	child := &Scope{
		LinkContext: parent.LinkContext,
		Symbols:     make(map[string]*Symbol),
	}
	child.Parent = parent
	parent.Scopes = append(parent.Scopes, child)
	return child
}

func (scope *Scope) search(name string) *Symbol {
	sym := scope.Symbols[name]

	if sym != nil {
		return sym
	}

	if scope.Parent != nil {
		return scope.Parent.search(name)
	}

	if scope.Extern != nil {
		sym = scope.Extern[name]
		if sym != nil {
			return sym
		}
	}

	return nil
}

func (scope *Scope) Visit(node js.Node) {
	mod, ok := node.(*js.Module)

	if ok {
		mod.Walk(scope.NewChildScope())
		return
	}

	fun, ok := node.(*js.FuncDecl)

	if ok {
		funScope := scope.NewChildScope()
		for _, arg := range fun.Args {
			name := string(*arg.Name)
			funScope.Symbols[name] = scope.parseType(arg.Type)
		}
		fun.Body.Walk(funScope)
	}

	block, ok := node.(*js.Block)

	if ok {
		block.Walk(scope.NewChildScope())
		return
	}

	namer, ok := node.(js.Namer)
	if ok && namer != nil {
		name := string(namer.GetName())

		sym := scope.bind(name, node)

		if sym == nil {
			scope.Symbols[name] = &Symbol{}
		} else if sym.Sub != "" {
			namer.SetName(js.Ident(sym.Sub))
		}
	}

	assign, ok := node.(*js.Assign)

	if ok {
		if namer, ok := assign.Lhs.(js.Namer); ok {
			name := string(namer.GetName())
			scope.Symbols[name] = scope.typeOf(assign.Rhs)
		}
	}

	v, ok := node.(*js.VarDecl)

	if ok {
		scope.Symbols[string(*v.Var.Name)] = scope.parseType(v.Var.Type)
	}

	f, ok := node.(*js.Field)

	if ok {
		f.Value = scope.zero(scope.parseType(f.Type))
	}

	if node != nil {
		node.Walk(scope)
	}
}

func (scope Scope) bind(name string, node js.Node) *Symbol {
	sel, ok := node.(*js.Selector)

	if ok {
		return scope.selector(sel)
	}

	sym := scope.search(name)

	if sym == nil {
		sym = scope.typeOf(node)
	}

	return sym
}

func (scope Scope) parseType(t string) *Symbol {
	parts := strings.Split(t, "[]")
	var sym *Symbol
	if len(parts) == 1 {
		sym = scope.search(t)
	} else {
		elem := scope.search(parts[0])
		if elem != nil {
			sym = &Symbol{
				Kind: Array,
				Members: map[string]*Symbol{
					"0": elem,
				},
			}
		}
	}

	if sym == nil {
		scope.Warn(fmt.Errorf("could not parse type %q", t))
	}

	return sym
}

func (scope Scope) resolveMembers(t *Symbol) map[string]*Symbol {
	if t == nil {
		return nil
	}

	if t.Kind == Ref {
		t = scope.parseType(t.Type)
		if t == nil {
			return nil
		}
	}

	if len(t.Extends) == 0 {
		return t.Members
	}

	combined := make(map[string]*Symbol)

	for _, parentName := range t.Extends {
		parent := scope.search(parentName)

		if parent == nil {
			continue
		}

		inherited := scope.resolveMembers(parent)

		for k, s := range inherited {
			combined[k] = s
		}
	}

	for k, s := range t.Members {
		combined[k] = s
	}

	return combined
}

func (scope Scope) selector(x *js.Selector) *Symbol {
	id := string(*x.Sel)
	return scope.resolveMembers(scope.typeOf(x.X))[id]
}

func (scope Scope) typeOf(node js.Node) *Symbol {
	assert := scope.Semantics[node].TypeAssertion
	if assert != "" {
		return scope.parseType(assert)
	}

	switch x := node.(type) {
	case *js.Ref:
		id := string(*x.Name)
		t := scope.search(id)
		if t != nil && t.Type != "" {
			return scope.search(t.Type)
		}
		return t
	case *js.Selector:
		return scope.selector(x)
	case *js.Call:
		funcType := scope.typeOf(x.Func)
		if funcType != nil {
			return scope.parseType(funcType.Type)
		}
	case *js.Index:
		operand := scope.typeOf(x.X)
		return operand.Members["0"]
	}

	return nil
}

func (scope Scope) zero(sym *Symbol) string {
	if sym == nil {
		return ""
	}

	switch sym.Kind {
	case Array:
		return "null"
	case BuiltIn:
		if sym.Type == "float64" {
			return "0.0"
		}
	}

	return ""
}
