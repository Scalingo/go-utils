package main

import (
	"context"
	"fmt"

	"github.com/Scalingo/go-utils/docs/examples/int64-async-gauge/metric"
	"github.com/Scalingo/go-utils/otel"
)

func main() {
	fmt.Println("start of program")

	ctx := context.Background()

	shutdown, err := otel.Init(ctx)
	if err != nil {
		fmt.Printf("init otel: %v\n", err)
		return
	}
	defer shutdown(ctx)

	metric.RegisterDatabaseAsyncGauge()

	fmt.Println("end of program")
}
