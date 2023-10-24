// Queue
// FIFO data structure using an array

package queue

import "fmt"

type Queue[T any] struct {
    data   []T
    front  int
    back   int
    length int
}

func NewQueue[T any](capacity int) *Queue[T] {
    return &Queue[T] {
        data:   make([]T, 0, capacity),
        front:  0,
        back:   0,
        length: 0,
    }
}

func (q *Queue[T]) IsEmpty() bool {
    return q.length == 0
}

func (q *Queue[T]) IsFull() bool {
    return q.length == cap(q.data)
}

func (q *Queue[T]) Enqueue(data T) error {
    if (q.IsFull()) {
        return fmt.Errorf("Failure to enqueue to a full queue")
    }

    q.data[q.back] = data
    q.length++
    q.back = (q.back + 1) % cap(q.data)

    return nil
}

func (q *Queue[T]) Dequeue() (T, error) {
    var result T
    var err error = fmt.Errorf("Failed to dequeue from empty queue")

    if (q.IsEmpty()) {
        return result, err
    }

    result = q.data[q.front]
    q.length--
    q.front = (q.front + 1) % cap(q.data)

    return result, nil
}

func (q *Queue[T]) PeekFront() (T, error) {
    var result T
    var err error = fmt.Errorf("Failed to obtain front of empty queue")

    if (q.IsEmpty()) {
        return result, err
    }
    result = q.data[q.front]

    return result, nil
}

func (q *Queue[T]) PeekBack() (T, error) {
    var result T
    var err error = fmt.Errorf("Failed to obtain back of empty queue")

    if (q.IsEmpty()) {
        return result, err
    }
    result = q.data[q.back]

    return result, nil
}

func (q *Queue[T]) Clear() {
    q.front = 0
    q.back = 0
    q.length = 0
}
