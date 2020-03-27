package retry

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetrier(t *testing.T) {
	t.Run("When the method works fine the first time", func(t *testing.T) {
		retrier := New()
		tries := 0
		err := retrier.Do(context.Background(), func(ctx context.Context) error {
			tries++
			return nil
		})

		assert.NoError(t, err)
		assert.Equal(t, tries, 1)
	})

	t.Run("When the method works fine the second time", func(t *testing.T) {
		retrier := New(WithWaitDuration(100 * time.Millisecond))
		tries := 0
		before := time.Now()
		err := retrier.Do(context.Background(), func(ctx context.Context) error {
			tries++
			if tries == 2 {
				return nil
			}
			return fmt.Errorf("Error attempt %v", tries)
		})
		duration := time.Now().Sub(before)

		assert.NoError(t, err)
		assert.Equal(t, tries, 2)
		if duration < 100*time.Millisecond {
			t.Fatalf("Test should take at least 100ms, took %v", duration)
		}
	})

	t.Run("When the method never returns", func(t *testing.T) {
		retrier := New(WithWaitDuration(1 * time.Millisecond))
		err := retrier.Do(context.Background(), func(ctx context.Context) error {
			return errors.New("nop")
		})

		assert.Error(t, err)
	})

	t.Run("When the context is canceled after the first try", func(t *testing.T) {
		retrier := New(WithWaitDuration(1 * time.Millisecond))
		tries := 0
		ctx, cancel := context.WithCancel(context.Background())
		err := retrier.Do(ctx, func(ctx context.Context) error {
			tries++
			if tries == 2 {
				return nil
			}
			cancel()
			return fmt.Errorf("Error attempt %v", tries)
		})

		assert.Error(t, err)
		assert.Equal(t, tries, 1)
	})

	t.Run("With timeout should ignore sleep", func(t *testing.T) {
		retrier := New(WithWaitDuration(1 * time.Second))
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		before := time.Now()
		err := retrier.Do(ctx, func(ctx context.Context) error {
			return errors.New("test")
		})

		assert.Error(t, err)
		assert.WithinDuration(t, time.Now(), before, 100*time.Millisecond)
		require.IsType(t, RetryErr{}, err)
		assert.EqualValues(t, err.(RetryErr).Scope, ContextScope)
		assert.Equal(t, err.(RetryErr).Err, context.DeadlineExceeded)
	})

	t.Run("With max duration it should ignore sleep", func(t *testing.T) {
		retrier := New(
			WithWaitDuration(1*time.Second),
			WithMaxDuration(50*time.Millisecond),
		)

		before := time.Now()
		err := retrier.Do(context.Background(), func(ctx context.Context) error {
			return errors.New("test")
		})

		assert.Error(t, err)
		assert.WithinDuration(t, time.Now(), before, 100*time.Millisecond)
		require.IsType(t, RetryErr{}, err)
		assert.EqualValues(t, err.(RetryErr).Scope, MaxDurationScope)
		assert.Equal(t, err.(RetryErr).Err, context.DeadlineExceeded)
	})

	t.Run("If both timeout are specified, the first one which is expired should exist the method", func(t *testing.T) {
		// Timeout from call context first
		retrier := New(
			WithWaitDuration(1*time.Second),
			WithMaxDuration(500*time.Millisecond),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		before := time.Now()
		err := retrier.Do(ctx, func(ctx context.Context) error {
			return errors.New("test")
		})

		assert.Error(t, err)
		assert.WithinDuration(t, time.Now(), before, 100*time.Millisecond)
		require.IsType(t, RetryErr{}, err)
		assert.EqualValues(t, err.(RetryErr).Scope, ContextScope)
		assert.Equal(t, err.(RetryErr).Err, context.DeadlineExceeded)

		// Timeout from MaxDuration first
		retrier = New(
			WithWaitDuration(1*time.Second),
			WithMaxDuration(50*time.Millisecond),
		)

		ctx, cancel = context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		before = time.Now()
		err = retrier.Do(ctx, func(ctx context.Context) error {
			return errors.New("test")
		})

		assert.Error(t, err)
		assert.WithinDuration(t, time.Now(), before, 100*time.Millisecond)
		require.IsType(t, RetryErr{}, err)
		assert.EqualValues(t, err.(RetryErr).Scope, MaxDurationScope)
		assert.Equal(t, err.(RetryErr).Err, context.DeadlineExceeded)
	})
}
