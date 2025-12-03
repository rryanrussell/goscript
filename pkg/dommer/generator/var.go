package generator

type ValueExpr = string

type Var struct {
	Ident *Ident
	Type  Type
	Value ValueExpr
}

func (v *Var) Walk(vis Visitor) {
	vis.Visit(v.Ident)
	vis.Visit(v.Type)
}

func (v *Var) Print(c PrintContext) {
	WriteAll(c, KVar, Space, v.Ident, Space, v.Type, Space, PrintIf(len(v.Value) > 0, Assign, Space, S(v.Value)))
}
