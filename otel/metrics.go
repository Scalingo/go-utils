package otel

import (
	sdkmetric "go.opentelemetry.io/otel/metric"
)

type Meter interface {
	GetOriginalMeter() sdkmetric.Meter
}

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

func (m *meter) GetOriginalMeter() sdkmetric.Meter {
	return m.sdkMeter
}
