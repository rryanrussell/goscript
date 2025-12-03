package js

type SwitchStmt struct {
	Tag  Expr
	Body *CaseBlock
}

type CaseBlock struct {
	List []Stmt // List of CaseClause
}

type CaseClause struct {
	List []Expr
	Body []Stmt
}

func (s *SwitchStmt) Print(c PrintContext) {
	WriteAll(c, Switch, Space, Lparen)
	if s.Tag != nil {
		WriteAll(c, s.Tag)
	} else {
		// JS switch requires an expression. If Go switch has no tag, it defaults to true.
		WriteAll(c, Token("true"))
	}
	WriteAll(c, Rparen, Space, s.Body)
}

func (s *SwitchStmt) Walk(v Visitor) {
	if s.Tag != nil {
		v.Visit(s.Tag)
	}
	v.Visit(s.Body)
}

func (s *SwitchStmt) stmtNode() {}

func (c *CaseBlock) Print(ctx PrintContext) {
	WriteAll(ctx, Lbrc)
	ctx.Indent()
	for _, s := range c.List {
		WriteAll(ctx, Newline, Indent, s)
	}
	ctx.Unindent()
	WriteAll(ctx, Newline, Indent, Rbrc)
}

func (c *CaseBlock) Walk(v Visitor) {
	for _, s := range c.List {
		v.Visit(s)
	}
}

func (c *CaseBlock) stmtNode() {}

func (c *CaseClause) Print(ctx PrintContext) {
	if c.List == nil {
		WriteAll(ctx, Default, Colon)
	} else {
		for i, expr := range c.List {
			if i > 0 {
				WriteAll(ctx, Newline, Indent)
			}
			WriteAll(ctx, Case, Space, expr, Colon)
		}
	}

	ctx.Indent()
	for _, s := range c.Body {
		WriteAll(ctx, Newline, Indent, s)
	}
	ctx.Unindent()
}

func (c *CaseClause) Walk(v Visitor) {
	for _, expr := range c.List {
		v.Visit(expr)
	}
	for _, s := range c.Body {
		v.Visit(s)
	}
}

func (c *CaseClause) stmtNode() {}
