package transpiler

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/rryanrussell/goscript/js"
)

type Context = ModuleContext

func exprs(ctx *Context, in []ast.Expr) []js.Expr {
	var out []js.Expr
	for _, x := range in {
		out = append(out, expr(ctx, x))
	}
	return out
}

func TypeExpr(ctx Complainer, e ast.Expr) string {
	switch x := e.(type) {
	case *ast.BasicLit:
		return x.Value
	case *ast.Ident:
		return x.Name
	case *ast.StarExpr:
		return TypeExpr(ctx, x.X)
	case *ast.ArrayType:
		return TypeExpr(ctx, x.Elt) + "[]"
	case *ast.StructType:
		return "any"
	}

	if ctx != nil {
		ctx.Warn(e, "could not create type name")
	}

	return ""
}

type Token struct {
	Node ast.Node
	Repr string
}

func (t Token) String() string {
	return t.Repr
}

func unknown(ctx *Context, n ast.Node) js.Expr {
	ctx.Error(n, fmt.Sprintf("unknown node %T", n))
	return js.Ident("__U__")
}

func ident(ctx *Context, i *ast.Ident) *js.Ident {
	val := js.Ident(i.Name)
	return &val
}

func selector(ctx *Context, s *ast.SelectorExpr) *js.Selector {
	return &js.Selector{
		X:   expr(ctx, s.X),
		Sel: ident(ctx, s.Sel),
	}
}

func index(ctx *Context, s *ast.IndexExpr) *js.Index {
	return &js.Index{
		X:     expr(ctx, s.X),
		Index: expr(ctx, s.Index),
	}
}

func sliceExpr(ctx *Context, s *ast.SliceExpr) js.Expr {
	var args []js.Expr

	if s.Low != nil {
		args = append(args, expr(ctx, s.Low))
	} else if s.High != nil {
		args = append(args, &js.BasicLit{Value: "0"})
	}

	if s.High != nil {
		args = append(args, expr(ctx, s.High))
	}

	return &js.Call{
		Func: &js.Selector{
			X:   expr(ctx, s.X),
			Sel: js.IdentP("slice"),
		},
		Args: args,
	}
}

func call(ctx *Context, c *ast.CallExpr) js.Expr {
	id, ok := c.Fun.(*ast.Ident)
	if ok && id.Name == "len" && len(c.Args) == 1 {
		return &js.Parens{
			X: &js.Binary{
				X: &js.Selector{
					X:             expr(ctx, c.Args[0]),
					Sel:           js.IdentP("length"),
					OptionalChain: true,
				},
				Op: js.HuhHuh,
				Y:  &js.BasicLit{Value: "0"},
			},
		}
	}

	if ok && id.Name == "append" && len(c.Args) == 2 {
		return &js.ArrayLit{
			Elts: []js.Expr{
				&js.Spread{
					X: &js.Parens{
						X: &js.Binary{
							X:  expr(ctx, c.Args[0]),
							Op: js.HuhHuh,
							Y:  &js.ArrayLit{},
						},
					},
				},
				expr(ctx, c.Args[1]),
			},
		}
	}

	return &js.Call{
		Func: expr(ctx, c.Fun),
		Args: exprs(ctx, c.Args),
	}
}

func composite(ctx *Context, c *ast.CompositeLit) js.Expr {
	switch c.Type.(type) {
	case *ast.ArrayType:
		return &js.ArrayLit{
			Elts: exprs(ctx, c.Elts),
		}
	}

	var lit Instance
	lit.Type = TypeExpr(ctx, c.Type)
	var class *js.Class

	for _, decl := range ctx.Module.Decls {
		cls, ok := decl.(*js.Class)
		if ok && cls.NameEq(lit.Type) {
			class = cls
			lit.Class = true
			break
		}
	}

	for i, elt := range c.Elts {
		var init InitValue
		kvp, ok := elt.(*ast.KeyValueExpr)
		if ok {
			init.Member = kvp.Key.(*ast.Ident).Name
			init.Value = expr(ctx, kvp.Value)
		} else {
			if class != nil {
				init.Member = class.Fields[i].NameString()
			}

			init.Value = expr(ctx, elt)
		}
		lit.Args = append(lit.Args, &init)
	}

	return &js.GuestExpr{Node: &lit}
}

