package otel

import (
	"context"
	"io"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
)

func TestSanitizingExporter_ExportRemovesEmptyAttributes(t *testing.T) {
	baseExporter := &capturingExporter{}
	log := logrus.New()
	log.SetOutput(io.Discard)

	exporter := sanitizingExporter{
		exporter: baseExporter,
		log:      log,
	}

	metrics := &metricdata.ResourceMetrics{
		Resource: resource.NewWithAttributes(
			"https://example.com/schema",
			attribute.String("", "empty-name"),
			attribute.String("resource.empty", ""),
			attribute.KeyValue{Key: "resource.unset"},
			attribute.String("resource.valid", "yes"),
			attribute.Int("resource.zero", 0),
			attribute.Bool("resource.false", false),
		),
		ScopeMetrics: []metricdata.ScopeMetrics{
			{
				Metrics: []metricdata.Metrics{
					{
						Name: "test_gauge",
						Data: metricdata.Gauge[int64]{
							DataPoints: []metricdata.DataPoint[int64]{
								{
									Attributes: attribute.NewSet(
										attribute.String("", "empty-name"),
										attribute.String("datapoint.empty", ""),
										attribute.KeyValue{Key: "datapoint.unset"},
										attribute.String("datapoint.valid", "yes"),
										attribute.Int("datapoint.zero", 0),
										attribute.Bool("datapoint.false", false),
									),
									Value: 42,
								},
							},
						},
					},
				},
			},
		},
	}

	err := exporter.Export(context.Background(), metrics)
	require.NoError(t, err)

	require.NotNil(t, baseExporter.exported)
	require.ElementsMatch(t, []attribute.KeyValue{
		attribute.String("resource.valid", "yes"),
		attribute.Int("resource.zero", 0),
		attribute.Bool("resource.false", false),
	}, baseExporter.exported.Resource.Attributes())

	gauge, ok := baseExporter.exported.ScopeMetrics[0].Metrics[0].Data.(metricdata.Gauge[int64])
	require.True(t, ok)
	require.ElementsMatch(t, []attribute.KeyValue{
		attribute.String("datapoint.valid", "yes"),
		attribute.Int("datapoint.zero", 0),
		attribute.Bool("datapoint.false", false),
	}, gauge.DataPoints[0].Attributes.ToSlice())
}

type capturingExporter struct {
	exported *metricdata.ResourceMetrics
}

func (e *capturingExporter) Temporality(_ sdkmetric.InstrumentKind) metricdata.Temporality {
	return metricdata.CumulativeTemporality
}

func (e *capturingExporter) Aggregation(kind sdkmetric.InstrumentKind) sdkmetric.Aggregation {
	return sdkmetric.DefaultAggregationSelector(kind)
}

func (e *capturingExporter) Export(_ context.Context, metrics *metricdata.ResourceMetrics) error {
	e.exported = metrics
	return nil
}

func (e *capturingExporter) ForceFlush(context.Context) error {
	return nil
}

func (e *capturingExporter) Shutdown(context.Context) error {
	return nil
}
