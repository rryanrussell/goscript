package js

type Call struct {
	Func Expr
	Args []Expr
}

func (e Call) Print(c PrintContext) {
	WriteAll(c, e.Func, Lparen)
	Join(c, e.Args, Comma, Space)
	Tok(c, Rparen)
}

func (e Call) Walk(v Visitor) {
	v.Visit(e.Func)
	Stroll(v, e.Args)
}
func (Call) exprNode() {}

type New struct {
	Class Expr
	Args  []Expr
}

func (e New) Print(c PrintContext) {
	WriteAll(c, KNew, Space, e.Class, Lparen)
	Join(c, e.Args, Comma, Space)
	Tok(c, Rparen)
}

func (e New) Walk(v Visitor) {
	v.Visit(e.Class)
	Stroll(v, e.Args)
}
func (New) exprNode() {}
