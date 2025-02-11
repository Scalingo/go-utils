package otel

import (
	sdkmetric "go.opentelemetry.io/otel/metric"
)

type Meter interface{}

type meter struct {
	name     string
	sdkMeter sdkmetric.Meter
}

func GetMeter(name string) Meter {
	return &meter{
		name:     name,
		sdkMeter: globalProviders.meterProvider.Meter(name),
	}
}
