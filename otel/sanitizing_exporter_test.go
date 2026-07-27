package otel

import (
	"context"
	"io"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.uber.org/mock/gomock"

	"github.com/Scalingo/go-utils/otel/otelmock"
)

func TestSanitizingExporter_ExportRemovesEmptyAttributes(t *testing.T) {
	ctrl := gomock.NewController(t)
	baseExporter := otelmock.NewMockExporter(ctrl)
	log := logrus.New()
	log.SetOutput(io.Discard)

	exporter := sanitizingExporter{
		exporter: baseExporter,
		log:      log,
	}
	baseExporter.EXPECT().Export(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, exported *metricdata.ResourceMetrics) error {
			require.NotNil(t, exported)
			require.ElementsMatch(t, []attribute.KeyValue{
				attribute.String("resource.valid", "yes"),
				attribute.Int("resource.zero", 0),
				attribute.Bool("resource.false", false),
			}, exported.Resource.Attributes())

			gauge, ok := exported.ScopeMetrics[0].Metrics[0].Data.(metricdata.Gauge[int64])
			require.True(t, ok)
			require.ElementsMatch(t, []attribute.KeyValue{
				attribute.String("datapoint.valid", "yes"),
				attribute.Int("datapoint.zero", 0),
				attribute.Bool("datapoint.false", false),
			}, gauge.DataPoints[0].Attributes.ToSlice())
			return nil
		},
	)

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

	err := exporter.Export(t.Context(), metrics)
	require.NoError(t, err)
}
