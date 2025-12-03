package ast

type TypeIdentExpr struct {
	Name *Ident
}

type GenericTypeExpr struct {
	Name        *Ident
	GenericArgs *GenericArgs
}

func (TypeIdentExpr) isInterfaceExpr()   {}
func (GenericTypeExpr) isInterfaceExpr() {}

type InterfaceExpr interface {
	TypeExpr
	isInterfaceExpr()
}

type KeyofTypeExpr struct {
	KeyOf *TokenExpr `keyword:"keyof"`
	Type  TypeExpr
}

type TypeofTypeExpr struct {
	TypeOf *TokenExpr `keyword:"typeof"`
	Type   TypeExpr
}

type TypeUnionExpr struct {
	Types []TypeExpr `delim:"Pipe" min:"2"`
}

type TypeIntersectionExpr struct {
	Types []TypeExpr `delim:"Ampersand" min:"2"`
}

type StringLiteralTypeExpr struct {
	Type *TokenExpr `class:"String"`
}

type ArrayExpr struct {
	Lbrkt *TokenExpr `symbol:"Lbrkt"`
	Rbrkt *TokenExpr `symbol:"Rbrkt"`
}

type TypeArrayExpr struct {
	Element     TypeExpr
	Array       *ArrayExpr
	DoubleArray *ArrayExpr `optional:"true"`
}

type TupleExpr struct {
	Lbrkt *TokenExpr `symbol:"Lbrkt"`
	Types []TypeExpr `delim:"Comma" min:"1"`
	Rbrkt *TokenExpr `symbol:"Rbrkt"`
}

type TypeBodyExpr struct {
	Lbrc    *TokenExpr `symbol:"Lbrc"`
	Members []TypeBodyDecl
	Rbrc    *TokenExpr      `symbol:"Rbrc"`
	Methods []*MethodDecl   `ignore:"true"`
	Attrs   []*VariableDecl `ignore:"true"`
}

type TypeExpr interface {
	isTypeExpr()
}

func (VariableDecl) isTypeBodyDecl() {}
func (MethodDecl) isTypeBodyDecl()   {}
func (IndexerDecl) isTypeBodyDecl()  {}
func (GetterDecl) isTypeBodyDecl()   {}
func (SetterDecl) isTypeBodyDecl()   {}
func (CallableDecl) isTypeBodyDecl() {}

type TypeBodyDecl interface {
	isTypeBodyDecl()
}

func (TypeBodyExpr) isTypeExpr()          {}
func (TypeUnionExpr) isTypeExpr()         {}
func (TypeIdentExpr) isTypeExpr()         {}
func (TypeArrayExpr) isTypeExpr()         {}
func (LambdaTypeExpr) isTypeExpr()        {}
func (ParenthesesTypeExpr) isTypeExpr()   {}
func (TypeIndex) isTypeExpr()             {}
func (StringLiteralTypeExpr) isTypeExpr() {}
func (TypeIntersectionExpr) isTypeExpr()  {}
func (TupleExpr) isTypeExpr()             {}
func (KeyofTypeExpr) isTypeExpr()         {}
func (TypeofTypeExpr) isTypeExpr()        {}
func (GenericTypeExpr) isTypeExpr()       {}

type TypeConstraint struct {
	Extends *TokenExpr `keyword:"extends"`
	Type    TypeExpr
}

type TypeArgDefault struct {
	Eq      *TokenExpr `symbol:"Eq"`
	Default TypeExpr
}

type TypeArg struct {
	Ident      *Ident
	Constraint *TypeConstraint `optional:"true"`
	Default    *TypeArgDefault `optional:"true"`
}

type GenericDecl struct {
	Lt   *TokenExpr `symbol:"Lt"`
	Args []*TypeArg `delim:"Comma"`
	Gt   *TokenExpr `symbol:"Gt"`
}

type LambdaTypeExpr struct {
	Generic *GenericDecl `optional:"true"`
	Lparen  *TokenExpr   `symbol:"Lparen"`
	Params  []*ParamExpr `delim:"Comma"`
	Rparen  *TokenExpr   `symbol:"Rparen"`
	ArrowEq *TokenExpr   `symbol:"Eq"`
	ArrowGt *TokenExpr   `symbol:"Gt"`
	Return  TypeExpr
}

type GenericArgs struct {
	Lt   *TokenExpr `symbol:"Lt"`
	Args []TypeExpr `delim:"Comma" min:"1"`
	Gt   *TokenExpr `symbol:"Gt"`
}

type ParenthesesTypeExpr struct {
	Lparen *TokenExpr `symbol:"Lparen"`
	Expr   TypeExpr
	Rparen *TokenExpr `symbol:"Rparen"`
}

type TypeIndex struct {
	Element TypeExpr
	Lbrkt   *TokenExpr `symbol:"Lbrkt"`
	Args    []*Ident   `delim:"Comma" min:"1"`
	Rbrkt   *TokenExpr `symbol:"Rbrkt"`
}
