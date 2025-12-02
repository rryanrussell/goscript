package js

type Stmt interface {
	Node
	stmtNode()
}

type Return struct {
	X Expr
}

func (r *Return) Print(c PrintContext) {
	WriteAll(c, Ret, PrintIf(r.X, Space, r.X))
}

func (r *Return) Walk(v Visitor) {
	v.Visit(r.X)
}

type ExprStmt struct {
	X Expr
}

func (s *ExprStmt) Print(c PrintContext) {
	if s.X != nil {
		WriteAll(c, s.X)
	}
}

func (s *ExprStmt) Walk(v Visitor) {
	if s.X != nil {
		v.Visit(s.X)
	}
}

type Assign struct {
	Define bool
	Lhs    Expr
	Rhs    Expr
}

func (a *Assign) Print(c PrintContext) {
	WriteAll(c, PrintIf(a.Define, Let, Space), a.Lhs, Space, Eq, Space, a.Rhs)
}

func (a *Assign) Walk(v Visitor) {
	v.Visit(a.Lhs)
	v.Visit(a.Rhs)
}

type IfStmt struct {
	Cond Expr
	Body *Block
}

func (i *IfStmt) Print(c PrintContext) {
	WriteAll(c, If, Space, Lparen, i.Cond, Rparen, i.Body)
}

func (i *IfStmt) Walk(v Visitor) {
	v.Visit(i.Cond)
	v.Visit(i.Body)
}

func (ExprStmt) stmtNode() {}
func (IfStmt) stmtNode()   {}
func (Assign) stmtNode()   {}
func (Return) stmtNode()   {}
