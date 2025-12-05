package generator

type Module struct {
	Ident   *Ident
	Structs []*Struct
	Enums   []*Enum
	Defs    []*Def
	Vars    []*Var
	Funcs   []*Method
}

func (m *Module) Walk(v Visitor) {
	v.Visit(m.Ident)
	Stroll(v, m.Structs)
	Stroll(v, m.Enums)
	Stroll(v, m.Defs)
	Stroll(v, m.Vars)
	Stroll(v, m.Funcs)
}

func (m *Module) Print(c PrintContext) {
	WriteAll(c, KPac, Space, m.Ident, Newline, Newline)

	Join(c, m.Enums, Newline, Newline)

	Newline.Print(c)
	Newline.Print(c)

	Join(c, m.Structs, Newline, Newline)

	Newline.Print(c)
	Newline.Print(c)

	Join(c, m.Defs, Newline)

	Newline.Print(c)
	Newline.Print(c)

	Join(c, m.Vars, Newline)

	Newline.Print(c)
	Newline.Print(c)

	Join(c, m.Funcs, Newline)
}
