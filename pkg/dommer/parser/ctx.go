package parser

import (
	"errors"
	"io"

	"github.com/rryanrussell/goscript/pkg/dommer/ast"
)

type Never struct{}

type Ctx interface {
	Next(whitespace, comments bool) ast.Token

	Begin()
	Accept()
	Reject()

	Debug() *Debug
	Visit() *VisitContext

	Position() int
}

type ImplCtx struct {
	stack   Stack[int]
	scanner Tokenizer
	debug   Debug
	visit   VisitContext
}

func (c *ImplCtx) Next(whitespace, comments bool) ast.Token {
	for {
		tok, err := c.scanner.Next()

		switch {
		case errors.Is(err, io.EOF):
			return ast.Token{Class: ast.EOF}
		case err != nil:
			panic(err)
		case tok.Class == ast.Whitespace && whitespace:
		case tok.Class == ast.BlockComment && comments:
		case tok.Class == ast.LineComment && comments:
		case tok.Class == ast.BlockComment, tok.Class == ast.LineComment:
			continue
		case tok.Class == ast.Whitespace:
			continue
		}

		return tok
	}
}

func (c *ImplCtx) Begin() {

	if c.stack.Len() > 100 {
		panic("max parse depth exceeded")
	}

	c.stack.Push(c.scanner.Position())
}

func (c *ImplCtx) Accept() {
	_, ok := c.stack.Pop()
	if !ok {
		panic("no transaction")
	}
}

func (c *ImplCtx) Reject() {
	x, ok := c.stack.Pop()
	if !ok {
		panic("no transaction")
	}

	c.scanner.SetPosition(x)
}

func (c *ImplCtx) Position() int {
	return c.scanner.Position()
}

func (c *ImplCtx) Debug() *Debug {
	return &c.debug
}

func (c *ImplCtx) Visit() *VisitContext {
	return &c.visit
}

func NewCtx(scanner Tokenizer) Ctx {
	return &ImplCtx{scanner: scanner}
}
