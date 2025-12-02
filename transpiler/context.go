package transpiler

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"os"

	"github.com/rryanrussell/goscript/js"
)

type Semantics struct {
	TypeAssertion   string
	PassByReference bool
}

type ModuleContext struct {
	*TranspileContext
	Module    js.Module
	Semantics map[js.Node]Semantics
	depth     int
}

func (c *ModuleContext) Indent() {
	c.depth++
}

func (c *ModuleContext) Unindent() {
	if c.depth == 0 {
		panic("indent underflow")
	}
	c.depth--
}

func (c *ModuleContext) Depth() int {
	return c.depth
}

type TranspileContext struct {
	io.Writer
	main  bool
	Fset  *token.FileSet
	Warns []Msg
	Errs  []Msg
}

func (c *TranspileContext) WriteString(s string) (n int, err error) {
	return io.WriteString(c.Writer, s)
}

func (c *TranspileContext) Main() {
	c.main = true
}

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[38;5;240m"
	GrayBg = "\033[48;5;240m"
	White  = "\033[97m"
)

func (c *TranspileContext) Exit() {
	print(Yellow)

	for _, warn := range c.Warns {
		if warn.at == nil {
			fmt.Fprintf(os.Stderr, "WARN: %s\n\n", warn.msg)
			continue
		}
		pos := c.Fset.Position(warn.at.Pos())
		end := c.Fset.Position(warn.at.End())
		fmt.Fprintf(os.Stderr, "WARN: %s at \n%v to %v\n\n", warn.msg, pos, end)
	}

	// var code int

	print(Red)

	for _, err := range c.Errs {
		// code = 1
		if err.at == nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n\n", err.msg)
			continue
		}
		pos := c.Fset.Position(err.at.Pos())
		end := c.Fset.Position(err.at.End())
		fmt.Fprintf(os.Stderr, "ERROR: %s at \n%v to %v\n\n", err.msg, pos, end)
	}

	print(Reset)

	// os.Exit(code)
}

type Complainer interface {
	Warn(at ast.Node, msg string)
	Error(at ast.Node, msg string)
	Fatal(at ast.Node, msg string)
}

func (c *TranspileContext) Warn(at ast.Node, msg string) {
	c.Warns = append(c.Warns, Msg{at, msg})
}

func (c *TranspileContext) Error(at ast.Node, msg string) {
	c.Errs = append(c.Errs, Msg{at, msg})
}

func (c *TranspileContext) Fatal(at ast.Node, msg string) {
	c.Error(at, msg)
	c.Exit()
}

type Msg struct {
	at  ast.Node
	msg string
}
