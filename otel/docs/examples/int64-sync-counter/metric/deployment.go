package metric

import (
	"context"
	"fmt"

	otelsdk "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/metric"
)

func WriteDeploymentMetric() bool {
	ctx := context.Background()

	deploymentMeter := otelsdk.Meter("deployment")
	fmt.Println("Meter() was called")

	deploymentCount, err := deploymentMeter.Int64Counter(
		"deployment_count",
		sdkmetric.WithDescription("Number of deployments"),
	)
	if err != nil {
		fmt.Printf("instrument creation failed: %v\n", err)
		return false
	}
	fmt.Println("Int64Counter() was called")

	appID := "3c8b0ca8-98b5-4fd3-a7d5-68ad974badb4"

	// Attribute like Heroku: https://opentelemetry.io/docs/specs/semconv/registry/attributes/heroku/#heroku-attributes
	deploymentCount.Add(ctx, 10, sdkmetric.WithAttributes(attribute.String("scalingo.app.id", appID)))
	deploymentCount.Add(ctx, 42, sdkmetric.WithAttributes(attribute.String("scalingo.app.id", appID)))

	return true
}
