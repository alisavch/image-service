package broker

import (
	"context"
	"time"
)

// Action int type.
type Action int

const (
	// Succeed indicates successful completion.
	Succeed Action = iota
	// Fail indicates unsuccessful completion.
	Fail
	// Retry indicates that you need to try again.
	Retry
)

// RetryPolicy returns Action.
type RetryPolicy func(err error) Action

// Repeater contains backoff and retry.
type Repeater struct {
	backoff     *Backoff
	retryPolicy RetryPolicy
}

// NewRepeater configures Repeater.
func NewRepeater(backoff *Backoff, policy RetryPolicy) Repeater {
	if policy == nil {
		policy = DefaultRetryPolicy
	}

	return Repeater{
		backoff:     backoff,
		retryPolicy: policy,
	}
}

func (r *Repeater) sleep(ctx context.Context, t <-chan time.Time) error {
	select {
	case <-t:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// DefaultRetryPolicy returns default retry policy.
func DefaultRetryPolicy(err error) Action {
	if err != nil {
		return Retry
	}
	return Succeed
}
