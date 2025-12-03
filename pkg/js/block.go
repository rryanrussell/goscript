package js

type Block struct {
	Lines  []Stmt
	Inline bool
}

func (b *Block) Print(c PrintContext) {
	WriteAll(c,
		Space,
		Lbrc,
	)

	c.Indent()

	for _, line := range b.Lines {
		WriteAll(c, Newline, Indent, line)
	}

	c.Unindent()

	WriteAll(c,
		Newline,
		Indent,
		Rbrc,
		PrintIf(!b.Inline, Newline),
	)
}

func (b *Block) Walk(v Visitor) {
	Stroll(v, b.Lines)
}
