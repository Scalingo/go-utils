package concurrency

import (
	"context"
	"fmt"

	"github.com/Scalingo/go-utils/logger"
)

func ExampleParallelWorker() {
	ctx := context.Background()
	res := make(chan int)
	over := make(chan struct{})

	slice := []string{"this", "is", "a", "test"}
	worker := NewParallelWorker(10, func() {
		close(res)
		// Ensure that our logics has finished reading all the values before
		// quitting Stop()
		<-over
	})

	for _, item := range slice {
		worker.Perform(func(ctx context.Context, item string) func() {
			return func() {
				log := logger.Get(ctx)
				log.Infof("performing worker for %v", item)
				res <- len(item)
			}
		}(ctx, item))
	}

	totalLength := 0
	go func() {
		defer close(over)
		for length := range res {
			totalLength += length
		}
	}()

	// Wait the end of the length computation
	worker.Stop()

	fmt.Println("totalLength:", totalLength)
	// Output:
	// totalLength: 11
}
