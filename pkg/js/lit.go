package js

type Prop struct {
	Key   Expr
	Value Expr
}

func (p *Prop) Print(c PrintContext) {
	WriteAll(c, p.Key, Colon, Space, p.Value)
}

func (p *Prop) Walk(v Visitor) {
	v.Visit(p.Key)
	v.Visit(p.Value)
}

func (Prop) exprNode() {}

type ObjectLit struct {
	Type  string
	Props []*Prop
}

func (e ObjectLit) Walk(v Visitor) {
	Stroll(v, e.Props)
}

func (e ObjectLit) Print(c PrintContext) {
	WriteAll(c, Lbrc, Space)
	Join(c, e.Props, Comma, Space)
	WriteAll(c, Space, Rbrc)
}

type ArrayLit struct {
	Elts []Expr
}

func (e ArrayLit) Walk(v Visitor) {
	Stroll(v, e.Elts)
}

func (e ArrayLit) Print(c PrintContext) {
	WriteAll(c, Lbrkt, Space)
	Join(c, e.Elts, Comma, Space)
	WriteAll(c, Space, Rbrkt)
}

type BasicLit struct {
	Value string
}

func (BasicLit) Walk(v Visitor) {}

func (e *BasicLit) Print(c PrintContext) {
	c.WriteString(e.Value)
}

func (ObjectLit) exprNode() {}
func (ArrayLit) exprNode()  {}
func (BasicLit) exprNode()  {}
