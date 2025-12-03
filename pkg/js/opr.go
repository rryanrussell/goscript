package js

type Unary struct {
	Op Token
	X  Expr
}

func (e *Unary) Walk(v Visitor) {
	v.Visit(e.X)
}

func (e *Unary) Print(c PrintContext) {
	WriteAll(c, e.Op, e.X)
}

type Parens struct {
	X Expr
}

func (e *Parens) Walk(v Visitor) {
	v.Visit(e.X)
}

func (e *Parens) Print(c PrintContext) {
	WriteAll(c, Lparen, e.X, Rparen)
}

type Spread struct {
	X Expr
}

func (e *Spread) Walk(v Visitor) {
	v.Visit(e.X)
}

func (e *Spread) Print(c PrintContext) {
	WriteAll(c, Dot, Dot, Dot, e.X)
}

func (Unary) exprNode()  {}
func (Binary) exprNode() {}
func (Parens) exprNode() {}
func (Spread) exprNode() {}
