package metric

import (
	"context"
	"fmt"

	otelsdk "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/metric"
)

func RegisterDatabaseAsyncGauge() bool {
	meter := otelsdk.Meter("database")

	_, err := meter.Int64ObservableGauge(
		"database_count",
		sdkmetric.WithDescription("Number of databases"),
		sdkmetric.WithInt64Callback(func(ctx context.Context, observer sdkmetric.Int64Observer) error {
			var databaseCount int64
			databaseCount = 42

			observer.Observe(databaseCount, sdkmetric.WithAttributes(attribute.String("env", "prod")))

			return nil
		}),
	)
	if err != nil {
		fmt.Printf("observable gauge creation failed: %v\n", err)
		return false
	}
	return true
}
