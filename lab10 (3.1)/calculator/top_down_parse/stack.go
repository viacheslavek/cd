package top_down_parse

import "errors"

type Stack[T any] struct {
	buffer []T
}

func NewStack[T any]() Stack[T] {
	return Stack[T]{
		buffer: make([]T, 0),
	}
}

func (s *Stack[T]) Push(elem T) {
	s.buffer = append(s.buffer, elem)
}

func (s *Stack[T]) Pop() (T, error) {
	if len(s.buffer) > 0 {
		elem := s.buffer[len(s.buffer)-1]
		s.buffer = s.buffer[:len(s.buffer)-1]
		return elem, nil
	}
	var tmp T
	return tmp, errors.New("empty buffer")
}

func (s *Stack[T]) GetElems() []T {
	return s.buffer
}
