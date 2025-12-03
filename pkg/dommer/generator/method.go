package generator

type Param struct {
	Ident *Ident
	Type  Type
}

func (p *Param) Walk(v Visitor) {
	v.Visit(p.Ident)
	v.Visit(p.Type)
}

func (p *Param) Print(c PrintContext) {
	WriteAll(c, p.Ident, Space, p.Type)
}

type Method struct {
	Ident  *Ident
	Recv   Type
	Params []*Param
	Type   Type
}

func (m *Method) Walk(v Visitor) {
	v.Visit(m.Type)
	Stroll(v, m.Params)
}

func (m *Method) Print(c PrintContext) {
	m.Ident.PrintRename(c)

	WriteAll(c, Newline, Indent, KFunc, Space, PrintIf(m.Recv, Lparen, m.Recv, Rparen, Space), m.Ident, Lparen)

	Join(c, m.Params, Comma, Space)

	WriteAll(c, Rparen, Space, m.Type, Space, Lbrc, Space, S(`panic("extern")`), Space, Rbrc)
}
