// StackList
// FILO data structure using a linked list

package stack_list

import "fmt"

type node[T any] struct {
	Data T
	Next *node[T]
}

type StackList[T any] struct {
	Head *node[T]
	size int
}

func NewStackList[T any]() *StackList[T] {
	return &StackList[T]{
		Head: nil,
		size: 0,
	}
}

func (s *StackList[T]) Push(data T) {
	s.Head = &node[T]{Data: data, Next: s.Head}
	s.size += 1
}

func (s *StackList[T]) Pop() (T, error) {
	var result T
	var err error = fmt.Errorf("Failed to pop from empty StackList.")

	if s.Head == nil {
		return result, err
	}

	result = s.Head.Data

	s.Head = s.Head.Next
	s.size -= 1

	return result, nil
}

func (s *StackList[T]) Top() (T, error) {
	var result T
	var err error = fmt.Errorf("StackList is empty.")

	if s.Head == nil {
		return result, err
	}

	result = s.Head.Data

	return result, nil
}

func (s *StackList[T]) Len() int {
	return s.size
}
