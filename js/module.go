package js

type Module struct {
	Named

	Decls []Decl
}

type Program struct {
	Modules []*Module
}

func (m *Module) Walk(v Visitor) {
	v.Visit(m.Name)
	Stroll(v, m.Decls)
}

func (m *Module) Print(c PrintContext) {
	Join(c, m.Decls, Newline)
}
