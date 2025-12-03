package generator

type Field struct {
	Ident *Ident
	Type  Type
	Tag   string
}

func (f Field) Walk(v Visitor) {
	v.Visit(f)
	v.Visit(f.Ident)
	v.Visit(f.Type)
}

func (f Field) Print(c PrintContext) {
	if f.Ident.PrintRename(c) {
		WriteAll(c, Newline, Indent)
	}

	WriteAll(c, f.Ident, Space, f.Type)

	if f.Tag != "" {
		WriteAll(c, Space, Backtick, S(f.Tag), Backtick)
	}
}
