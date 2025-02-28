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

func (m *MockMeterWrapper) Int64Gauge(name string, options ...metric.Int64GaugeOption) (metric.Int64Gauge, error) {
	rawGauge, err := m.meter.Int64Gauge(name, options...)
	return &MockInt64GaugeWrapper{gauge: rawGauge}, err
}

func (m *MockMeterWrapper) Int64Histogram(name string, options ...metric.Int64HistogramOption) (metric.Int64Histogram, error) {
	rawHistogram, err := m.meter.Int64Histogram(name, options...)
	return &MockInt64HistogramWrapper{histogram: rawHistogram}, err
}

func (m *MockMeterWrapper) Int64UpDownCounter(name string, options ...metric.Int64UpDownCounterOption) (metric.Int64UpDownCounter, error) {
	rawUpDownCounter, err := m.meter.Int64UpDownCounter(name, options...)
	return &MockInt64UpDownCounterWrapper{upDownCounter: rawUpDownCounter}, err
}
