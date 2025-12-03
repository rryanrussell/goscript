package dommer

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/rryanrussell/goscript/pkg/dommer/ast"
	g "github.com/rryanrussell/goscript/pkg/dommer/generator"
)

type Name struct {
	ns    *g.Ident
	value string
}

type Binding struct {
	ident *g.Ident
	bound bool
}

type TypeRef struct {
	ident   *g.Ident
	binding g.Type
}

func (t *TypeRef) Walk(v g.Visitor) {
	v.Visit(t.ident)
}

func (m *TypeRef) Print(c g.PrintContext) {
	m.ident.Print(c)
}

func (m *TypeRef) IsType() {}

type GenerateContext struct {
	module g.Module

	names map[Name]*Binding
	types map[Name]g.Type
}

var JSAny = ast.Any.Value
var JSVoid = ast.Void.Value
var JSString = ast.String.Value
var JSNull = ast.Null.Value
var JSNumber = ast.Number.Value
var JSBoolean = ast.Boolean.Value
var JSFloat32List = "Float32List"
var JSInt32List = "Int32List"

var Global *g.Ident = nil

type BuiltInType struct {
	g.Symbol
}

func (t *BuiltInType) Walk(v g.Visitor) {}
func (m *BuiltInType) IsType()          {}

func (c *GenerateContext) Error(values ...any) {
	panic(errors.New(fmt.Sprint(values...)))
}

func (c *GenerateContext) BindType(ns *g.Ident, name string, t g.Type) g.Type {
	key := Name{ns, name}

	_, exists := c.types[key]
	if exists {
		// c.Error(name, "is already defined")
	} else {
		c.types[key] = t
	}

	return t
}

func (c *GenerateContext) Define(ns *g.Ident, name string, conv g.NameConvention, fallbackSuffixes ...string) *g.Ident {
	currentName := name
	alias := ""

	var suffix int

	for {
		key := Name{ns, currentName}

		_, exists := c.names[key]

		if !exists {
			binding := &Binding{&g.Ident{Name: currentName, NameConvention: conv, SourceAlias: alias}, true}
			c.names[key] = binding
			return binding.ident
		}

		if suffix >= len(fallbackSuffixes) {
			c.Error(name, " is already defined")
		}

		alias = name
		currentName = name + fallbackSuffixes[suffix]
		suffix++
	}
}

func (c *GenerateContext) DefRef(ns *g.Ident, name string, conv g.NameConvention) *g.Ident {
	key := Name{ns, name}

	binding, exists := c.names[key]
	if exists {
		return binding.ident
	} else {
		binding = &Binding{&g.Ident{Name: name, NameConvention: conv}, true}
		c.names[key] = binding
	}

	return binding.ident
}

func (c *GenerateContext) Ref(ns *g.Ident, name string) *g.Ident {
	key := Name{ns, name}

	binding, exists := c.names[key]
	if exists {
		return binding.ident
	} else {
		binding = &Binding{&g.Ident{Name: name}, false}
		c.names[key] = binding
	}

	return binding.ident
}

func (c *GenerateContext) Overload(ns *g.Ident, name string, types []g.Type) *g.Ident {
	currentName := name
	pc := &InlinePrintContext{}
	var idx int

	for {
		pc.Reset()
		pc.WriteString(currentName)
		g.Space.Print(pc)
		g.Join(pc, types, g.Space, g.Comma)
		exactKey := Name{ns, pc.String()}
		nameKey := Name{ns, currentName}

		_, exact := c.names[exactKey]

		if exact {
			return nil
		}

		_, collision := c.names[nameKey]

		if !collision {
			id := &g.Ident{Name: currentName, SourceAlias: name}
			binding := &Binding{id, true}
			c.names[exactKey] = binding
			c.names[nameKey] = binding
			return id
		}

		currentName = name + strconv.Itoa(idx)
		idx++
	}
}

func (c *GenerateContext) RefType(ns *g.Ident, name string) g.Type {
	key := Name{ns, name}

	ref, exists := c.types[key]
	if exists {
		return ref
	} else {
		ref = &TypeRef{ident: &g.Ident{Name: name}}
		c.types[key] = ref
	}

	return ref
}

