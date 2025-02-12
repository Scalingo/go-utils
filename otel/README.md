# Package `otel` v0.1.0 

## Usage

```go
package main

import (
    "context"

	otelsdk "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/metric"

    "github.com/Scalingo/go-utils/otel"
)

func main() {
    ctx := context.Background()

    // Initialize OpenTelemetry SDK
    shutdown, err := otel.Init(ctx)
    if err != nil {
        fmt.Printf("init otel: %v\n", err)
        return
    }
    // Handle collection of metrics properly when service shut down
    defer shutdown(ctx)

	// Create a meter
	meter := otelsdk.Meter("deployment")

	// Create an instrument, based on the meter previously created
	deploymentCount, err := meter.Int64Counter("deployment_count", sdkmetric.WithDescription("Number of deployments"))
	if err != nil {
		fmt.Printf("instrument creation failed: %v\n", err)
		return
	}

	// Create measurements on the instrument
	deploymentCount.Add(ctx, 10, sdkmetric.WithAttributes(attribute.String("app_id", "caaaefb0-dcaa-4866-83d2-b581228169d8")))
	deploymentCount.Add(ctx, 42, sdkmetric.WithAttributes(attribute.String("app_id", "caaaefb0-dcaa-4866-83d2-b581228169d8")))
}
```
