package parser

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/rryanrussell/goscript/pkg/dommer/ast"
)

func test(c Ctx, node Node) bool {
	_, ok := try(c, node)
	return ok
}

func try(c Ctx, node Node) (any, bool) {
	c.Begin()

	val, ok := node.Process(c)

	if ok {
		c.Accept()
		return val, true
	}

	c.Reject()
	return nil, false
}

func tryFunc(c Ctx, node FuncNode) (any, bool) {
	return try(c, node)
}

func typeOf[T any]() reflect.Type {
	var zero T
	return reflect.TypeOf(&zero).Elem()
}

func popOption(tag reflect.StructTag, key string) reflect.StructTag {
	var remain []string
	for _, option := range strings.Split(string(tag), " ") {
		if !strings.HasPrefix(option, key+":") {
			remain = append(remain, option)
		}
	}
	return reflect.StructTag(strings.Join(remain, " "))
}

type Entry struct {
	Goal    reflect.Type
	Options reflect.StructTag
}

type Registry = map[Entry]Node

func nodeIface(registry Registry, entry Entry) Node {
	if entry.Goal.Kind() != reflect.Interface {
		panic("iface must be an interface")
	}

	return &Parallel{}
}

func nodeOptional(registry Registry, entry Entry) Node {
	return &Optional{nodeType(registry, Entry{
		Goal:    entry.Goal,
		Options: popOption(entry.Options, "optional"),
	})}
}

func nodeToken(registry Registry, entry Entry) Node {
	symbol := entry.Options.Get("symbol")
	if len(symbol) > 0 {
		var token ast.Token

		switch symbol {
		case "QuestionMark":
			token = ast.QuestionMark
		case "Lparen":
			token = ast.Lparen
		case "Rparen":
			token = ast.Rparen
		case "Semi":
			token = ast.Semi
		case "Colon":
			token = ast.Colon
		case "Lbrc":
			token = ast.Lbrc
		case "Rbrc":
			token = ast.Rbrc
		case "Lt":
			token = ast.Lt
		case "Gt":
			token = ast.Gt
		case "Pipe":
			token = ast.Pipe
		case "Rbrkt":
			token = ast.Rbrkt
		case "Lbrkt":
			token = ast.Lbrkt
		case "Eq":
			token = ast.Eq
		case "SglQte":
			token = ast.SglQte
		case "DblQte":
			token = ast.DblQte
		case "Dot":
			token = ast.Dot
		default:
			panic(fmt.Errorf("symbol '%s' value not mapped", symbol))
		}

		return MatchesToken{token}
	}

	switch entry.Options.Get("class") {
	case "Word":
		return TokenValue{ast.Word}
	case "Keyword":
		return TokenValue{ast.Keyword}
	case "KeywordWord":
		return TokenKeywordWord{}
	case "String":
		return TokenValue{ast.ClsString}
	case "":
		break
	default:
		panic("class value not mapped")
	}

	switch entry.Options.Get("keyword") {
	case "declare":
		return MatchesToken{ast.Declare}
	case "var":
		return MatchesToken{ast.Var}
	case "extends":
		return MatchesToken{ast.Extends}
	case "interface":
		return MatchesToken{ast.Interface}
	case "keyof":
		return MatchesToken{ast.Keyof}
	case "typeof":
		return MatchesToken{ast.Typeof}
	case "readonly":
		return MatchesToken{ast.Readonly}
	case "type":
		return MatchesToken{ast.Type}
	case "function":
		return MatchesToken{ast.Function}
	case "namespace":
		return MatchesToken{ast.Namespace}
	case "const":
		return MatchesToken{ast.Const}
	case "":
		break
	default:
		panic("keyword value not mapped")
	}

	switch entry.Options.Get("word") {
	case "":
		break
	default:
		return MatchesToken{Token: ast.Token{Value: entry.Options.Get("word"), Class: ast.Word}}
	}

	panic(fmt.Errorf("unexpected token options: %s", entry.Options))
}

func bindStructs(registry Registry) {
	for entry, node := range registry {
		automaton, ok := node.(*Automaton)

		if !ok {
			continue
		}

		for i := 0; i < entry.Goal.NumField(); i++ {
			field := entry.Goal.Field(i)
			if field.Tag.Get("ignore") == "true" {
				continue
			}
			automaton.Sequence = append(automaton.Sequence, Step{
				Target: field,
				Node:   nodeType(registry, Entry{Goal: field.Type, Options: field.Tag}),
			})
		}
	}
}

