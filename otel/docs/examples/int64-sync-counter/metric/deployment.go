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

	deploymentCount.Add(ctx, 10, sdkmetric.WithAttributes(attribute.String("env", "prod")))
	deploymentCount.Add(ctx, 42, sdkmetric.WithAttributes(attribute.String("env", "staging")))

	return true
}
