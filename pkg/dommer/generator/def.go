package generator

type Def struct {
	Ident  *Ident
	Target Type
	Alias  bool
}

func (d *Def) Walk(v Visitor) {
	v.Visit(d.Ident)
	v.Visit(d.Target)
}

func (d *Def) Print(c PrintContext) {
	WriteAll(c, KType, Space, d.Ident, Space, PrintIf(d.Alias, Assign, Space), d.Target)
}

func (Def) IsType() {}
