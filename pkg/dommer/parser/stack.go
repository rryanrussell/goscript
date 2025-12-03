package parser

type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(val T) {
	s.items = append(s.items, val)
}

func (s *Stack[T]) Queue() []T {
	return s.items
}

func (s *Stack[T]) Pop() (T, bool) {
	val, ok := s.Peek()

	if !ok {
		return val, false
	}

	s.items = s.items[:len(s.items)-1]

	return val, true
}

func (s *Stack[T]) Peek() (T, bool) {
	if s.Empty() {
		var t T
		return t, false
	}

	return s.items[len(s.items)-1], true
}

func (s *Stack[T]) Len() int {
	return len(s.items)
}

func (s *Stack[T]) Empty() bool {
	return len(s.items) == 0
}
