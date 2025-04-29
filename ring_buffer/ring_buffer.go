package ring_buffer

import (
	"fmt"
	"strings"
)

type RingBuffer[T any] struct {
	data   []T
	front  int
	back   int
	length int
}

func NewRingBuffer[T any](capacity int) *RingBuffer[T] {
	return &RingBuffer[T]{
		data:   make([]T, capacity, capacity),
		front:  0,
		back:   0,
		length: 0,
	}
}

type RingBufferFullError struct{
	lastElement any
}

func (e *RingBufferFullError) Error() string {
	return fmt.Sprintf("Ring buffer is full. Last element: %v", e.lastElement)
}

type RingBufferEmptyError struct{}

func (e *RingBufferEmptyError) Error() string {
	return "Ring buffer is empty."
}

func (q *RingBuffer[T]) IsEmpty() bool {
	return q.length == 0
}

func (q *RingBuffer[T]) IsFull() bool {
	return q.length == cap(q.data)
}

func (q *RingBuffer[T]) PushBack(element T) error {
	if q.IsFull() {
		return &RingBufferFullError{ lastElement: element }
	}
	if !q.IsEmpty() {
		q.back = (q.back + 1) % cap(q.data)
	}

	q.data[q.back] = element
	q.length++

	return nil
}

func (q *RingBuffer[T]) Enqueue(element T) error {
	return q.PushBack(element)
}

func (q *RingBuffer[T]) PopFront() (T, error) {
	var result T

	if q.IsEmpty() {
		return result, &RingBufferEmptyError{}
	}

	result = q.data[q.front]
	q.length--
	q.front = (q.front + 1) % cap(q.data)

	return result, nil
}

func (q *RingBuffer[T]) Dequeue() (T, error) {
	return q.PopFront()
}

func (q *RingBuffer[T]) PushBackOver(element T) T {
	var result T

	if q.IsFull() {
		result, _ = q.PopFront()
	}
	if !q.IsEmpty() {
		q.back = (q.back + 1) % cap(q.data)
	}

	q.data[q.back] = element
	q.length++

	return result
}

func (q *RingBuffer[T]) PeekFront() (T, error) {
	var result T

	if q.IsEmpty() {
		return result, &RingBufferEmptyError{}
	}
	result = q.data[q.front]

	return result, nil
}

func (q *RingBuffer[T]) PeekBack() (T, error) {
	var result T

	if q.IsEmpty() {
		return result, &RingBufferEmptyError{}
	}
	result = q.data[q.back]

	return result, nil
}

func (q *RingBuffer[T]) Peek() (T, error) {
	return q.PeekFront()
}

func (q *RingBuffer[T]) Clear() {
	q.front = 0
	q.back = 0
	q.length = 0
}

func (q *RingBuffer[T]) Print() string {
	var sb strings.Builder
	sb.WriteString("RingBuffer: ")

	if !q.IsEmpty() {
		head := q.front

		for head != q.back {
			sb.WriteString(fmt.Sprintf("%v", q.data[head]))
			sb.WriteString("->")
			head = (head + 1) % cap(q.data)
		}

		sb.WriteString(fmt.Sprintf("%v", q.data[q.back]))
	}

	sb.WriteByte('\n')

	return sb.String()
}
