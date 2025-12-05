package js

type Expr interface {
	Node
	exprNode()
}

type Binary struct {
	X  Expr
	Op Token
	Y  Expr
}

func (e Binary) Walk(v Visitor) {
	v.Visit(e.X)
	v.Visit(e.Y)
}

func (e Binary) Print(c PrintContext) {
	WriteAll(c, e.X, e.Op, e.Y)
}

type Selector struct {
	X             Expr
	Sel           *Ident
	OptionalChain bool
}

func (s *Selector) GetName() Ident {
	return *s.Sel
}

func (s *Selector) SetName(name Ident) {
	s.Sel = &name
}

func (e *Selector) Print(c PrintContext) {
	WriteAll(c, e.X, PrintIf(e.OptionalChain, QstnMrk), Dot, e.Sel)
}

func (e *Selector) Walk(v Visitor) {
	v.Visit(e.X)
	v.Visit(e.Sel)
}

type Index struct {
	X, Index Expr
}

func (e Index) Walk(v Visitor) {
	v.Visit(e.X)
	v.Visit(e.Index)
}

func (e Index) Print(c PrintContext) {
	WriteAll(c, e.X, Lbrkt, e.Index, Rbrkt)
}

func (Ident) exprNode()    {}
func (Index) exprNode()    {}
func (Selector) exprNode() {}

type AwaitExpr struct {
	X Expr
}

func (e *AwaitExpr) Print(c PrintContext) {
	WriteAll(c, Token("await"), Space, e.X)
}

func (e *AwaitExpr) Walk(v Visitor) {
	v.Visit(e.X)
}

func (AwaitExpr) exprNode() {}
