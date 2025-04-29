package ring_buffer

import (
	"errors"
	"testing"
)

func TestRingBuffer(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		queue := NewRingBuffer[rune](32)

		if queue == nil {
			t.Fatal("Ring Buffer creation failed.")
		}
	})

	t.Run("EnqueueDequeue", func(t *testing.T) {
		queue := NewRingBuffer[rune](4)

		if queue == nil {
			t.Fatal("Ring Buffer creation failed.")
		}
		runeList := [4]rune{'a', 'b', 'c', 'd'}

		for _, r := range runeList {
			err := queue.Enqueue(r)
			if err != nil {
				t.Error("Ring Buffer Enqueue failed.")
			}
		}

		if err := queue.Enqueue('e'); !errors.Is(err, &RingBufferFullError{}) {
			t.Errorf("\n\nGot:      %v\nExpected: RingBufferFullError.\n\n", err)
		}

		for _, r := range runeList {
			s, err := queue.Dequeue()
			if err != nil {
				t.Error("Ring Buffer Dequeue failed.")
			}
			if s != r {
				t.Error("Ring Buffer is not FIFO")
			}
		}

		if _, err := queue.Dequeue(); !errors.Is(err, &RingBufferEmptyError{}) {
			t.Errorf("Got: %v  Expected: RingBufferEmptyError.", err)
		}
	})
}
