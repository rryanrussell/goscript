package link

type Kind int

const (
	Unknown Kind = iota
	Enum
	Struct
	Var
	Alias
	Array
	Ref
	BuiltIn
)

func (k Kind) String() string {
	return [...]string{"Unknown", "Enum", "Struct", "Var", "Alias", "Array", "Ref", "BuiltIn"}[k]
}

type Symbol struct {
	Members map[string]*Symbol
	Extends []string
	Kind    Kind
	Type    string
	Sub     string
}
