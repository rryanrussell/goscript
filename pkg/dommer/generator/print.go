package generator

import (
	"io"
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
	default:
		return Printers(items)
	}
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

func CleanName(str string) string {
	switch str {
	case "type", "func":
		return str + "0"
	default:
		return str
	}
}

func ExportName(str string) string {
	str = CleanName(str)

	if len(str) > 0 && str[0] >= 'a' && str[0] <= 'z' {
		return string(str[0]+'A'-'a') + str[1:]
	}

	if len(str) > 0 && str[0] >= '0' && str[0] <= '9' {
		return "_" + str
	}

	return str
}
