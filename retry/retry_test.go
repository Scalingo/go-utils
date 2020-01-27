package retry

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
		retrier := New(WithWaitDuration(1 * time.Millisecond))
		tries := 0
		err := retrier.Do(context.Background(), func(ctx context.Context) error {
			tries++
			if tries == 2 {
				return nil
			}
			return fmt.Errorf("Error attempt %v", tries)
		})

		assert.NoError(t, err)
		assert.Equal(t, tries, 2)
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
}
