package retry

import (
	"context"
	stderrors "errors"
	"fmt"
	"testing"
	"testing/synctest"
	"time"

	"github.com/sirupsen/logrus"
	logrustest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Scalingo/go-utils/logger"
)

func TestRetrier(t *testing.T) {
	t.Run("When the method works fine the first time", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			retrier := New()
			tries := 0
			err := retrier.Do(t.Context(), func(ctx context.Context) error {
				tries++
				return nil
			})
			synctest.Wait()

			require.NoError(t, err)
			assert.Equal(t, 1, tries)
		})
	})

	t.Run("When the method works fine the second time", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			retrier := New(WithWaitDuration(100 * time.Millisecond))
			tries := 0
			before := time.Now()
			err := retrier.Do(t.Context(), func(ctx context.Context) error {
				tries++
				if tries == 2 {
					return nil
				}
				return fmt.Errorf("Error attempt %v", tries)
			})
			duration := time.Since(before)
			synctest.Wait()

			require.NoError(t, err)
			assert.Equal(t, 2, tries)
			if duration < 100*time.Millisecond {
				t.Fatalf("Test should take at least 100ms, took %v", duration)
			}
		})
	})

	t.Run("When the method works fine the third time and exponential backoff is enabled", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			retrier := New(
				WithWaitDuration(200*time.Millisecond),
				WithExponentialBackoff(2),
			)

			tries := 0
			before := time.Now()
			err := retrier.Do(t.Context(), func(ctx context.Context) error {
				tries++
				if tries == 3 {
					return nil
				}
				return fmt.Errorf("Error attempt %v", tries)
			})
			duration := time.Since(before)
			synctest.Wait()

			require.NoError(t, err)
			assert.Equal(t, 3, tries)
			expectedDuration := 600 * time.Millisecond
			if duration < expectedDuration {
				t.Fatalf("Test should take at least %v, took %v", expectedDuration, duration)
			}
		})
	})

	t.Run("When the method never succeeds", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			retrier := New(WithWaitDuration(1 * time.Millisecond))
			err := retrier.Do(t.Context(), func(ctx context.Context) error {
				return stderrors.New("nop")
			})
			synctest.Wait()

			require.Error(t, err)
		})
	})

	t.Run("It should cancel the retry if a RetryCancelError is retuned", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			retrier := New(WithWaitDuration(1 * time.Millisecond))
			count := 0
			err := retrier.Do(t.Context(), func(ctx context.Context) error {
				count++
				return NewRetryCancelError(stderrors.New("nop"))
			})
			synctest.Wait()

			require.Error(t, err)
			assert.Equal(t, 1, count)
		})
	})

	t.Run("When the context is canceled after the first try", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			retrier := New(WithWaitDuration(1 * time.Millisecond))
			tries := 0
			ctx, cancel := context.WithCancel(t.Context())
			err := retrier.Do(ctx, func(ctx context.Context) error {
				tries++
				if tries == 2 {
					return nil
				}
				cancel()
				return fmt.Errorf("Error attempt %v", tries)
			})
			synctest.Wait()

			require.Error(t, err)
			assert.Equal(t, 1, tries)
		})
	})

	t.Run("With timeout should ignore sleep", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			retrier := New(WithWaitDuration(1 * time.Second))
			ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)
			defer cancel()

			before := time.Now()
			err := retrier.Do(ctx, func(ctx context.Context) error {
				return stderrors.New("retry test error")
			})
			synctest.Wait()

			require.Error(t, err)
			assert.WithinDuration(t, time.Now(), before, 100*time.Millisecond)

			var retryError RetryError
			require.ErrorAs(t, err, &retryError)
			assert.Equal(t, ContextScope, retryError.Scope)
			assert.Equal(t, context.DeadlineExceeded, retryError.Err)
			require.ErrorContains(t, retryError.LastErr, "retry test error")
		})
	})

	t.Run("With max duration it should ignore sleep", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			retrier := New(
				WithWaitDuration(1*time.Second),
				WithMaxDuration(50*time.Millisecond),
			)

			before := time.Now()
			err := retrier.Do(t.Context(), func(ctx context.Context) error {
				return stderrors.New("max duration error")
			})
			synctest.Wait()

			require.Error(t, err)
			assert.WithinDuration(t, time.Now(), before, 100*time.Millisecond)

			var retryError RetryError
			require.ErrorAs(t, err, &retryError)
			assert.Equal(t, MaxDurationScope, retryError.Scope)
			assert.Equal(t, retryError.Err, context.DeadlineExceeded)
			require.ErrorContains(t, retryError.LastErr, "max duration error")
		})
	})

	t.Run("With a callback", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			callbackCalls := 0
			retrier := New(
				WithMaxAttempts(1),
				WithWaitDuration(50*time.Millisecond),
				WithWaitDuration(100*time.Millisecond),
				WithErrorCallback(func(ctx context.Context, err error, currentAttempt, maxAttempts int) {
					callbackCalls++
					require.ErrorContains(t, err, "TestError")
					assert.Equal(t, 0, currentAttempt)
					assert.Equal(t, 1, maxAttempts)
				}),
			)

			err := retrier.Do(t.Context(), func(ctx context.Context) error {
				return stderrors.New("TestError")
			})
			synctest.Wait()

			require.Error(t, err)
			assert.Equal(t, 1, callbackCalls)
		})
	})

	t.Run("If both timeout are specified, the first one which is expired should exist the method", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			// Timeout from call context first
			retrier := New(
				WithWaitDuration(1*time.Second),
				WithMaxDuration(500*time.Millisecond),
			)

			ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)
			defer cancel()

			before := time.Now()
			err := retrier.Do(ctx, func(ctx context.Context) error {
				return stderrors.New("test")
			})
			synctest.Wait()

			require.Error(t, err)
			assert.WithinDuration(t, time.Now(), before, 100*time.Millisecond)
			var retryError RetryError
			require.ErrorAs(t, err, &retryError)
			assert.Equal(t, ContextScope, retryError.Scope)
			assert.Equal(t, retryError.Err, context.DeadlineExceeded)

			// Timeout from MaxDuration first
			retrier = New(
				WithWaitDuration(1*time.Second),
				WithMaxDuration(50*time.Millisecond),
			)

			ctx, cancel = context.WithTimeout(t.Context(), 500*time.Millisecond)
			defer cancel()

			before = time.Now()
			err = retrier.Do(ctx, func(ctx context.Context) error {
				return stderrors.New("test")
			})
			synctest.Wait()

			require.Error(t, err)
			assert.WithinDuration(t, time.Now(), before, 100*time.Millisecond)
			require.ErrorAs(t, err, &retryError)
			assert.Equal(t, MaxDurationScope, retryError.Scope)
			assert.Equal(t, retryError.Err, context.DeadlineExceeded)
		})
	})

	t.Run("If logging on attempt error is set, it should log the error", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			log, hook := logrustest.NewNullLogger()
			ctx := logger.ToCtx(t.Context(), log)

			retrier := New(
				WithMaxAttempts(2),
				WithWaitDuration(10*time.Millisecond),
				WithLoggingOnAttemptError(logrus.ErrorLevel),
			)

			err := retrier.Do(ctx, func(_ context.Context) error {
				return stderrors.New("TestError")
			})
			synctest.Wait()

			require.Error(t, err)
			assert.Len(t, hook.Entries, 2)
			assert.Contains(t, "Attempt failed", hook.Entries[0].Message)
			assert.Equal(t, logrus.ErrorLevel, hook.Entries[0].Level)
		})
	})
}

func TestRetryErrorUnwrap(t *testing.T) {
	t.Run("RetryError should unwrap to Err", func(t *testing.T) {
		baseErr := stderrors.New("base error")
		err := RetryError{
			Scope:   ContextScope,
			Err:     baseErr,
			LastErr: stderrors.New("last error"),
		}

		require.ErrorIs(t, err, baseErr)
		assert.Equal(t, baseErr, stderrors.Unwrap(err))
	})

	t.Run("RetryCancelError should unwrap to inner error", func(t *testing.T) {
		baseErr := stderrors.New("cancel error")
		err := NewRetryCancelError(baseErr)

		require.ErrorIs(t, err, baseErr)
		assert.Equal(t, baseErr, stderrors.Unwrap(err))
	})
}
