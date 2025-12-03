package generator

import (
	"go/token"
	"strings"
)

type Token token.Token

func (t Token) Print(c PrintContext) {
	c.WriteString(token.Token(t).String())
}

func T(t token.Token) Token {
	return Token(t)
}

type Symbol string

func (s Symbol) Print(c PrintContext) {
	c.WriteString(string(s))
}

func S(s string) Symbol {
	return Symbol(s)
}

type IndentType struct{}

func (IndentType) Print(c PrintContext) {
	c.WriteString(strings.Repeat(string(IndentUnit), c.Depth()))
}

var (
	Lbrc   = T(token.LBRACE)
	Rbrc   = T(token.RBRACE)
	Lparen = T(token.LPAREN)
	Rparen = T(token.RPAREN)
	Lbrkt  = T(token.LBRACK)
	Rbrkt  = T(token.RBRACK)
	Comma  = T(token.COMMA)
	Assign = T(token.ASSIGN)

	Backtick = S("`")
	Space    = S(" ")
	Newline  = S("\n")
	DblQte   = S(`"`)

	IndentUnit = S("	")
	Indent     = IndentType{}

	KType    = S("type")
	KStruct  = S("struct")
	KString  = S("string")
	KVar     = S("var")
	KConst   = S("const")
	KPac     = S("package")
	KAny     = S("any")
	KFunc    = S("func")
	KFloat64 = S("float64")
	KBool    = S("bool")
)
