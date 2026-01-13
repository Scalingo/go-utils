package nsqconsumer

import (
	"context"
	"time"

	otelsdk "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/Scalingo/go-utils/errors/v3"
)

type telemetry struct {
	messageDuration metric.Float64Histogram
}

const (
	telemetryInstrumentationName = "scalingo.nsq_consumer"
	messageDurationMetricName    = "scalingo.nsq_consumer.message.duration"
)

const (
	topicAttributeKey       = "scalingo.nsq.topic"
	channelAttributeKey     = "scalingo.nsq.channel"
	messageTypeAttributeKey = "scalingo.nsq.message_type"
	statusAttributeKey      = "scalingo.nsq.status"
	unknownMessageType      = "unknown"
	statusSuccess           = "success"
	statusError             = "error"
)

func newTelemetry(ctx context.Context) (*telemetry, error) {
	meter := otelsdk.Meter(telemetryInstrumentationName)

	messageDuration, err := meter.Float64Histogram(
		messageDurationMetricName,
		metric.WithDescription("NSQ message handling duration in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "create message duration histogram")
	}

	return &telemetry{
		messageDuration: messageDuration,
	}, nil
}

func (t *telemetry) record(ctx context.Context, startedAt time.Time, topic, channel, messageType string, err error) {
	if messageType == "" {
		messageType = unknownMessageType
	}
	status := statusSuccess
	if err != nil {
		status = statusError
	}
	attrs := metric.WithAttributes(
		attribute.String(topicAttributeKey, topic),
		attribute.String(channelAttributeKey, channel),
		attribute.String(messageTypeAttributeKey, messageType),
		attribute.String(statusAttributeKey, status),
	)

	t.messageDuration.Record(ctx, time.Since(startedAt).Seconds(), attrs)
}
