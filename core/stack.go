package core

type Stack[T any] struct {
	seq []T
	len int
}

func NewStack[T any]() Stack[T] {
	return Stack[T]{
		[]T{},
		0,
	}
}

func (s *Stack[T]) Push(item T) {
	s.seq = append(s.seq, item)
	s.len = len(s.seq)
}

func (s *Stack[T]) Pop() (T, bool) {
	if s.len == 0 {
		return *new(T), false
	}
	res := s.Peek()
	s.seq = s.seq[:s.len-1]
	s.len -= 1
	return res, true
}

func (s *Stack[T]) Peek() T {
	res := s.seq[s.len-1]
	return res
}
