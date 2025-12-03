package js

type Var struct {
	Named
	Type string
}

func (v *Var) Walk(vis Visitor) {

}

func (v *Var) Print(c PrintContext) {
	v.Name.Print(c)
}
func (Var) exprNode() {}

type VarDecl struct {
	Var    *Var
	Const  bool
	Value  Expr
	Inline bool
}

func (v *VarDecl) Walk(vis Visitor) {
	vis.Visit(v.Var)
	if v.Value != nil {
		vis.Visit(v.Value)
	}
}

func (v *VarDecl) Print(c PrintContext) {
	WriteAll(c,
		PrintIf(v.Const, Const), PrintIf(!v.Const, Let), Space,
		v.Var.Name, PrintIf(v.Value != nil, Eq, v.Value), PrintIf(!v.Inline, Newline),
	)
}

func (VarDecl) declNode() {}
