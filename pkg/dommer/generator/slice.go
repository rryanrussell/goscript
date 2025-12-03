package generator

type Slice struct {
	Element Type
}

func (s *Slice) Walk(v Visitor) {
	v.Visit(s.Element)
}

func (s *Slice) Print(c PrintContext) {
	WriteAll(c, Lbrkt, Rbrkt, s.Element)
}

func (Slice) IsType() {}
