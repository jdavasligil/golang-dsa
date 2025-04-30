// Adaptive rate FIFO MPSC queue
package bucket

import (
	"fmt"
	"math"
	"strings"
	"sync/atomic"
	"time"
)

type MaxWaitError struct{}

func (e *MaxWaitError) Error() string {
	return "Max wait time exceeded: no packet dropped."
}

func (e *MaxWaitError) Is(target error) bool {
	_, ok := target.(*MaxWaitError)
	return ok
}

type BucketFullError struct{}

func (e *BucketFullError) Error() string {
	return "Bucket is full. Try again later."
}

func (e *BucketFullError) Is(target error) bool {
	_, ok := target.(*BucketFullError)
	return ok
}

type BucketClosedError struct{}

func (e *BucketClosedError) Error() string {
	return "Bucket is closed. No more drops remain."
}

func (e *BucketClosedError) Is(target error) bool {
	_, ok := target.(*BucketClosedError)
	return ok
}

// A leaky bucket with adaptive rate smoothing to handle backpressure.
type Bucket[T any] struct {
	// Contents of the bucket.
	packets chan T

	// Blocks until the next drop is ready.
	dropTimer *time.Timer

	// Time between packet drops (ns / drop).
	dropInterval time.Duration

	// Maximum time between packet drops (ns / drop).
	maxDropInterval time.Duration

	// Minimum time between packet drops (ns / drop).
	minDropInterval time.Duration

	// Factor to bias rate from equilibrium. 0 < dropBias < 1
	dropBias float64

	// Prevents drop timer and packets from blocking forever.
	waitTimer *time.Timer

	// Max time AwaitDrop will block for.
	maxWaitTime time.Duration

	// Triggers drop timer adaption updates.
	updateTimer *time.Timer

	// Minimum time to wait between updates. Ideally, a multiple of your expected burst time.
	updateInterval time.Duration

	// Counts the number of drops coming in over the update interval.
	packetCount atomic.Int32

	// Tracks the maximum number of packets in the bucket over the update interval.
	packetMax atomic.Int32
}

// Bucket Constraints (unenforced):
//
//		(max expected rate IN) * updateInterval << capacity.
//		0 < target < capacity.
//		minDropInterval < dropInterval < maxDropInterval << updateInterval.
//	 0 < dropBias ~ 1 (0.9 - 1.1)
func NewBucket[T any](
	capacity int,
	dropInterval time.Duration,
	maxDropInterval time.Duration,
	minDropInterval time.Duration,
	dropBias float64,
	maxWaitTime time.Duration,
	updateInterval time.Duration,
) *Bucket[T] {
	return &Bucket[T]{
		packets:         make(chan T, capacity),
		dropTimer:       time.NewTimer(dropInterval),
		dropInterval:    dropInterval,
		maxDropInterval: maxDropInterval,
		minDropInterval: minDropInterval,
		dropBias:        dropBias,
		waitTimer:       time.NewTimer(maxWaitTime),
		maxWaitTime:     maxWaitTime,
		updateTimer:     time.NewTimer(updateInterval),
		updateInterval:  updateInterval,
	}
}

// Multiple producers may add drops to the bucket.
func (b *Bucket[T]) AddDrop(packet T) error {
	if len(b.packets) == cap(b.packets) {
		return &BucketFullError{}
	}
	b.packets <- packet
	b.packetCount.Add(1)

	pmax := b.packetMax.Load()
	b.packetMax.Store(max(int32(len(b.packets)), pmax))

	return nil
}

// A single consumer may await drops from the bucket.
func (b *Bucket[T]) AwaitDrop() (T, error) {
	var packet T

	b.dropTimer.Reset(b.dropInterval)
	b.waitTimer.Reset(b.maxWaitTime)

	// Wait for drop or max wait time, whichever comes first.
	select {
	case <-b.dropTimer.C:
	case <-b.waitTimer.C:
		return packet, &MaxWaitError{}
	}

	// Wait for packet to be ready or max wait time, whichever comes first.
	select {
	case p, ok := <-b.packets:
		packet = p
		if !ok {
			return packet, &BucketClosedError{}
		}
	case <-b.waitTimer.C:
		return packet, &MaxWaitError{}
	}

	select {
	case <-b.updateTimer.C:
		b.adapt()
	default:
	}

	return packet, nil
}

// Change the drop interval in response to new information.
// Currently makes potentially sudden changes.
// Might be improved with interpolation.
func (b *Bucket[T]) adapt() {
	// Normalized error between max burst amount detected and target (zero)
	// -1 <= epsilon <= 1
	// Zero:     equilibrium. Rate should not change.
	// Negative: below equilibrium. Rate should decrease.
	// Positive: above equilibrium. Rate should increase.
	epsilon := float64(b.packetMax.Swap(0)) / float64(cap(b.packets))
	fmt.Printf("EPSILON: %f\n", epsilon)

	// Prevent singularity problems by clamping at 1.
	packetsSinceLast := max(1, b.packetCount.Swap(0))

	// The estimated average update interval for maintaining equilibrium (ns / drop)
	tau := float64(b.updateInterval) / float64(packetsSinceLast)
	fmt.Printf("TAU: %v (ms)\n", time.Duration(tau).Milliseconds())

	// If epsilon is 0, seek equilibrium.
	// If epsilon is -1 (len is furthest below target), approach zero
	// If epsilon is +1 (len is furthest above target), approach double rate
	fmt.Printf("DI BEFORE: %v (ms)\n", b.dropInterval.Milliseconds())
	b.dropInterval = time.Duration(math.Ceil(tau * (1 - epsilon) * b.dropBias))
	fmt.Printf("DI AFTER:  %v (ms)\n", b.dropInterval.Milliseconds())

	// Clamp between min and max
	b.dropInterval = max(b.minDropInterval, min(b.maxDropInterval, b.dropInterval))
	fmt.Printf("DI CLAMPED:  %v (ms)\n", b.dropInterval.Milliseconds())

	b.updateTimer.Reset(b.updateInterval)
}

// Close must be called by the producer, not the consumer.
func (b *Bucket[T]) Close() {
	close(b.packets)
}

// Drain must be called by the consumer after the producer has closed the bucket.
func (b *Bucket[T]) Drain() []T {
	dropsRemaining := make([]T, 0, len(b.packets))

	for p := range b.packets {
		dropsRemaining = append(dropsRemaining, p)
	}

	return dropsRemaining
}

// Status must be called by Consumer.
func (b *Bucket[T]) Status() string {
	var sb strings.Builder
	sb.WriteString("Bucket Status:\n")

	sb.WriteString(fmt.Sprintf("\tCapacity:             (%d / %d)\n", len(b.packets), cap(b.packets)))
	sb.WriteString(fmt.Sprintf("\tDrops Since Update:    %d\n", b.packetCount.Load()))
	sb.WriteString(fmt.Sprintf("\tCurrent Drop Rate(ms): %d\n", b.dropInterval.Milliseconds()))

	return sb.String()
}
