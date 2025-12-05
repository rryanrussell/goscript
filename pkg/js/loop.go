package js

type WhileStmt struct {
	Cond Expr
	Body *Block
}

func (i *WhileStmt) Print(c PrintContext) {
	WriteAll(c, While, Space, Lparen, i.Cond, Rparen, i.Body)
}

func (i *WhileStmt) Walk(v Visitor) {
	v.Visit(i.Cond)
	v.Visit(i.Body)
}

type ForStmt struct {
	Init Stmt
	Cond Expr
	Post Stmt
	Body *Block
}

func (f *ForStmt) Print(c PrintContext) {
	WriteAll(c, For, Space, Lparen,
		PrintIf(f.Init, f.Init), Semi,
		PrintIf(f.Cond, f.Cond), Semi,
		PrintIf(f.Post, f.Post),
		Rparen,
		f.Body,
	)
}

func (f *ForStmt) Walk(v Visitor) {
	if f.Init != nil {
		v.Visit(f.Init)
	}
	if f.Cond != nil {
		v.Visit(f.Cond)
	}
	if f.Post != nil {
		v.Visit(f.Post)
	}

	v.Visit(f.Body)
}

type ForEach struct {
	Key      *VarDecl
	Value    *VarDecl
	Iterable Expr
	Body     *Block
}

func (f *ForEach) Walk(v Visitor) {
	if f.Key != nil {
		v.Visit(f.Key)
	}
	if f.Value != nil {
		v.Visit(f.Value)
	}

	v.Visit(f.Iterable)

	v.Visit(f.Body)
}

func (f *ForEach) Print(c PrintContext) {
	if f.Key != nil && f.Value != nil {
		WriteAll(c,
			For, Space, Lparen,
			Let, Space, Lbrkt, f.Key.Var, Comma, Space, f.Value.Var, Rbrkt,
			Space, Of, Space, f.Iterable,
			Rparen, f.Body,
		)
		return
	}
	WriteAll(c, For, Space, Lparen,
		PrintIf(f.Key, f.Key, Space, In),
		PrintIf(f.Value, f.Value, Space, Of),
		Space,
		f.Iterable,
		Rparen,
		f.Body,
	)
}

func (WhileStmt) stmtNode() {}
func (ForStmt) stmtNode()   {}
func (ForEach) stmtNode()   {}
