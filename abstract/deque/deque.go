// Package deque provides an interface for deque data structures.
//
// Deque (Double-ended queue) is an abstract data structure that supports the
// following operations.

package stack

type Deque[T any] interface {
	PushBack(element T) error
	PushFront(element T) error
	PopBack() (T, error)
	PopFront() (T, error)
	Back() (T, error)
	Front() (T, error)
}
