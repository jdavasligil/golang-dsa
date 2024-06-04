// ArrayList
// FILO data structure using a dynamic array

package arraylist

import "fmt"

type ArrayList[T any] struct {
	Data []T
}

func NewArrayList[T any]() *ArrayList[T] {
	return &ArrayList[T]{
		Data: make([]T, 0, 1024),
	}
}

func (s *ArrayList[T]) Push(data T) {
	s.Data = append(s.Data, data)
}

func (s *ArrayList[T]) Pop() (T, error) {
	var result T
	var err error = fmt.Errorf("Failed to pop from empty ArrayList.")

	if len(s.Data) == 0 {
		return result, err
	}

	result = s.Data[len(s.Data)-1]
	s.Data = s.Data[0 : len(s.Data)-1]

	return result, nil
}

func (s *ArrayList[T]) Top() (T, error) {
	var result T
	var err error = fmt.Errorf("ArrayList is empty.")

	if len(s.Data) == 0 {
		return result, err
	}

	result = s.Data[len(s.Data)-1]

	return result, nil
}

func (s *ArrayList[T]) Len() int {
	return len(s.Data)
}
