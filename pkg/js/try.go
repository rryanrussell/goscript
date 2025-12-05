package js

type TryStmt struct {
	Body    *Block
	Catch   *CatchClause
	Finally *Block
}

func (s *TryStmt) Print(c PrintContext) {
	WriteAll(c, Try, Space, s.Body)
	if s.Catch != nil {
		WriteAll(c, Space, s.Catch)
	}
	if s.Finally != nil {
		WriteAll(c, Space, Finally, Space, s.Finally)
	}
}

func (s *TryStmt) Walk(v Visitor) {
	v.Visit(s.Body)
	if s.Catch != nil {
		v.Visit(s.Catch)
	}
	if s.Finally != nil {
		v.Visit(s.Finally)
	}
}

func (TryStmt) stmtNode() {}

type CatchClause struct {
	Param *Ident
	Body  *Block
}

func (c *CatchClause) Print(ctx PrintContext) {
	WriteAll(ctx, Catch, Space)
	if c.Param != nil {
		WriteAll(ctx, Lparen, c.Param, Rparen, Space)
	}
	WriteAll(ctx, c.Body)
}

func (c *CatchClause) Walk(v Visitor) {
	if c.Param != nil {
		v.Visit(c.Param)
	}
	v.Visit(c.Body)
}

type ThrowStmt struct {
	X Expr
}

func (s *ThrowStmt) Print(c PrintContext) {
	WriteAll(c, Throw, Space, s.X)
}

func (s *ThrowStmt) Walk(v Visitor) {
	v.Visit(s.X)
}

func (ThrowStmt) stmtNode() {}
