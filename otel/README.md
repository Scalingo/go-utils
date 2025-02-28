# Package `otel` v0.1.0 

## Usage

### Collect a metric synchronously (with a counter instrument)

```go
package main

import (
    "context"
	"fmt"

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

### Collect a metric asynchronously (with a gauge instrument)

```go
package main

import (
    "context"
	"fmt"

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

	meter := otelsdk.Meter("database")

	_, err = meter.Int64ObservableGauge(
		"database_count",
		sdkmetric.WithDescription("Number of databases"),
		sdkmetric.WithInt64Callback(func(ctx context.Context, observer sdkmetric.Int64Observer) error {
			var databaseCount int64
			databaseCount = 42

			observer.Observe(
				databaseCount, sdkmetric.WithAttributes(
					attribute.String("app_id", "caaaefb0-dcaa-4866-83d2-b581228169d8"),
                ),
            )

			return nil
		}),
	)
	if err != nil {
		fmt.Printf("observable gauge creation failed: %v\n", err)
	}
}
```

## Development of this package

### Generate mocks

Run the following command to generate mocks for the `otel` package inside the `otelmock` package:
```bash
gomock_generator
```

> [!IMPORTANT]
> And then, please make sure that on each mock file, that each struct embed the original interface.

For example, the `gomock` library underneath generates the following code for the `metric.Int64Gauge` interface:
```go
// MockInt64Gauge is a mock of Int64Gauge interface.
type MockInt64Gauge struct {
	ctrl     *gomock.Controller
	recorder *MockInt64GaugeMockRecorder
}
```

Please make sure that the `MockInt64Gauge` struct embeds the `metric.Int64Gauge` interface like this, otherwise the generated mocks will not be usable:
```go
type MockInt64Gauge struct {
	metric.Int64Gauge

	ctrl     *gomock.Controller
	recorder *MockInt64GaugeMockRecorder
}
```

This is actually a missing feature of `gomock` that is not yet implemented, see [this issue](https://github.com/uber-go/mock/issues/64).
