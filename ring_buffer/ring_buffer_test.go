package ring_buffer

import "testing"

func TestNewRingBuffer(t *testing.T) {
	queue := NewRingBuffer[rune](32)

	if queue == nil {
		t.Error("Ring Buffer creation failed.")
	}
}

func TestRingBufferEnqueue(t *testing.T) {
	queue := NewRingBuffer[rune](4)
	runeList := [4]rune{'a', 'b', 'c', 'd'}

	for _, r := range runeList {
		err := queue.Enqueue(r)
		if err != nil {
			t.Error("Ring Buffer Enqueue failed.")
		}
	}
	if queue.Enqueue('e') == nil {
		t.Error("Ring Buffer Enqueue should have failed.")
	}
}

func TestRingBufferDequeue(t *testing.T) {
	queue := NewRingBuffer[rune](4)
	runeList := [4]rune{'a', 'b', 'c', 'd'}
	for _, r := range runeList {
		queue.Enqueue(r)
	}
}
