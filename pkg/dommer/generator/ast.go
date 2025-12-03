package generator

type Visitor interface {
	Visit(node Node)
}

type Node interface {
	Printer
	Walk(visitor Visitor)
}

func Stroll[T Node](v Visitor, nodes []T) {
	for _, node := range nodes {
		v.Visit(node)
	}
}
