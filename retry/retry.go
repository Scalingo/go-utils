package retry

import (
	"context"
	"time"
)

type Retryable func(ctx context.Context) error

type Retry interface {
	Do(ctx context.Context, method Retryable) error
}

type Retryer struct {
	waitDuration time.Duration
	maxAttempts  int
}

type RetryerOptsFunc func(r *Retryer)

func WithWaitDuration(duration time.Duration) RetryerOptsFunc {
	return func(r *Retryer) {
		r.waitDuration = duration
	}
}

func WithMaxAttempts(maxAttempts int) RetryerOptsFunc {
	return func(r *Retryer) {
		r.maxAttempts = maxAttempts
	}
}

func New(opts ...RetryerOptsFunc) Retryer {
	r := &Retryer{
		waitDuration: 10 * time.Second,
		maxAttempts:  5,
	}

	for _, opt := range opts {
		opt(r)
	}

	return *r
}

func (r Retryer) Do(ctx context.Context, method Retryable) error {
	var err error
	for i := 0; i < r.maxAttempts; i++ {
		err = method(ctx)
		if err == nil {
			return nil
		}

		timer := time.NewTimer(r.waitDuration)
		select {
		case <-timer.C:
		case <-ctx.Done():
			deadLineErr := ctx.Err()
			if deadLineErr != nil {
				return deadLineErr
			}
		}
	}
	return err
}
