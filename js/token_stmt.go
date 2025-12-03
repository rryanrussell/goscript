package js

type TokenStmt struct {
	Token Token
	Label *Ident
}

func (t *TokenStmt) Print(c PrintContext) {
	WriteAll(c, t.Token)
	if t.Label != nil {
		WriteAll(c, Space, t.Label)
	}
}

func (t *TokenStmt) Walk(v Visitor) {}

func (t *TokenStmt) stmtNode() {}
