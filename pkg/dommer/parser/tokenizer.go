package parser

import (
	"errors"
	"io"

	"github.com/rryanrussell/goscript/pkg/dommer/ast"
)

type Tokenizer interface {
	Next() (ast.Token, error)
	Position() int
	SetPosition(pos int)
	ReadPosition() int64
}

type ByteScannerTokenizer struct {
	data        []byte
	token       []byte
	tokens      []ast.Token
	dataOffset  int64
	tokenOffset int
}

func NewTokenizer(data []byte) Tokenizer {
	return &ByteScannerTokenizer{
		data: data,
	}
}

var ErrNope = errors.New("nope")

func (t *ByteScannerTokenizer) tryScanLc() (ast.Token, error) {
	var state int

	t.token = t.token[:0]

	for {
		b, err := t.readByte()
		t.dataOffset++

		switch {
		case state < 2 && b == '/':
			state++
		case state < 2:
			return ast.Token{}, ErrNope
		case b == '\n', errors.Is(err, io.EOF):
			return ast.Token{Value: string(t.token), Class: ast.LineComment}, err
		case err != nil:
			return ast.Token{}, err
		default:
			t.token = append(t.token, b)
		}
	}
}

func (t *ByteScannerTokenizer) tryScanBc() (ast.Token, error) {
	var state int

	t.token = t.token[:0]

	for {
		b, err := t.readByte()
		t.dataOffset++

		switch {
		case state == 0 && b == '/':
			state++
		case state == 0:
			return ast.Token{}, ErrNope
		case state == 1 && b == '*':
			state++
		case state == 1:
			return ast.Token{}, ErrNope
		case state == 2 && b == '*':
			state++
		case state == 3 && b == '/':
			return ast.Token{Value: string(t.token), Class: ast.BlockComment}, err
		case state == 3:
			t.token = append(t.token, '*', b)
			state = 2
		case err != nil:
			return ast.Token{}, err
		default:
			t.token = append(t.token, b)
		}
	}
}

func (t *ByteScannerTokenizer) tryScanString() (ast.Token, error) {
	var state int

	t.token = t.token[:0]

	for {
		b, err := t.readByte()
		t.dataOffset++

		switch {
		case state == 0 && b == '"':
			state++
		case state == 0:
			return ast.Token{}, ErrNope
		case state == 1 && b == '"':
			return ast.Token{Value: string(t.token), Class: ast.ClsString}, err
		default:
			t.token = append(t.token, b)
		}
	}
}

func (t *ByteScannerTokenizer) readByte() (byte, error) {
	if t.dataOffset >= int64(len(t.data)) {
		return 0, io.EOF
	}

	return t.data[t.dataOffset], nil
}

func (t *ByteScannerTokenizer) scan() (ast.Token, error) {
	start := t.dataOffset

	lc, err := t.tryScanLc()
	if errors.Is(err, ErrNope) {
		t.dataOffset = start
	} else {
		return lc, err
	}

	bc, err := t.tryScanBc()
	if errors.Is(err, ErrNope) {
		t.dataOffset = start
	} else {
		return bc, err
	}

	str, err := t.tryScanString()
	if errors.Is(err, ErrNope) {
		t.dataOffset = start
	} else {
		return str, err
	}

	var cls ast.TokenClass
	t.token = t.token[:0]

	for {
		char, err := t.readByte()

		if err != nil {
			return ast.T(string(t.token), cls), err
		}

		t.dataOffset++

		var charClass ast.TokenClass

		switch string(char) {
		case " ", "\n", "\t":
			charClass = ast.Whitespace
		case "(", ")", ",", ".", ";", ":", "{", "}", "[", "]", "=", "'", `"`, "`", "*", "?", "<", ">", "|", "&", "/":
			charClass = ast.Symbol
		default:
			charClass = ast.Word
		}

		switch {
		case cls == charClass:
			t.token = append(t.token, char)
		case cls == ast.None && charClass == ast.Symbol:
			return ast.T(string(char), ast.Symbol), nil
		case cls == ast.None:
			cls = charClass
			t.token = append(t.token, char)
		default:
			tok := ast.T(string(t.token), cls)
			t.dataOffset--
			return tok, nil
		}
	}
}

func (t *ByteScannerTokenizer) Position() int {
	return t.tokenOffset
}

func (t *ByteScannerTokenizer) ReadPosition() int64 {
	return t.dataOffset
}

func (t *ByteScannerTokenizer) SetPosition(pos int) {
	t.tokenOffset = pos
}

func (t *ByteScannerTokenizer) Next() (ast.Token, error) {
	if t.tokenOffset < len(t.tokens) {
		tok := t.tokens[t.tokenOffset]
		t.tokenOffset++
		return tok, nil
	}

	token, err := t.scan()

	if err != nil {
		t.tokens = append(t.tokens, token)
		t.tokenOffset++
		return token, err
	}

	if token.Class != ast.Word {
		t.tokens = append(t.tokens, token)
		t.tokenOffset++
		return token, nil
	}

	switch token.Value {
	case ast.Readonly.Value,
		ast.Interface.Value,
		ast.Var.Value,
		ast.Declare.Value,
		ast.Extends.Value,
		ast.Keyof.Value,
		ast.Typeof.Value,
		ast.Type.Value,
		ast.Namespace.Value,
		ast.Function.Value,
		ast.Const.Value:
		token.Class = ast.Keyword
	}

	t.tokens = append(t.tokens, token)
	t.tokenOffset++
	return token, nil
}
