package js

import "strings"

type Token string

func (t Token) Print(c PrintContext) {
	c.WriteString(string(t))
}

func (t Token) Walk(v Visitor) {}

const (
	Space    Token      = " "
	Dot      Token      = "."
	Comma    Token      = ","
	Colon    Token      = ":"
	Eq       Token      = "="
	If       Token      = "if"
	Else     Token      = "else"
	For      Token      = "for"
	While    Token      = "while"
	Let      Token      = "let"
	Const    Token      = "const"
	Ret      Token      = "return"
	Function Token      = "function"
	In       Token      = "in"
	Of       Token      = "of"
	QstnMrk  Token      = "?"
	HuhHuh   Token      = "??"
	This     Token      = "this"
	KClass   Token      = "class"
	KNew     Token      = "new"
	Lbrkt    Token      = "["
	Rbrkt    Token      = "]"
	Lbrc     Token      = "{"
	Rbrc     Token      = "}"
	Lparen   Token      = "("
	Rparen   Token      = ")"
	Indent   IndentType = "    "
	Semi     Token      = ";"
	Inc      Token      = "++"
	Dec      Token      = "--"
	Newline  Token      = "\n"
	Switch   Token      = "switch"
	Case     Token      = "case"
	Default  Token      = "default"
	Break    Token      = "break"
	Continue Token      = "continue"
)

type IndentType string

func (unit IndentType) Print(c PrintContext) {
	c.WriteString(strings.Repeat(string(unit), c.Depth()))
}
