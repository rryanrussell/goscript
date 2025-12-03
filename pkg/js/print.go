package js

import (
	"io"
	"reflect"
)

type Writer interface {
	io.Writer
	WriteString(s string) (int, error)
}

type Indenter interface {
	Indent()
	Unindent()
	Depth() int
}

type WriterIndenter interface {
	Writer
	Indenter
}

type PrintContext = WriterIndenter

type Printer interface {
	Print(c PrintContext)
}

func WriteAll(c PrintContext, items ...Printer) {
	for _, item := range items {
		item.Print(c)
	}
}

type Printers []Printer

func (p Printers) Print(c PrintContext) {
	WriteAll(c, p...)
}

func PrintIf(cond any, items ...Printer) Printer {
	switch cond {
	case nil, false:
		return Printers(nil)
	case true:
		return Printers(items)
	default:
		rcond := reflect.ValueOf(cond)
		if rcond.IsNil() {
			return Printers(nil)
		}
		return Printers(items)
	}
}

func Tok(c PrintContext, t Token) {
	t.Print(c)
}

func Join[T Printer](c PrintContext, items []T, sep ...Printer) {
	for i, item := range items {
		item.Print(c)
		if i < len(items)-1 {
			for _, sep := range sep {
				sep.Print(c)
			}
		}
	}
}
