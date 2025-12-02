package js

type FuncDecl struct {
	Named
	Body   *Block
	Args   []*Var
	Member bool
}

func (f *FuncDecl) Print(c PrintContext) {
	WriteAll(c,
		PrintIf(
			!f.Member,
			Function,
			Space),
		f.Name,
		Lparen,
	)

	Join(c, f.Args, Comma, Space)

	Tok(c, Rparen)

	f.Body.Print(c)
}

func (f *FuncDecl) Walk(v Visitor) {
	v.Visit(f.Name)
	Stroll(v, f.Args)
	v.Visit(f.Body)
}

func (FuncDecl) exprNode() {}
func (FuncDecl) declNode() {}
