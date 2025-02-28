package oteltest

import (
	"context"

	"go.opentelemetry.io/otel/metric"
)

type MockInt64CounterWrapper struct {
	metric.Int64Counter

	counter metric.Int64Counter
}

func (m *MockInt64CounterWrapper) Add(ctx context.Context, value int64, options ...metric.AddOption) {
	m.counter.Add(ctx, value, options...)
}
