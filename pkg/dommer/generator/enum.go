package generator

type Enum struct {
	Ident   *Ident
	Members []*EnumMember
	Kind    EnumKind
}

type EnumKind int

const (
	StringKind EnumKind = iota
)

type EnumMember struct {
	Ident *Ident
	Value string
}

func (m *EnumMember) Walk(v Visitor) {
	v.Visit(m.Ident)
}

func (m *EnumMember) Print(c PrintContext) {
	m.Ident.Print(c)
}

func (Enum) IsType() {}

func (e *Enum) Print(c PrintContext) {
	if e.Kind != StringKind {
		panic("not implemented")
	}

	WriteAll(c, KType, Space, e.Ident, Space, KString,
		Newline,
		Newline,
		KVar, Lparen,
	)

	c.Indent()

	for _, member := range e.Members {
		WriteAll(
			c, Newline, Indent,
			member.Ident, Space, e.Ident, Space, Assign, Space, DblQte, S(member.Value), DblQte,
		)
	}

	c.Unindent()

	WriteAll(c, Newline, Rparen, Newline)
}

func (s *Enum) Walk(v Visitor) {
	v.Visit(s.Ident)
	Stroll(v, s.Members)
}
