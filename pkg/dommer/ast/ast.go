package ast

type Ident struct {
	Name *TokenExpr `class:"Word"`
}

type KeywordIdent struct {
	Name *TokenExpr `class:"KeywordWord"`
}

func (id KeywordIdent) Ident() string {
	return id.Name.Value
}

type VariableModifierExpr struct {
	Mod *TokenExpr `keyword:"readonly"`
}

type TokenExpr Token

func (t TokenExpr) String() string {
	return Token(t).String()
}

type Ellipsis struct {
	Dot1 *TokenExpr `symbol:"Dot"`
	Dot2 *TokenExpr `symbol:"Dot"`
	Dot3 *TokenExpr `symbol:"Dot"`
}

type ParamExpr struct {
	Ellipsis *Ellipsis `optional:"true"`
	Ident    *KeywordIdent
	Optional *TokenExpr      `symbol:"QuestionMark" optional:"true"`
	Type     *TypeAnnotation `optional:"true"`
}

type GetterDecl struct {
	Get    *TokenExpr `word:"get"`
	Ident  *Ident
	Lparen *TokenExpr `symbol:"Lparen"`
	Rparen *TokenExpr `symbol:"Rparen"`
	Return *TypeAnnotation
	Semi   *TokenExpr `symbol:"Semi" optional:"true"`
}

type SetterDecl struct {
	Set    *TokenExpr `word:"set"`
	Ident  *Ident
	Lparen *TokenExpr `symbol:"Lparen"`
	Param  *ParamExpr
	Rparen *TokenExpr `symbol:"Rparen"`
	Semi   *TokenExpr `symbol:"Semi" optional:"true"`
}

type MethodDecl struct {
	Ident   *Ident
	Generic *GenericDecl `optional:"true"`
	Lparen  *TokenExpr   `symbol:"Lparen"`
	Params  []*ParamExpr `delim:"Comma"`
	Rparen  *TokenExpr   `symbol:"Rparen"`
	Return  *TypeAnnotation
	Semi    *TokenExpr `symbol:"Semi" optional:"true"`
}

type CallableDecl struct {
	Generic *GenericDecl `optional:"true"`
	Lparen  *TokenExpr   `symbol:"Lparen"`
	Params  []*ParamExpr `delim:"Comma"`
	Rparen  *TokenExpr   `symbol:"Rparen"`
	Return  *TypeAnnotation
	Semi    *TokenExpr `symbol:"Semi" optional:"true"`
}

type VariableIdent interface {
	Ident() string
}

type StringIdent struct {
	Name *TokenExpr `class:"String"`
}

func (id StringIdent) Ident() string {
	return id.Name.Value
}

type VariableDecl struct {
	Mod      *VariableModifierExpr `optional:"true"`
	Ident    VariableIdent
	Optional *TokenExpr `symbol:"QuestionMark" optional:"true"`
	Type     *TypeAnnotation
	Semi     *TokenExpr `symbol:"Semi" optional:"true"`
}

type IndexerDecl struct {
	Lbrkt *TokenExpr `symbol:"Lbrkt"`
	Param *ParamExpr
	Rbrkt *TokenExpr `symbol:"Rbrkt"`
	Type  *TypeAnnotation
	Semi  *TokenExpr `symbol:"Semi" optional:"true"`
}

type ExtendsExpr struct {
	Extends    *TokenExpr      `keyword:"extends"`
	Interfaces []InterfaceExpr `delim:"Comma"`
}

type InterfaceDecl struct {
	Interface *TokenExpr `keyword:"interface"`
	Ident     *Ident
	Generic   *GenericDecl `optional:"true"`
	Extends   *ExtendsExpr `optional:"true"`
	Body      *TypeBodyExpr
}

type GlobalVarDecl struct {
	Declare  *TokenExpr `keyword:"declare" optional:"true"`
	Var      *TokenExpr `keyword:"var"`
	Variable *VariableDecl
}

type GlobalConstDecl struct {
	Declare  *TokenExpr `keyword:"declare" optional:"true"`
	Const    *TokenExpr `keyword:"const"`
	Variable *VariableDecl
}

type FunctionDecl struct {
	Declare *TokenExpr `keyword:"declare" optional:"true"`
	Var     *TokenExpr `keyword:"function"`
	Method  *MethodDecl
}

type TypeDecl struct {
	Type    *TokenExpr `keyword:"type"`
	Ident   *Ident
	Generic *GenericDecl `optional:"true"`
	Eq      *TokenExpr   `symbol:"Eq"`
	Value   TypeExpr
	Semi    *TokenExpr `symbol:"Semi" optional:"true"`
}

func (TypeDecl) isDtsDecl()        {}
func (GlobalVarDecl) isDtsDecl()   {}
func (GlobalConstDecl) isDtsDecl() {}
func (InterfaceDecl) isDtsDecl()   {}
func (NamespaceDecl) isDtsDecl()   {}
func (FunctionDecl) isDtsDecl()    {}

type NamespaceDecl struct {
	Declare   *TokenExpr `keyword:"declare"`
	Namespace *TokenExpr `keyword:"namespace"`
	Ident     *Ident
	Lbrc      *TokenExpr `symbol:"Lbrc"`
	Decls     []DtsDecl
	Rbrc      *TokenExpr `symbol:"Rbrc"`
}

type DtsDecl interface {
	isDtsDecl()
}

type DtsModule struct {
	Decls []DtsDecl
}

type TypeAnnotation struct {
	Colon *TokenExpr `symbol:"Colon"`
	Type  TypeExpr
}
