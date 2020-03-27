package retry

import (
	"context"
	"fmt"
	"math"
	"time"
)

type RetryErrScope string

const (
	MaxDurationScope RetryErrScope = "max-duration"
	ContextScope     RetryErrScope = "context"
)

type RetryErr struct {
	Scope RetryErrScope
	Err   error
}

func (err RetryErr) Error() string {
	return fmt.Sprintf("retry error (%v): %v", err.Scope, err.Err)
}

type Retryable func(ctx context.Context) error

type Retry interface {
	Do(ctx context.Context, method Retryable) error
}

type Retryer struct {
	waitDuration time.Duration
	maxDuration  time.Duration
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

func WithMaxDuration(duration time.Duration) RetryerOptsFunc {
	return func(r *Retryer) {
		r.maxDuration = duration
	}
}

func WithoutMaxAttempts() RetryerOptsFunc {
	return func(r *Retryer) {
		r.maxAttempts = math.MaxInt32
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

// Do execute method following rules of the Retry struct
// Two timeouts co-exist:
// * The one given as param of 'method': can be the scope of the current
// http.Request for instance
// * The one defined with the option WithMaxDuration, which would cancel the
// retry loop if it has expired.
func (r Retryer) Do(ctx context.Context, method Retryable) error {
	timeoutCtx := context.Background()
	if r.maxDuration != 0 {
		var cancel func()
		timeoutCtx, cancel = context.WithTimeout(ctx, r.maxDuration)
		defer cancel()
	}

	var err error
	for i := 0; i < r.maxAttempts; i++ {
		err = method(ctx)
		if err == nil {
			return nil
		}

		timer := time.NewTimer(r.waitDuration)
		select {
		case <-timer.C:
		case <-timeoutCtx.Done():
			return RetryErr{
				Scope: MaxDurationScope,
				Err:   timeoutCtx.Err(),
			}
		case <-ctx.Done():
			return RetryErr{
				Scope: ContextScope,
				Err:   ctx.Err(),
			}
		}
	}
	return err
}