func isNullTypeExpr(expr ast.TypeExpr) bool {
	id, ok := ptr(expr).(*ast.TypeIdentExpr)
	return ok && id.Name.Name.Value == JSNull
}

func (c *GenerateContext) RefTypeExpr(expr ast.TypeExpr, nameHint string, generic *ast.GenericDecl) g.Type {
	switch t := ptr(expr).(type) {
	case *ast.StringLiteralTypeExpr:
		return c.RefEnumMember(nameHint, t.Type.Value)
	case *ast.TypeIdentExpr:

		if generic != nil {
			for _, arg := range generic.Args {
				if arg.Ident.Name.Value == t.Name.Name.Value {
					if arg.Constraint == nil {
						return c.RefType(nil, JSAny)
					}

					return c.RefTypeExpr(arg.Constraint.Type, nameHint, generic)
				}
			}
		}

		return c.RefType(nil, t.Name.Name.Value)
	case *ast.TypeArrayExpr:
		slice := &g.Slice{Element: c.RefTypeExpr(t.Element, nameHint, generic)}
		return c.BindType(c.module.Ident, ToName(slice), slice)
	case *ast.TypeUnionExpr:
		if len(t.Types) != 2 {
			break
		}

		firstNull := isNullTypeExpr(t.Types[0])
		secondNull := isNullTypeExpr(t.Types[1])

		var baseType g.Type

		switch {
		case firstNull == secondNull:
		case !firstNull:
			baseType = c.RefTypeExpr(t.Types[0], nameHint, generic)
		case !secondNull:
			baseType = c.RefTypeExpr(t.Types[1], nameHint, generic)
		}

		if baseType == nil {
			break
		}

		// TODO: Wrap in a nillable type ref
		return baseType
	}

	return c.RefType(nil, JSAny)
}

func (c *GenerateContext) GetType(ns *g.Ident, name string) g.Type {
	return c.types[Name{ns, name}]
}

func (c *GenerateContext) RefEnumMember(enumName, member string) g.Type {
	if enumName == "" {
		enumName = member + "Type"
	}

	existing := c.GetType(c.module.Ident, enumName)

	if existing != nil {
		return existing
	}

	var enum *g.Enum
	for i, e := range c.module.Enums {
		if e.Ident.Name == enumName {
			enum = c.module.Enums[i]
			break
		}
	}
	if enum == nil {
		c.module.Enums = append(c.module.Enums, &g.Enum{Ident: c.DefRef(c.module.Ident, enumName, g.NameExport)})
		enum = c.module.Enums[len(c.module.Enums)-1]
	}
	enum.Members = append(enum.Members, &g.EnumMember{Ident: c.DefRef(c.module.Ident, member, g.NameExport), Value: member})

	return c.RefType(nil, enumName)
}

func (c *GenerateContext) DefineStruct(iface *ast.InterfaceDecl) {
	node := g.Struct{
		Ident: c.Define(c.module.Ident, iface.Ident.Name.Value, g.NameExport, "Type"),
	}
	c.module.Structs = append(c.module.Structs, &node)

	if iface.Extends != nil {
		for _, parent := range iface.Extends.Interfaces {
			node.Extends = append(node.Extends, c.RefTypeExpr(parent, "", iface.Generic))
		}
	}

	for _, member := range iface.Body.Members {
		switch mem := ptr(member).(type) {
		case *ast.VariableDecl:
			node.Fields = append(node.Fields, &g.Field{
				Ident: c.Define(node.Ident, mem.Ident.Ident(), g.NameExport),
				Type:  c.RefTypeExpr(mem.Type.Type, "", iface.Generic),
			})
		case *ast.MethodDecl:
			var params []*g.Param
			var overloadTypeArgs []g.Type

			for _, param := range mem.Params {
				t := c.RefTypeExpr(param.Type.Type, param.Ident.Name.Value, mem.Generic)
				params = append(params, &g.Param{
					Ident: &g.Ident{Name: param.Ident.Name.Value, NameConvention: g.NameParameter},
					Type:  t,
				})
				overloadTypeArgs = append(overloadTypeArgs, t)
			}

			id := c.Overload(node.Ident, mem.Ident.Name.Value, overloadTypeArgs)

			if id == nil {
				// Resolved overload already exists
				continue
			}

			m := g.Method{
				Ident:  id,
				Recv:   &TypeRef{ident: node.Ident},
				Type:   c.RefTypeExpr(mem.Return.Type, "", mem.Generic),
				Params: params,
			}
			node.Methods = append(node.Methods, &m)
		}
	}
}

