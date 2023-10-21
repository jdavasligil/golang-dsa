// Stack
// FILO data structure using a linked list

package stack

import "fmt"

type node[T any] struct {
    Data T
    Next *node[T]
}

type Stack[T any] struct {
    Head *node[T]
    size uint
}

func NewStack[T any]() *Stack[T] {
    return &Stack[T] {
        Head: nil,
        size: 0,
    }
}

func (s *Stack[T]) Push(data T) {
    s.Head = &node[T] { Data: data, Next: s.Head }
    s.size += 1
}

func (s *Stack[T]) Pop() (T, error) {
    var result T
    var err error = fmt.Errorf("Failed to pop from empty stack.")

    if (s.Head == nil) {
        return result, err
    }

    result = s.Head.Data

    s.Head = s.Head.Next
    s.size -= 1

    return result, nil
}

func (s *Stack[T]) Top() (T, error) {
    var result T
    var err error = fmt.Errorf("Stack is empty.")

    if (s.Head == nil) {
        return result, err
    }

    result = s.Head.Data

    return result, nil
}

func (s *Stack[T]) Len() uint {
    return s.size
}
