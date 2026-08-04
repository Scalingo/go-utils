package otel

import (
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"

	"github.com/Scalingo/go-utils/otel/otelmock"
)

var _ metric.Float64Counter = (*otelmock.MockFloat64Counter)(nil)
var _ metric.Float64Gauge = (*otelmock.MockFloat64Gauge)(nil)
var _ metric.Float64Histogram = (*otelmock.MockFloat64Histogram)(nil)
var _ metric.Float64ObservableCounter = (*otelmock.MockFloat64ObservableCounter)(nil)
var _ metric.Float64ObservableGauge = (*otelmock.MockFloat64ObservableGauge)(nil)
var _ metric.Float64ObservableUpDownCounter = (*otelmock.MockFloat64ObservableUpDownCounter)(nil)
var _ metric.Float64Observer = (*otelmock.MockFloat64Observer)(nil)
var _ metric.Float64UpDownCounter = (*otelmock.MockFloat64UpDownCounter)(nil)
var _ metric.Int64Counter = (*otelmock.MockInt64Counter)(nil)
var _ metric.Int64Gauge = (*otelmock.MockInt64Gauge)(nil)
var _ metric.Int64Histogram = (*otelmock.MockInt64Histogram)(nil)
var _ metric.Int64ObservableCounter = (*otelmock.MockInt64ObservableCounter)(nil)
var _ metric.Int64ObservableGauge = (*otelmock.MockInt64ObservableGauge)(nil)
var _ metric.Int64ObservableUpDownCounter = (*otelmock.MockInt64ObservableUpDownCounter)(nil)
var _ metric.Int64Observer = (*otelmock.MockInt64Observer)(nil)
var _ metric.Int64UpDownCounter = (*otelmock.MockInt64UpDownCounter)(nil)
var _ metric.Meter = (*otelmock.MockMeter)(nil)
var _ metric.MeterProvider = (*otelmock.MockMeterProvider)(nil)
var _ metric.Observer = (*otelmock.MockObserver)(nil)
var _ metric.Registration = (*otelmock.MockRegistration)(nil)
var _ sdkmetric.Exporter = (*otelmock.MockExporter)(nil)
