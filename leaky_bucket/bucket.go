// Adaptive rate FIFO MPSC queue
//
// TODO:
// 		- Unit Testing
// 		- Rate Interpolation

package bucket

import (
	"errors"
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

type BucketConstraintError struct {
	Constraint string
}

func (e *BucketConstraintError) Error() string {
	return fmt.Sprintf("Constraint violated: %s", e.Constraint)
}

func (e *BucketConstraintError) Is(target error) bool {
	_, ok := target.(*BucketConstraintError)
	return ok
}

// A leaky bucket with adaptive rate smoothing to handle backpressure.
type Bucket[T any] struct {
	// Contents of the bucket.
	packets chan T

	// Enables dynamic rate increase proportional to the burst size.
	lowLatency bool

	// Factor to bias rate from equilibrium. 0 < dropBias < 1
	dropBias float64

	// Blocks until the next drop is ready.
	dropTimer *time.Timer

	// Time between packet drops (ns / drop).
	dropInterval time.Duration

	// Minimum time between packet drops (ns / drop).
	minDropInterval time.Duration

	// Maximum time between packet drops (ns / drop).
	maxDropInterval time.Duration

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

type BucketOptions struct {
	Capacity        int
	LowLatency      bool
	DropBias        float64
	DropInterval    time.Duration
	MinDropInterval time.Duration
	MaxDropInterval time.Duration
	MaxWaitTime     time.Duration
	UpdateInterval  time.Duration
}

func NewBucket[T any](opts *BucketOptions) (*Bucket[T], error) {
	var errs error

	if opts.Capacity <= 0 {
		errs = errors.Join(errs, &BucketConstraintError{
			fmt.Sprintf("Capacity %d > 0", opts.Capacity),
		})
	}

	if opts.DropBias <= 0 || opts.DropBias > 1 {
		errs = errors.Join(errs, &BucketConstraintError{
			fmt.Sprintf("0 < Drop Bias %.2f <= 1", opts.DropBias),
		})
	}

	if opts.DropInterval < opts.MinDropInterval {
		errs = errors.Join(errs, &BucketConstraintError{
			fmt.Sprintf("Drop Interval %d > Minimum %d", opts.DropInterval.Milliseconds(), opts.MinDropInterval.Milliseconds()),
		})
	}

	if opts.DropInterval > opts.MaxDropInterval {
		errs = errors.Join(errs, &BucketConstraintError{
			fmt.Sprintf("Drop Interval %d < Maximum %d", opts.DropInterval.Milliseconds(), opts.MaxDropInterval.Milliseconds()),
		})
	}

	if opts.UpdateInterval <= opts.MaxDropInterval {
		errs = errors.Join(errs, &BucketConstraintError{
			fmt.Sprintf("Update Interval %d > Max Drop Interval %d", opts.UpdateInterval.Milliseconds(), opts.MaxDropInterval.Milliseconds()),
		})
	}

	if opts.MaxWaitTime <= opts.MaxDropInterval {
		errs = errors.Join(errs, &BucketConstraintError{
			fmt.Sprintf("Max Wait Time %d > Max Drop Interval %d", opts.MaxWaitTime.Milliseconds(), opts.MaxDropInterval.Milliseconds()),
		})
	}

	if errs != nil {
		return nil, errs
	}

	return &Bucket[T]{
		packets:         make(chan T, opts.Capacity),
		lowLatency:      opts.LowLatency,
		dropTimer:       time.NewTimer(opts.DropInterval),
		dropInterval:    opts.DropInterval,
		minDropInterval: opts.MinDropInterval,
		maxDropInterval: opts.MaxDropInterval,
		dropBias:        opts.DropBias,
		waitTimer:       time.NewTimer(opts.MaxWaitTime),
		maxWaitTime:     opts.MaxWaitTime,
		updateTimer:     time.NewTimer(opts.UpdateInterval),
		updateInterval:  opts.UpdateInterval,
	}, nil
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
//
// Note that this function is blocking (until maxWaitTime).
//
// Expect BucketClosedError once the producer shuts down the bucket.
//
// The bucket will continue to drain even after shutdown. Detect shutdown
// immediately, use a separate channel.
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
	// Ratio of max burst amount to total capacity.
	// 0 < epsilon <= 1
	// Zero:     equilibrium. Rate should not change.
	// Positive: above equilibrium. Rate should increase.
	var epsilon float64

	if b.lowLatency {
		epsilon = float64(b.packetMax.Swap(0)) / float64(cap(b.packets))
	}

	// Prevent singularity problems by clamping at 1.
	packetsSinceLast := max(1, b.packetCount.Swap(0))

	// The estimated average update interval for maintaining equilibrium (ns / drop)
	tau := float64(b.updateInterval) / float64(packetsSinceLast)

	// Set interval to estimated equilibrium scaled by burst ratio and drop bias.
	b.dropInterval = time.Duration(math.Ceil(tau * (1 - epsilon) * b.dropBias))

	// Clamp between min and max
	b.dropInterval = max(b.minDropInterval, min(b.maxDropInterval, b.dropInterval))

	b.updateTimer.Reset(b.updateInterval)
}

// Close must be called by the producer, not the consumer.
func (b *Bucket[T]) Close() {
	close(b.packets)
}

// Drain may be called by the consumer after the producer has closed the bucket.
func (b *Bucket[T]) Drain() []T {
	dropsRemaining := make([]T, 0, len(b.packets))

	for p := range b.packets {
		dropsRemaining = append(dropsRemaining, p)
	}

	return dropsRemaining
}

// Status must be called by the consumer.
func (b *Bucket[T]) Status() string {
	var sb strings.Builder
	sb.WriteString("Bucket Status:\n")

	sb.WriteString(fmt.Sprintf("\tCapacity:             (%d / %d)\n", len(b.packets), cap(b.packets)))
	sb.WriteString(fmt.Sprintf("\tDrops Since Update:    %d\n", b.packetCount.Load()))
	sb.WriteString(fmt.Sprintf("\tCurrent Drop Rate(ms): %d\n", b.dropInterval.Milliseconds()))

	return sb.String()
}
