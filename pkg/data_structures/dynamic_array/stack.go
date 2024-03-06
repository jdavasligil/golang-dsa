// Stack
// FILO data structure using a dynamic array

package stack

import "fmt"

type Stack[T any] struct {
    Data []T
}

func NewStack[T any]() *Stack[T] {
    return &Stack[T] {
        Data: make([]T, 0, 1024),
    }
}

func (s *Stack[T]) Push(data T) {
    s.Data = append(s.Data, data)
}

func (s *Stack[T]) Pop() (T, error) {
    var result T
    var err error = fmt.Errorf("Failed to pop from empty stack.")

    if (len(s.Data) == 0) {
        return result, err
    }

    result = s.Data[len(s.Data) - 1]
    s.Data = s.Data[0:len(s.Data) - 1]

    return result, nil
}

func (s *Stack[T]) Top() (T, error) {
    var result T
    var err error = fmt.Errorf("Stack is empty.")

    if (len(s.Data) == 0) {
        return result, err
    }

    result = s.Data[len(s.Data) - 1]

    return result, nil
}

func (s *Stack[T]) Len() int {
    return len(s.Data)
}
