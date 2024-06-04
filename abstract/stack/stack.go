// Package stack provides an interface for stack data structures.
//
// Stack is an abstract LIFO data structure that supports the following
// operations.

package stack

type Stack[T any] interface {
	Push(element T) error
	Pop() (T, error)
	Top() (T, error)
}
