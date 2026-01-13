package nsqproducer

import (
	"context"
	"time"

	otelsdk "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type telemetry struct {
	publishCounter       metric.Int64Counter
	publishErrorsCounter metric.Int64Counter
	publishDuration      metric.Float64Histogram
}

const (
	telemetryInstrumentationName = "scalingo.nsq_producer"
	publishCountMetricName       = "scalingo.nsq_producer.publish.count"
	publishErrorsMetricName      = "scalingo.nsq_producer.publish.errors"
	publishDurationMetricName    = "scalingo.nsq_producer.publish.duration"
	unknownMessageType           = "unknown"
)

const (
	topicAttributeKey       = "scalingo.nsq.topic"
	messageTypeAttributeKey = "scalingo.nsq.message_type"
	publishTypeAttributeKey = "scalingo.nsq.publish_type"
)

const (
	publishTypeImmediate = "immediate"
	publishTypeDeferred  = "deferred"
)

func newTelemetry() (*telemetry, error) {
	meter := otelsdk.Meter(telemetryInstrumentationName)

	publishCounter, err := meter.Int64Counter(
		publishCountMetricName,
		metric.WithDescription("Number of NSQ messages published"),
	)
	if err != nil {
		return nil, err
	}

	publishErrorsCounter, err := meter.Int64Counter(
		publishErrorsMetricName,
		metric.WithDescription("Number of NSQ message publish errors"),
	)
	if err != nil {
		return nil, err
	}

	publishDuration, err := meter.Float64Histogram(
		publishDurationMetricName,
		metric.WithDescription("NSQ message publish duration in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	return &telemetry{
		publishCounter:       publishCounter,
		publishErrorsCounter: publishErrorsCounter,
		publishDuration:      publishDuration,
	}, nil
}

func (t *telemetry) record(ctx context.Context, startedAt time.Time, topic, messageType, publishType string, err error) {
	if messageType == "" {
		messageType = unknownMessageType
	}
	attrs := metric.WithAttributes(
		attribute.String(topicAttributeKey, topic),
		attribute.String(messageTypeAttributeKey, messageType),
		attribute.String(publishTypeAttributeKey, publishType),
	)

	t.publishCounter.Add(ctx, 1, attrs)
	if err != nil {
		t.publishErrorsCounter.Add(ctx, 1, attrs)
	}
	t.publishDuration.Record(ctx, time.Since(startedAt).Seconds(), attrs)
}
