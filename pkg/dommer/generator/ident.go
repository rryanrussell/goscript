package generator

import "fmt"

type NameConvention int

const (
	NameExport NameConvention = iota
	NameParameter
	NameModule
)

type Ident struct {
	Name           string
	SourceAlias    string
	NameConvention NameConvention
}

func (id *Ident) PrintRename(c PrintContext) bool {
	switch {
	case id.NameConvention != NameExport:
		return false
	case id.Name == ExportName(id.Name):
		return false
	case id.SourceAlias != "" && id.SourceAlias != id.Name:
		fmt.Fprintf(c, "// rename:%s", id.SourceAlias)
	default:
		fmt.Fprintf(c, "// rename:%s", id.Name)
	}

	return true
}

func (id *Ident) Print(c PrintContext) {
	switch id.NameConvention {
	case NameExport:
		c.WriteString(ExportName(id.Name))
	case NameModule, NameParameter:
		c.WriteString(CleanName(id.Name))
	}
}

func (id *Ident) Walk(visitor Visitor) {
	visitor.Visit(id)
}
