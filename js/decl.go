package js

type Decl interface {
	Node
	declNode()
}