func ptr(val any) any {
	rval := reflect.ValueOf(val)
	if rval.Kind() == reflect.Struct {
		ptr := reflect.New(rval.Type())
		ptr.Elem().Set(rval)
		return ptr.Interface()
	}
	return val
}

func Generate(w io.Writer, results []Result, names ...string) {
	c := &GenerateContext{
		types: make(map[Name]g.Type),
		names: make(map[Name]*Binding),
	}
	c.module.Ident = &g.Ident{Name: "bindings", NameConvention: g.NameModule}

	c.BindType(Global, JSVoid, &BuiltInType{g.S("")})
	c.BindType(Global, JSAny, &BuiltInType{g.KAny})
	c.BindType(Global, JSString, &BuiltInType{g.KString})
	c.BindType(Global, JSNumber, &BuiltInType{g.KFloat64})
	c.BindType(Global, JSBoolean, &BuiltInType{g.KBool})
	c.BindType(Global, JSFloat32List, &BuiltInType{g.S("[]float32")})
	c.BindType(Global, JSInt32List, &BuiltInType{g.S("[]int32")})
	c.BindType(c.module.Ident, "namespaceURI", &BuiltInType{g.KString})

	c.module.Vars = append(c.module.Vars, &g.Var{
		Ident: c.Define(c.module.Ident, "Document", g.NameExport),
		Type:  c.RefType(c.module.Ident, "Document"),
	})

	for _, result := range results {
		for _, declIface := range result.Value.Decls {

			switch decl := ptr(declIface).(type) {
			case (*ast.InterfaceDecl):
				switch decl.Ident.Name.Value {
				case "WebGLRenderingContextBase", "WebGLRenderingContext", "HTMLCanvasElement", "WebGLRenderingContextOverloads", "Document":
					c.DefineStruct(decl)
				}
			case (*ast.TypeDecl):
				switch decl.Ident.Name.Value {
				case "GLenum", "GLsizei", "GLuint", "GLfloat", "GLint", "GLboolean":
					c.module.Defs = append(c.module.Defs, &g.Def{
						Ident:  c.Define(c.module.Ident, decl.Ident.Name.Value, g.NameExport),
						Target: c.RefTypeExpr(decl.Value, "", nil),
						Alias:  true,
					})
				}
			}
		}
	}

BindLoop:
	for _, t := range c.types {
		ref, ok := t.(*TypeRef)

		if !ok || ref.binding != nil {
			continue
		}

		// TODO: bind using stronger names
		for _, s := range c.module.Structs {
			if s.Ident.Name == ref.ident.Name {
				ref.binding = s
				continue BindLoop
			}
		}

		for _, e := range c.module.Enums {
			if e.Ident.Name == ref.ident.Name {
				ref.binding = e
				continue BindLoop
			}
		}

		for _, d := range c.module.Defs {
			if d.Ident.Name == ref.ident.Name {
				ref.binding = d
				continue BindLoop
			}
		}

		c.module.Defs = append(c.module.Defs, &g.Def{Ident: ref.ident, Target: &BuiltInType{Symbol: g.S("struct {}")}})
	}

	c.module.Print(&WriterPrintContext{w, 0})
}

func ToName(parts ...g.Printer) string {
	pc := &InlinePrintContext{}
	g.Join(pc, parts)
	return pc.String()
}

type InlinePrintContext struct {
	strings.Builder
}

func (InlinePrintContext) Indent()    {}
func (InlinePrintContext) Unindent()  {}
func (InlinePrintContext) Depth() int { return 0 }

type WriterPrintContext struct {
	io.Writer
	depth int
}

func (pc *WriterPrintContext) WriteString(s string) (n int, err error) {
	return io.WriteString(pc.Writer, s)
}

func (pc *WriterPrintContext) Indent() {
	pc.depth++
}

func (pc *WriterPrintContext) Unindent() {
	pc.depth--
}

func (pc *WriterPrintContext) Depth() int {
	return pc.depth
}
