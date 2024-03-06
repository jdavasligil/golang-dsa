// Package queue provides an interface for queue data structures.
//
// Queue is an abstract FIFO data structure that supports the following
// operations.

package queue

type Queue[T any] interface {
    Enqueue(element T) error
    Dequeue() (T, error)
    Peek() (T, error)
}
