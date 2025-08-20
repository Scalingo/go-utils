package main

import (
	"context"
	"fmt"

	"github.com/Scalingo/go-utils/docs/examples/int64-async-gauge/metric"
	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/go-utils/otel"
)

func main() {
	fmt.Println("start of program")

	ctx := context.Background()
	log := logger.Get(ctx)

	shutdown := otel.Init(ctx)
	defer func() {
		err := shutdown()
		if err != nil {
			log.WithError(err).Error("Shutdown OpenTelemetry")
		}
	}()

	metric.RegisterDatabaseAsyncGauge()

	fmt.Println("end of program")
}