func op(ctx *Context, tok token.Token) js.Token {
	return js.Token(tok.String())
}

func binary(ctx *Context, b *ast.BinaryExpr) *js.Binary {
	return &js.Binary{
		X:  expr(ctx, b.X),
		Op: op(ctx, b.Op),
		Y:  expr(ctx, b.Y),
	}
}

func ref(ctx *Context, i *ast.Ident) *js.Ref {
	return &js.Ref{Named: js.Named{Name: (*js.Ident)(&i.Name)}}
}

func expr(ctx *Context, w ast.Expr) js.Expr {
	switch x := w.(type) {
	case *ast.CallExpr:
		return call(ctx, x)
	case *ast.StarExpr:
		asExpr := expr(ctx, x.X)
		if ctx.Semantics == nil {
			ctx.Semantics = make(map[js.Node]Semantics)
		}
		sem := ctx.Semantics[asExpr]
		sem.PassByReference = true
		ctx.Semantics[asExpr] = sem
		return asExpr
	case *ast.SelectorExpr:
		return selector(ctx, x)
	case *ast.BasicLit:
		return &js.BasicLit{Value: x.Value}
	case *ast.Ident:
		return ref(ctx, x)
	case *ast.UnaryExpr:
		if x.Op == token.AND {
			asExpr := expr(ctx, x.X)
			if ctx.Semantics == nil {
				ctx.Semantics = make(map[js.Node]Semantics)
			}
			sem := ctx.Semantics[asExpr]
			sem.PassByReference = true
			ctx.Semantics[asExpr] = sem
			return asExpr
		}

		return &js.Unary{
			X:  expr(ctx, x.X),
			Op: op(ctx, x.Op),
		}
	case *ast.IndexExpr:
		return index(ctx, x)
	case *ast.SliceExpr:
		return sliceExpr(ctx, x)
	case *ast.CompositeLit:
		return composite(ctx, x)
	case *ast.BinaryExpr:
		return binary(ctx, x)
	case *ast.FuncLit:
		return funlit(ctx, x)
	case *ast.TypeAssertExpr:
		asExpr := expr(ctx, x.X)
		if ctx.Semantics == nil {
			ctx.Semantics = make(map[js.Node]Semantics)
		}
		sem := ctx.Semantics[asExpr]
		sem.TypeAssertion = TypeExpr(ctx, x.Type)

		ctx.Semantics[asExpr] = sem
		return asExpr
	default:
		return unknown(ctx, w)
	}
}

func block(ctx *Context, b *ast.BlockStmt) *js.Block {
	block := &js.Block{}
	for _, stmt := range b.List {
		block.Lines = append(block.Lines, statement(ctx, stmt))
	}
	return block
}

func ifs(ctx *Context, i *ast.IfStmt) js.Stmt {
	var els js.Stmt
	if i.Else != nil {
		els = statement(ctx, i.Else)
	}

	ifStmt := &js.IfStmt{
		Cond: expr(ctx, i.Cond),
		Body: block(ctx, i.Body),
		Else: els,
	}

	if i.Init != nil {
		block := &js.Block{}
		block.Lines = append(block.Lines, statement(ctx, i.Init))
		block.Lines = append(block.Lines, ifStmt)
		return block
	}

	return ifStmt
}

