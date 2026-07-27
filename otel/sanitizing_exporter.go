package otel

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"

	"github.com/Scalingo/go-utils/logger"
)

type sanitizingExporter struct {
	exporter sdkmetric.Exporter
	log      logrus.FieldLogger
}

func newSanitizingExporter(ctx context.Context, exporter sdkmetric.Exporter) sdkmetric.Exporter {
	return sanitizingExporter{
		exporter: exporter,
		log:      logger.Get(ctx),
	}
}

func (e sanitizingExporter) Temporality(kind sdkmetric.InstrumentKind) metricdata.Temporality {
	return e.exporter.Temporality(kind)
}

func (e sanitizingExporter) Aggregation(kind sdkmetric.InstrumentKind) sdkmetric.Aggregation {
	return e.exporter.Aggregation(kind)
}

func (e sanitizingExporter) Export(ctx context.Context, metrics *metricdata.ResourceMetrics) error {
	e.sanitizeResourceMetrics(metrics)
	return e.exporter.Export(ctx, metrics)
}

func (e sanitizingExporter) ForceFlush(ctx context.Context) error {
	return e.exporter.ForceFlush(ctx)
}

func (e sanitizingExporter) Shutdown(ctx context.Context) error {
	return e.exporter.Shutdown(ctx)
}

func (e sanitizingExporter) sanitizeResourceMetrics(metrics *metricdata.ResourceMetrics) {
	if metrics == nil {
		return
	}

	if metrics.Resource != nil {
		attrs, changed := e.sanitizeAttributes("resource", "", metrics.Resource.Attributes())
		if changed {
			metrics.Resource = resource.NewWithAttributes(metrics.Resource.SchemaURL(), attrs...)
		}
	}

	for scopeIdx := range metrics.ScopeMetrics {
		scopeMetrics := &metrics.ScopeMetrics[scopeIdx]
		for metricIdx := range scopeMetrics.Metrics {
			e.sanitizeMetric(&scopeMetrics.Metrics[metricIdx])
		}
	}
}

func (e sanitizingExporter) sanitizeMetric(metric *metricdata.Metrics) {
	switch data := metric.Data.(type) {
	case metricdata.Gauge[int64]:
		sanitizeDataPoints(e, metric.Name, data.DataPoints)
		metric.Data = data
	case metricdata.Gauge[float64]:
		sanitizeDataPoints(e, metric.Name, data.DataPoints)
		metric.Data = data
	case metricdata.Sum[int64]:
		sanitizeDataPoints(e, metric.Name, data.DataPoints)
		metric.Data = data
	case metricdata.Sum[float64]:
		sanitizeDataPoints(e, metric.Name, data.DataPoints)
		metric.Data = data
	case metricdata.Histogram[int64]:
		sanitizeHistogramDataPoints(e, metric.Name, data.DataPoints)
		metric.Data = data
	case metricdata.Histogram[float64]:
		sanitizeHistogramDataPoints(e, metric.Name, data.DataPoints)
		metric.Data = data
	case metricdata.ExponentialHistogram[int64]:
		sanitizeExponentialHistogramDataPoints(e, metric.Name, data.DataPoints)
		metric.Data = data
	case metricdata.ExponentialHistogram[float64]:
		sanitizeExponentialHistogramDataPoints(e, metric.Name, data.DataPoints)
		metric.Data = data
	case metricdata.Summary:
		sanitizeSummaryDataPoints(e, metric.Name, data.DataPoints)
		metric.Data = data
	}
}

func sanitizeDataPoints[N int64 | float64](exporter sanitizingExporter, metricName string, dataPoints []metricdata.DataPoint[N]) {
	for idx := range dataPoints {
		attrs, changed := exporter.sanitizeAttributes("datapoint", metricName, dataPoints[idx].Attributes.ToSlice())
		if changed {
			dataPoints[idx].Attributes = attribute.NewSet(attrs...)
		}
	}
}

func sanitizeHistogramDataPoints[N int64 | float64](exporter sanitizingExporter, metricName string, dataPoints []metricdata.HistogramDataPoint[N]) {
	for idx := range dataPoints {
		attrs, changed := exporter.sanitizeAttributes("datapoint", metricName, dataPoints[idx].Attributes.ToSlice())
		if changed {
			dataPoints[idx].Attributes = attribute.NewSet(attrs...)
		}
	}
}

func sanitizeExponentialHistogramDataPoints[N int64 | float64](exporter sanitizingExporter, metricName string, dataPoints []metricdata.ExponentialHistogramDataPoint[N]) {
	for idx := range dataPoints {
		attrs, changed := exporter.sanitizeAttributes("datapoint", metricName, dataPoints[idx].Attributes.ToSlice())
		if changed {
			dataPoints[idx].Attributes = attribute.NewSet(attrs...)
		}
	}
}

func sanitizeSummaryDataPoints(exporter sanitizingExporter, metricName string, dataPoints []metricdata.SummaryDataPoint) {
	for idx := range dataPoints {
		attrs, changed := exporter.sanitizeAttributes("datapoint", metricName, dataPoints[idx].Attributes.ToSlice())
		if changed {
			dataPoints[idx].Attributes = attribute.NewSet(attrs...)
		}
	}
}

func (e sanitizingExporter) sanitizeAttributes(location, metricName string, attrs []attribute.KeyValue) ([]attribute.KeyValue, bool) {
	sanitizedAttrs := make([]attribute.KeyValue, 0, len(attrs))
	changed := false

	for _, attr := range attrs {
		if !isExportableAttribute(attr) {
			changed = true
			e.logDroppedAttribute(location, metricName, attr)
			continue
		}

		sanitizedAttrs = append(sanitizedAttrs, attr)
	}

	return sanitizedAttrs, changed
}

func isExportableAttribute(attr attribute.KeyValue) bool {
	if string(attr.Key) == "" {
		return false
	}

	if attr.Value.Type() == attribute.EMPTY {
		return false
	}

	return attr.Value.String() != ""
}

func (e sanitizingExporter) logDroppedAttribute(location, metricName string, attr attribute.KeyValue) {
	fields := logrus.Fields{
		"attribute_key":      string(attr.Key),
		"attribute_location": location,
		"attribute_type":     attr.Value.Type().String(),
		"attribute_value":    attr.Value.String(),
	}
	if metricName != "" {
		fields["metric_name"] = metricName
	}

	e.log.WithFields(fields).Error("OpenTelemetry attribute with empty name or value removed")
}
