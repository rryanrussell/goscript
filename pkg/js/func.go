package js

type FuncDecl struct {
	Named
	Body   *Block
	Args   []*Var
	Member bool
	Async  bool
}

func (f *FuncDecl) Print(c PrintContext) {
	WriteAll(c,
		PrintIf(f.Async, Token("async"), Space),
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

type ArrowFunc struct {
	Body  *Block
	Args  []*Var
	Async bool
}

func (f *ArrowFunc) Print(c PrintContext) {
	WriteAll(c,
		PrintIf(f.Async, Token("async"), Space),
		Lparen,
	)
	Join(c, f.Args, Comma, Space)
	WriteAll(c, Rparen, Space, Token("=>"), Space, f.Body)
}

func (f *ArrowFunc) Walk(v Visitor) {
	Stroll(v, f.Args)
	v.Visit(f.Body)
}

func (ArrowFunc) exprNode() {}
