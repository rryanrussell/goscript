package transpiler

import (
	"github.com/rryanrussell/goscript/pkg/js"
)

type InitValue struct {
	Member string
	Value  js.Expr
}

func (val *InitValue) Walk(v js.Visitor) {
	v.Visit(val.Value)
}

func (val *InitValue) Print(c js.PrintContext) {
	js.WriteAll(c, js.Token(val.Member), js.Colon, val.Value)
}

type Instance struct {
	Type  string
	Args  []*InitValue
	Class bool
}

func (i *Instance) Walk(v js.Visitor) {
	js.Stroll(v, i.Args)
}

func (i Instance) Print(c js.PrintContext) {
	if i.Class {
		var args []js.Expr

		for _, arg := range i.Args {
			args = append(args, arg.Value)
		}

		js.New{
			Class: js.IdentP(i.Type),
			Args:  args,
		}.Print(c)
	} else {
		var args []js.Expr

		for _, arg := range i.Args {
			args = append(args, js.GuestExpr{Node: arg})
		}

		js.ObjectLit{
			Type: i.Type,
			Elts: args,
		}.Print(c)
	}
}