func switches(ctx *Context, s *ast.SwitchStmt) js.Stmt {
	// Handling Switch
	sw := &js.SwitchStmt{}
	if s.Tag != nil {
		sw.Tag = expr(ctx, s.Tag)
	}

	sw.Body = &js.CaseBlock{}

	for _, stmt := range s.Body.List {
		cc, ok := stmt.(*ast.CaseClause)
		if !ok {
			continue
		}

		jsCc := &js.CaseClause{}
		jsCc.List = exprs(ctx, cc.List)

		// Handle Body
		hasFallthrough := false
		for i, st := range cc.Body {
			// Check for fallthrough (must be last statement)
			if branch, ok := st.(*ast.BranchStmt); ok && branch.Tok == token.FALLTHROUGH {
				if i == len(cc.Body)-1 {
					hasFallthrough = true
					// Don't emit fallthrough statement in JS
					continue
				}
			}
			jsCc.Body = append(jsCc.Body, statement(ctx, st))
		}

		// Add break if no fallthrough
		if !hasFallthrough {
			jsCc.Body = append(jsCc.Body, &js.TokenStmt{Token: js.Break})
		}

		sw.Body.List = append(sw.Body.List, jsCc)
	}

	if s.Init != nil {
		block := &js.Block{}
		block.Lines = append(block.Lines, statement(ctx, s.Init))
		block.Lines = append(block.Lines, sw)
		return block
	}

	return sw
}

func fors(ctx *Context, f *ast.ForStmt) js.Stmt {
	if f.Post == nil && f.Init == nil {
		return &js.WhileStmt{
			Cond: expr(ctx, f.Cond),
			Body: block(ctx, f.Body),
		}
	} else {
		return &js.ForStmt{
			Init: statement(ctx, f.Init),
			Cond: expr(ctx, f.Cond),
			Post: statement(ctx, f.Post),
			Body: block(ctx, f.Body),
		}
	}
}

func ranges(ctx *Context, r *ast.RangeStmt) js.Stmt {
	var loop js.ForEach
	loop.Body = block(ctx, r.Body)
	loop.Iterable = expr(ctx, r.X)
	var key, value string

	k, ok := r.Key.(*ast.Ident)
	if ok && k.Name != "_" {
		key = k.Name
	}

	v, ok := r.Value.(*ast.Ident)
	if ok && v.Name != "_" {
		value = v.Name
	}

	switch {
	case key != "" && value != "":
		return &js.ExprStmt{X: unknown(ctx, r)}
	case key == "" && value == "":
		return &js.ExprStmt{X: unknown(ctx, r)}
	case key == "":
		loop.Value = &js.VarDecl{
			Var:    &js.Var{Named: js.NewNamed(value)},
			Const:  true,
			Inline: true,
		}
	default:
		loop.Key = &js.VarDecl{
			Var:    &js.Var{Named: js.NewNamed(key)},
			Const:  true,
			Inline: true,
		}
	}

	return &loop
}

func returns(ctx *Context, r *ast.ReturnStmt) *js.Return {
	var ret js.Return

	switch len(r.Results) {
	case 0:
	case 1:
		ret.X = expr(ctx, r.Results[0])
	default:
		x := js.ArrayLit{}

		for _, res := range r.Results {
			x.Elts = append(x.Elts, expr(ctx, res))
		}

		ret.X = &x
	}

	return &ret
}

func incdec(ctx *Context, o *ast.IncDecStmt) js.Stmt {
	var op js.Token

	switch o.Tok {
	case token.INC:
		op = js.Inc
	case token.DEC:
		op = js.Dec
	default:
		ctx.Error(o, fmt.Sprintf("unsupported token %s", o.Tok))
	}

	return &js.ExprStmt{
		X: &js.Unary{
			Op: op,
			X:  expr(ctx, o.X),
		},
	}
}

