package js

type TokenStmt struct {
	Token Token
}

func (t *TokenStmt) Print(c PrintContext) {
	WriteAll(c, t.Token)
}

func (t *TokenStmt) Walk(v Visitor) {}

func (t *TokenStmt) stmtNode() {}
