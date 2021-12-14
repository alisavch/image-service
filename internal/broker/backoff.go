package broker

import (
	"math"
	"math/rand"
	"time"
)

// Stop stops backoff.
const Stop time.Duration = -1

// FnBackoff returns functions duration.
type FnBackoff func(attemptNum int, min, max time.Duration) time.Duration

// Backoff contains values for backoff.
type Backoff struct {
	min, max   time.Duration
	maxAttempt int
	attemptNum int
	backoff    FnBackoff
}

// NewBackoff configures Backoff.
func NewBackoff(min, max time.Duration, maxAttempt int, backoff FnBackoff) *Backoff {
	if backoff == nil {
		backoff = ExponentialBackoff
	}
	return &Backoff{
		min:        min,
		max:        max,
		maxAttempt: maxAttempt,
		backoff:    backoff,
	}
}

// Next interrupts all attempts.
func (b *Backoff) Next() time.Duration {
	if b.attemptNum > b.maxAttempt {
		return Stop
	}
	b.attemptNum++
	return b.backoff(b.attemptNum, b.min, b.max)
}

// Reset resets all attempts.
func (b *Backoff) Reset() {
	b.attemptNum = 1
}

// ExponentialBackoff is performed exponentially.
func ExponentialBackoff(attemptNum int, min, max time.Duration) time.Duration {
	factor := 2.0
	rand.Seed(time.Now().UnixNano())
	delay := time.Duration(math.Pow(factor, float64(attemptNum)) * float64(min))
	jitter := time.Duration(rand.Float64() * float64(min) * float64(attemptNum))

	delay += jitter
	if delay > max {
		delay = max
	}

	return delay
}
