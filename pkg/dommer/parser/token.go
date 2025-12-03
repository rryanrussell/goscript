package parser

import (
	"fmt"

	"github.com/rryanrussell/goscript/pkg/dommer/ast"
)

type MatchesToken struct {
	ast.Token
}

func (t MatchesToken) Process(c Ctx) (any, bool) {
	return tryFunc(c, FuncNode{func() (any, bool) {
		tok := c.Next(false, false)
		match := t.Class == tok.Class && t.Value == tok.Value

		c.Debug().Match(t, tok)

		if !match {
			return nil, false
		} else {
			expr := ast.TokenExpr(tok)
			return &expr, true
		}
	}, fmt.Sprintf("MatchesToken: %s", t.Token)})
}

type TokenValue struct {
	ast.TokenClass
}

func (t TokenValue) Process(c Ctx) (any, bool) {
	return tryFunc(c, FuncNode{func() (any, bool) {
		tok := c.Next(false, false)
		ok := t.TokenClass == tok.Class

		c.Debug().Match(t, tok)

		if ok {
			expr := ast.TokenExpr(tok)
			return &expr, true
		}

		return nil, false
	}, fmt.Sprintf("TokenValue: %s", t.TokenClass)})
}

type TokenKeywordWord struct{}

func (t TokenKeywordWord) Process(c Ctx) (any, bool) {
	return tryFunc(c, FuncNode{func() (any, bool) {
		tok := c.Next(false, false)

		c.Debug().Match(t, tok)

		switch tok.Class {
		case ast.Word, ast.Keyword:
			expr := ast.TokenExpr(tok)
			return &expr, true
		}

		return nil, false
	}, "TokenKeywordWord"})
}
