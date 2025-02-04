# How to use

```go

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Scalingo/go-utils/otel"
)

func main() {
	ctx := context.Background()
	// Set up OpenTelemetry SDK
	// Method 1: Using the env
	otelShutdown, err := otel.InitFromEnv(ctx)
	if err != nil {
		return
	}
	// Method 2: configuring it ourselves
	// otelConfig := otel.ConfigFromEnv()
	// otelConfig.Exporter = ...
	// otelShutdown, err := otel.InitGlobalWrapper(ctx, otelConfig)

	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	// 4. Get Global Meter and Send Metrics
	if err := globalWrapperSendMetrics(ctx); err != nil {
		log.Fatalf("failed to send metrics: %v", err)
	}

	fmt.Println("Metrics successfully sent!")
}

// globalWrapperSendMetrics sends an example metric
func globalWrapperSendMetrics(ctx context.Context) error {
	// Example of sending a counter metric
	err := otel.SendMetric(ctx, "example.counter", 1.0, "An example counter metric")
	if err != nil {
		return err
	}

	// Simulating some load
	time.Sleep(2 * time.Second)

	err = otel.SendMetric(ctx, "example.counter", 2.5, "Another counter increment")
	if err != nil {
		return err
	}

	return nil
}
```