func statement(ctx *Context, s ast.Stmt) js.Stmt {
	switch stmt := s.(type) {
	case *ast.ExprStmt:
		return &js.ExprStmt{
			X: expr(ctx, stmt.X),
		}
	case *ast.AssignStmt:
		return assign(ctx, stmt)
	case *ast.ReturnStmt:
		return returns(ctx, stmt)
	case *ast.IncDecStmt:
		return incdec(ctx, stmt)
	case *ast.IfStmt:
		return ifs(ctx, stmt)
	case *ast.SwitchStmt:
		return switches(ctx, stmt)
	case *ast.ForStmt:
		return fors(ctx, stmt)
	case *ast.RangeStmt:
		return ranges(ctx, stmt)
	case *ast.BlockStmt:
		return block(ctx, stmt)
	case *ast.DeclStmt:
		decl := stmt.Decl.(*ast.GenDecl)
		if len(decl.Specs) != 1 {
			break
		}
		val := decl.Specs[0].(*ast.ValueSpec)

		if decl.Tok == token.VAR {
			return &js.Assign{
				Define: true,
				Lhs:    single(ctx, val.Names),
				Rhs:    &js.ObjectLit{Type: TypeExpr(ctx, val.Type)},
			}
		}
	}

	return &js.ExprStmt{
		X: unknown(ctx, s),
	}
}

func single[T ast.Expr](ctx *Context, el []T) js.Expr {
	if len(el) > 1 {
		ctx.Fatal(el[1], "did not expect second element")
	}

	return expr(ctx, el[0])
}

func assign(ctx *Context, a *ast.AssignStmt) *js.Assign {
	return &js.Assign{
		Define: a.Tok == token.DEFINE,
		Lhs:    single(ctx, a.Lhs),
		Rhs:    single(ctx, a.Rhs),
	}
}

func fun(ctx *Context, t *ast.FuncType, body *ast.BlockStmt, name string) *js.FuncDecl {
	var f js.FuncDecl

	f.Name = (*js.Ident)(&name)

	for _, param := range t.Params.List {
		for _, name := range param.Names {
			var v js.Var
			v.SetName(js.Ident(name.Name))
			v.Type = TypeExpr(ctx, param.Type)
			f.Args = append(f.Args, &v)
		}
	}

	if ctx.Depth() == 0 && name == "main" {
		ctx.Main()
	}

	f.Body = block(ctx, body)

	return &f
}

func funlit(ctx *Context, f *ast.FuncLit) *js.FuncDecl {
	lit := fun(ctx, f.Type, f.Body, "")
	lit.Body.Inline = true
	return lit
}

func fundcl(ctx *Context, f *ast.FuncDecl) *js.FuncDecl {
	return fun(ctx, f.Type, f.Body, f.Name.Name)
}

func clsdecl(ctx *Context, decl *ast.GenDecl) *js.Class {
	var cls js.Class
	t := decl.Specs[0].(*ast.TypeSpec)
	cls.Name = (*js.Ident)(&t.Name.Name)

	s := t.Type.(*ast.StructType)

	for _, f := range s.Fields.List {
		cls.Fields = append(cls.Fields, &js.Field{
			Type:  TypeExpr(ctx, f.Type),
			Named: js.Named{Name: (*js.Ident)(&f.Names[0].Name)},
		})
	}

	return &cls
}

func Decl(ctx *Context, w ast.Decl) {
	switch d := w.(type) {
	case *ast.FuncDecl:
		ctx.Module.Decls = append(ctx.Module.Decls, fundcl(ctx, d))
	case *ast.GenDecl:
		if d.Tok == token.IMPORT {
			return
		}

		if d.Tok == token.TYPE {
			ctx.Module.Decls = append(ctx.Module.Decls, clsdecl(ctx, d))
			return
		}

		if d.Tok == token.VAR {
			for _, spec := range d.Specs {
				value := spec.(*ast.ValueSpec)
				for _, name := range value.Names {
					nameIdent := js.Ident(name.Name)
					ctx.Module.Decls = append(ctx.Module.Decls, &js.VarDecl{Var: &js.Var{
						Named: js.Named{Name: &nameIdent},
						Type:  TypeExpr(ctx, value.Type),
					}})
				}

			}
			return
		}

		unknown(ctx, d)

	default:
		unknown(ctx, d)
	}
}
