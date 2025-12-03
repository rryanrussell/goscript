package js

type Ident string

func (i Ident) Print(c PrintContext) {
	c.WriteString(string(i))
}

func IdentP(s string) *Ident {
	id := Ident(s)
	return &id
}

func (Ident) Walk(Visitor) {}

type Namer interface {
	GetName() Ident
	SetName(name Ident)
}

type Named struct {
	Name *Ident
}

func NewNamed(s string) Named {
	return Named{IdentP(s)}
}

func (n *Named) NameEq(s string) bool {
	return string(*n.Name) == s
}

func (n *Named) NameString() string {
	return string(*n.Name)
}

func (n *Named) GetName() Ident {
	return *n.Name
}

func (n *Named) SetName(name Ident) {
	n.Name = &name
}

type Ref struct {
	Named
}

func (r *Ref) Print(c PrintContext) {
	r.Name.Print(c)
}

func (r *Ref) Walk(v Visitor) {
	v.Visit(r.Name)
}

func (Ref) exprNode() {}
