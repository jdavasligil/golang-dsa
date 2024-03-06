package ring_buffer

import "testing"

func TestNewRingBuffer (t *testing.T) {
    queue := NewRingBuffer[rune](32)

    if (queue == nil) {
        t.Error("StackList creation failed.")
    }
}
