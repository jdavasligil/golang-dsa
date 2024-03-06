package ring_buffer

import "fmt"

type RingBuffer[T any] struct {
    data   []T
    front  int
    back   int
    length int
}

func NewRingBuffer[T any](capacity int) *RingBuffer[T] {
    return &RingBuffer[T] {
        data:   make([]T, capacity, capacity),
        front:  0,
        back:   0,
        length: 0,
    }
}

func (q *RingBuffer[T]) IsEmpty() bool {
    return q.length == 0
}

func (q *RingBuffer[T]) IsFull() bool {
    return q.length == cap(q.data)
}

func (q *RingBuffer[T]) PushFront(element T) error {
    if (q.IsFull()) {
        return fmt.Errorf("Failure to push element to front: buffer is full.")
    }
    if (!q.IsEmpty()) {
        q.back = (q.back + 1) % cap(q.data)
    }

    q.data[q.back] = element
    q.length++

    return nil
}

func (q *RingBuffer[T]) PopBack() (T, error) {
    var result T
    var err error = fmt.Errorf("Failed to pop element from back: buffer is empty.")

    if (q.IsEmpty()) {
        return result, err
    }

    result = q.data[q.front]
    q.length--
    q.front = (q.front + 1) % cap(q.data)

    return result, nil
}

func (q *RingBuffer[T]) Enqueue(element T) error {
    return q.PushFront(element)
}

func (q *RingBuffer[T]) PushFrontOver(element T) T {
    var result T

    if (q.IsFull()) {
        result, _ = q.PopBack()
    }
    if (!q.IsEmpty()) {
        q.back = (q.back + 1) % cap(q.data)
    }

    q.data[q.back] = element
    q.length++

    return result
}

func (q *RingBuffer[T]) PeekFront() (T, error) {
    var result T
    var err error = fmt.Errorf("Failed to obtain front of empty queue")

    if (q.IsEmpty()) {
        return result, err
    }
    result = q.data[q.front]

    return result, nil
}

func (q *RingBuffer[T]) PeekBack() (T, error) {
    var result T
    var err error = fmt.Errorf("Failed to obtain back of empty queue")

    if (q.IsEmpty()) {
        return result, err
    }
    result = q.data[q.back]

    return result, nil
}

func (q *RingBuffer[T]) Clear() {
    q.front = 0
    q.back = 0
    q.length = 0
}

func (q *RingBuffer[T]) Print() {
    print("RingBuffer: ")

    if !q.IsEmpty() {
        head := q.front

        for head != q.back {
            print(q.data[head])
            print("->")
            head = (head + 1) % cap(q.data)
        }

        print(q.data[q.back])
    }

    print("\n")
}
