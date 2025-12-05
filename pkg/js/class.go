package js

import "strings"

type Class struct {
	Named
	PkgName string
	Fields  []*Field
}

func (cls *Class) Walk(v Visitor) {
	v.Visit(cls.Name)
	Stroll(v, cls.Fields)
}

func (cls *Class) Print(c PrintContext) {
	WriteAll(c,
		KClass,
		Space,
		cls.Name,
		Space,
		Lbrc,
	)

	c.Indent()

	ctor := &FuncDecl{Named: Named{Name: IdentP("constructor")}}
	ctor.Body = &Block{}
	ctor.Member = true

	// Add __goType
	typeName := string(*cls.Name)
	if cls.PkgName != "" {
		typeName = cls.PkgName + "." + typeName
	}
	ctor.Body.Lines = append(ctor.Body.Lines, &Assign{
		Lhs: &Selector{Sel: IdentP("__gs_goType"), X: IdentP(string(This))},
		Rhs: &BasicLit{Value: "\"" + typeName + "\""},
	})

	for _, f := range cls.Fields {
		WriteAll(c, Newline, Indent, f)

		pName := IdentP(string(*f.Name))
		pNameLower := IdentP(strings.ToLower(string(*f.Name)))
		ctor.Args = append(ctor.Args, &Var{Named: Named{pNameLower}})
		ctor.Body.Lines = append(ctor.Body.Lines, &Assign{Lhs: &Selector{Sel: pName, X: IdentP(string(This))}, Rhs: pNameLower})
	}

	WriteAll(c, Newline, Newline, Indent, ctor)

	c.Unindent()

	WriteAll(c, Rbrc, Newline)
}

func (Class) declNode() {}

type Field struct {
	Named
	Type  string
	Value string
}

func (v *Field) Walk(vis Visitor) {

}

func (v *Field) Print(c PrintContext) {
	val := v.Value
	if val == "" {
		val = "null"
	}
	WriteAll(c, v.Name, Eq, Token(val))
}
func (Field) exprNode() {}
