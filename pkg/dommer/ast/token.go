package ast

import (
	"fmt"
	"strings"
)

type TokenClass int

const (
	None TokenClass = iota
	Symbol
	Whitespace
	Word
	Keyword
	ClsString
	BlockComment
	LineComment
	EOF
)

func (c TokenClass) String() string {
	return [...]string{
		"None",
		"Symbol",
		"Whitespace",
		"Word",
		"Keyword",
		"EOF",
	}[c]
}

type Token struct {
	Value string
	Class TokenClass
}

func T(value string, class TokenClass) Token {
	return Token{
		Value: value,
		Class: class,
	}
}

var (
	Readonly  = Token{"readonly", Keyword}
	Interface = Token{"interface", Keyword}
	Var       = Token{"var", Keyword}
	Declare   = Token{"declare", Keyword}
	Extends   = Token{"extends", Keyword}
	Keyof     = Token{"keyof", Keyword}
	Typeof    = Token{"typeof", Keyword}
	Type      = Token{"type", Keyword}
	Function  = Token{"function", Keyword}
	Namespace = Token{"namespace", Keyword}
	Const     = Token{"const", Keyword}
	Any       = Token{"any", Keyword}
	String    = Token{"string", Keyword}
	Null      = Token{"null", Keyword}
	Void      = Token{"void", Keyword}
	Number    = Token{"number", Keyword}
	Boolean   = Token{"boolean", Keyword}
)

var (
	Lparen       = Token{"(", Symbol}
	Rparen       = Token{")", Symbol}
	Comma        = Token{",", Symbol}
	Dot          = Token{".", Symbol}
	Semi         = Token{";", Symbol}
	Colon        = Token{":", Symbol}
	Pipe         = Token{"|", Symbol}
	Ampersand    = Token{"&", Symbol}
	Lbrc         = Token{"{", Symbol}
	Rbrc         = Token{"}", Symbol}
	Lbrkt        = Token{"[", Symbol}
	Rbrkt        = Token{"]", Symbol}
	Eq           = Token{"=", Symbol}
	SglQte       = Token{"'", Symbol}
	DblQte       = Token{`"`, Symbol}
	Backtick     = Token{"`", Symbol}
	Star         = Token{"*", Symbol}
	QuestionMark = Token{"?", Symbol}
	Lt           = Token{"<", Symbol}
	Gt           = Token{">", Symbol}
)

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

func (t Token) String() string {
	if t.Class == Whitespace {
		return fmt.Sprintf(`%s%#v%s`, GrayBg, t.Value, Reset)
	}

	if t.Class == Keyword {
		return fmt.Sprintf("%s%s%s", Cyan, t.Value, Reset)
	}

	if t.Class == Symbol {
		return fmt.Sprintf("%s%s%s", Yellow, t.Value, Reset)
	}

	if t.Class == EOF {
		return "EOF"
	}

	return fmt.Sprintf("%s%s%s", Purple, t.Value, Reset)
}

func ToksString(toks []Token) string {
	var parts []string

	for _, tok := range toks {
		parts = append(parts, tok.String())
	}

	return strings.Join(parts, " ")
}
