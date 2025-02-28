package oteltest

import (
	"context"

	"go.opentelemetry.io/otel/metric"
)

// MockInt64CounterWrapper is a wrapper around a metric.Int64Counter that allows for mocking.
type MockInt64CounterWrapper struct {
	metric.Int64Counter

	counter metric.Int64Counter
}

func (m *MockInt64CounterWrapper) Add(ctx context.Context, value int64, options ...metric.AddOption) {
	m.counter.Add(ctx, value, options...)
}

// MockInt64GaugeWrapper is a wrapper around a metric.Int64Gauge that allows for mocking.
type MockInt64GaugeWrapper struct {
	metric.Int64Gauge

	gauge metric.Int64Gauge
}

func (m *MockInt64GaugeWrapper) Record(ctx context.Context, value int64, options ...metric.RecordOption) {
	m.gauge.Record(ctx, value, options...)
}

// MockInt64HistogramWrapper is a wrapper around a metric.Int64Histogram that allows for mocking.
type MockInt64HistogramWrapper struct {
	metric.Int64Histogram

	histogram metric.Int64Histogram
}

func (m *MockInt64HistogramWrapper) Record(ctx context.Context, value int64, options ...metric.RecordOption) {
	m.histogram.Record(ctx, value, options...)
}

type MockInt64UpDownCounterWrapper struct {
	metric.Int64UpDownCounter

	upDownCounter metric.Int64UpDownCounter
}

func (m *MockInt64UpDownCounterWrapper) Add(ctx context.Context, value int64, options ...metric.AddOption) {
	m.upDownCounter.Add(ctx, value, options...)
}
