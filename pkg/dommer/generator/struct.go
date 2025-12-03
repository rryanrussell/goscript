package generator

type Struct struct {
	Ident   *Ident
	Extends []Type
	Fields  []*Field
	Methods []*Method
}

func (Struct) IsType() {}

func (s *Struct) Print(c PrintContext) {
	c.Indent()

	WriteAll(c, KType, Space, s.Ident, Space, KStruct, Space, Lbrc, Newline, Indent)

	for _, t := range s.Extends {
		WriteAll(c, t, Newline, Indent)
	}

	Join(c, s.Fields, Newline, Indent)

	c.Unindent()

	WriteAll(c, Newline, Rbrc, Newline, Newline)

	Join(c, s.Methods, Newline)
}

func (s *Struct) Walk(v Visitor) {
	v.Visit(s)
	v.Visit(s.Ident)
	Stroll(v, s.Fields)
	Stroll(v, s.Methods)
}
