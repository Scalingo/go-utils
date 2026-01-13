package nsqproducer

import (
	"context"
	"time"

	otelsdk "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type telemetry struct {
	publishDuration metric.Float64Histogram
}

const (
	telemetryInstrumentationName = "scalingo.nsq_producer"
	publishDurationMetricName    = "scalingo.nsq_producer.publish.duration"
	unknownMessageType           = "unknown"
)

const (
	topicAttributeKey       = "scalingo.nsq.topic"
	messageTypeAttributeKey = "scalingo.nsq.message_type"
	publishTypeAttributeKey = "scalingo.nsq.publish_type"
	statusAttributeKey      = "scalingo.nsq.status"
)

const (
	publishTypeImmediate = "immediate"
	publishTypeDeferred  = "deferred"
)

func newTelemetry() (*telemetry, error) {
	meter := otelsdk.Meter(telemetryInstrumentationName)

	publishDuration, err := meter.Float64Histogram(
		publishDurationMetricName,
		metric.WithDescription("NSQ message publish duration in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	return &telemetry{
		publishDuration: publishDuration,
	}, nil
}

func (t *telemetry) record(ctx context.Context, startedAt time.Time, topic, messageType, publishType string, err error) {
	if messageType == "" {
		messageType = unknownMessageType
	}
	status := "success"
	if err != nil {
		status = "error"
	}
	attrs := metric.WithAttributes(
		attribute.String(topicAttributeKey, topic),
		attribute.String(messageTypeAttributeKey, messageType),
		attribute.String(publishTypeAttributeKey, publishType),
		attribute.String(statusAttributeKey, status),
	)

	t.publishDuration.Record(ctx, time.Since(startedAt).Seconds(), attrs)
}
