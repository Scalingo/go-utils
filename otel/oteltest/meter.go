package oteltest

import (
	"go.opentelemetry.io/otel/metric"
)

type MockMeterWrapper struct {
	metric.Meter

	meter metric.Meter
}

func (m *MockMeterWrapper) Int64Counter(name string, options ...metric.Int64CounterOption) (metric.Int64Counter, error) {
	rawCounter, err := m.meter.Int64Counter(name, options...)
	return &MockInt64CounterWrapper{counter: rawCounter}, err
}
