package otel

import (
	"context"

	sdkmetric "go.opentelemetry.io/otel/metric"

	"github.com/Scalingo/go-utils/errors/v2"
)

type Meter interface {
	NewInstrument(ctx context.Context, kind InstrumentKind, name string, description string, unit string) (Instrument, error)
}

type meter struct {
	name     string
	sdkMeter sdkmetric.Meter
}

type InstrumentKind string

const (
	Counter            InstrumentKind = "counter"
	AsyncCounter       InstrumentKind = "async_counter"
	Histogram          InstrumentKind = "histogram"
	Gauge              InstrumentKind = "gauge"
	AsyncGauge         InstrumentKind = "async_gauge"
	UpDownCounter      InstrumentKind = "up_down_counter"
	AsyncUpDownCounter InstrumentKind = "async_up_down_counter"
)

type Instrument interface {
	Add(ctx context.Context, value float64)
	Remove(ctx context.Context, value float64)
	Record(ctx context.Context, value float64)
}

type instrument struct {
	instrumentKind InstrumentKind
	name           string
	description    string
	unit           string
	counter        sdkmetric.Float64Counter
	histogram      sdkmetric.Float64Histogram
	gauge          sdkmetric.Float64Gauge
	upDownCounter  sdkmetric.Float64UpDownCounter
}

func GetMeter(name string) Meter {
	return &meter{
		name:     name,
		sdkMeter: globalProviders.meterProvider.Meter(name),
	}
}

func (m *meter) NewInstrument(ctx context.Context, kind InstrumentKind, name string, description string, unit string) (Instrument, error) {
	var err error
	instrument := &instrument{
		instrumentKind: kind,
		name:           name,
		description:    description,
		unit:           unit,
	}
	switch kind {
	case Counter:
		instrument.counter, err = m.sdkMeter.Float64Counter(
			name,
			sdkmetric.WithDescription(description),
			sdkmetric.WithUnit(unit))
		if err != nil {
			return nil, err
		}
	case Histogram:
		instrument.histogram, err = m.sdkMeter.Float64Histogram(
			name,
			sdkmetric.WithDescription(description),
			sdkmetric.WithUnit(unit))
		if err != nil {
			return nil, err
		}
	case Gauge:
		instrument.gauge, err = m.sdkMeter.Float64Gauge(
			name,
			sdkmetric.WithDescription(description),
			sdkmetric.WithUnit(unit))
		if err != nil {
			return nil, err
		}
	case UpDownCounter:
		instrument.upDownCounter, err = m.sdkMeter.Float64UpDownCounter(
			name,
			sdkmetric.WithDescription(description),
			sdkmetric.WithUnit(unit))
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.Newf(ctx, "unknown instrument kind: %v", kind)
	}
	return instrument, nil
}

func (i *instrument) Add(ctx context.Context, value float64) {
	switch i.instrumentKind {
	case Counter:
		i.counter.Add(ctx, value)
	case Histogram:
		i.histogram.Record(ctx, value)
	case Gauge:
		i.gauge.Record(ctx, value)
	case UpDownCounter:
		i.upDownCounter.Add(ctx, value)
	}
}

func (i *instrument) Remove(ctx context.Context, value float64) {
	i.Add(ctx, -value)
}

func (i *instrument) Record(ctx context.Context, value float64) {
	i.Add(ctx, value)
}