func nodeStruct(registry Registry, entry Entry) Node {
	if entry.Goal.Kind() != reflect.Struct {
		panic("goal must be a struct")
	}

	node := &Automaton{
		Goal: entry.Goal,
	}

	registry[entry] = node

	for i := 0; i < entry.Goal.NumField(); i++ {
		field := entry.Goal.Field(i)
		if field.Tag.Get("ignore") == "true" {
			continue
		}
		nodeType(registry, Entry{Goal: field.Type, Options: field.Tag})
	}

	return node
}

func nodeType(registry Registry, entry Entry) Node {
	node, ok := registry[entry]

	if ok {
		return node
	}

	if entry.Options.Get("optional") == "true" {
		return nodeOptional(registry, entry)
	}

	switch entry.Goal.Kind() {
	case reflect.Pointer:
		return nodeType(registry, Entry{entry.Goal.Elem(), entry.Options})
	case reflect.Interface:
		node = nodeIface(registry, entry)
	case reflect.Slice:
		node = nodeSlice(registry, entry)
	case reflect.Struct:
		if entry.Goal == typeOf[ast.TokenExpr]() {
			node = nodeToken(registry, entry)
		} else {
			if entry.Options != "" {
				panic("unexpected options")
			}
			return nodeStruct(registry, entry)
		}
	default:
		panic("unexpected ast type")
	}

	registry[entry] = node

	return node
}

func nodeSlice(r Registry, entry Entry) Node {
	var delim Node
	trailing := entry.Options.Get("trailing") == "true"
	var minimum int

	minStr := entry.Options.Get("min")

	if len(minStr) > 0 {
		parsed, err := strconv.Atoi(minStr)
		if err != nil {
			panic(err)
		}
		minimum = parsed
	}

	switch entry.Options.Get("delim") {
	case "Comma":
		delim = MatchesToken{ast.Comma}
	case "Pipe":
		delim = MatchesToken{ast.Pipe}
	case "Ampersand":
		delim = MatchesToken{ast.Ampersand}
	case "":
		break
	default:
		panic("delim value not mapped")
	}

	return &Repeat{
		PayNode:             nodeType(r, Entry{Goal: entry.Goal.Elem()}),
		Delim:               delim,
		AcceptTrailingDelim: trailing,
		Minimum:             minimum,
	}
}

func NewRegistry() Registry {
	return make(map[Entry]Node)
}

func bindIface(registry Registry, iface reflect.Type, impls ...reflect.Type) {
	node := nodeType(registry, Entry{Goal: iface, Options: ""}).(*Parallel)

	for _, impl := range impls {
		if !impl.Implements(iface) {
			panic("impl does not implement iface")
		}
		implNode := nodeType(registry, Entry{Goal: impl, Options: ""})
		node.Paths = append(node.Paths, implNode)
	}
}

func NewParser() Node {
	reg := NewRegistry()

	root := nodeType(reg, Entry{Goal: typeOf[ast.DtsModule](), Options: ""})

	bindIface(reg, typeOf[ast.InterfaceExpr](),
		typeOf[ast.GenericTypeExpr](),
		typeOf[ast.TypeIdentExpr](),
	)

	bindIface(reg, typeOf[ast.TypeExpr](),
		typeOf[ast.TypeBodyExpr](),
		typeOf[ast.TupleExpr](),
		typeOf[ast.TypeUnionExpr](),
		typeOf[ast.LambdaTypeExpr](),
		typeOf[ast.GenericTypeExpr](),
		typeOf[ast.TypeIntersectionExpr](),
		typeOf[ast.KeyofTypeExpr](),
		typeOf[ast.TypeofTypeExpr](),
		typeOf[ast.TypeArrayExpr](),
		typeOf[ast.ParenthesesTypeExpr](),
		typeOf[ast.TypeIndex](),
		typeOf[ast.StringLiteralTypeExpr](),
		typeOf[ast.TypeIdentExpr](),
	)

	bindIface(reg, typeOf[ast.TypeBodyDecl](),
		typeOf[ast.VariableDecl](),
		typeOf[ast.SetterDecl](),
		typeOf[ast.GetterDecl](),
		typeOf[ast.MethodDecl](),
		typeOf[ast.CallableDecl](),
		typeOf[ast.IndexerDecl](),
	)

	bindIface(reg, typeOf[ast.DtsDecl](),
		typeOf[ast.InterfaceDecl](),
		typeOf[ast.GlobalConstDecl](),
		typeOf[ast.GlobalVarDecl](),
		typeOf[ast.TypeDecl](),
		typeOf[ast.NamespaceDecl](),
		typeOf[ast.FunctionDecl](),
	)

	bindIface(reg, typeOf[ast.VariableIdent](),
		typeOf[ast.StringIdent](),
		typeOf[ast.KeywordIdent](),
	)

	bindStructs(reg)

	return root
}
