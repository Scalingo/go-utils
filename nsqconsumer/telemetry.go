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
	messagesCounter      metric.Int64Counter
	messageErrorsCounter metric.Int64Counter
	messageDuration      metric.Float64Histogram
}

const (
	telemetryInstrumentationName = "scalingo.nsq_consumer"
	messageCountMetricName       = "scalingo.nsq_consumer.message.count"
	messageErrorsMetricName      = "scalingo.nsq_consumer.message.errors"
	messageDurationMetricName    = "scalingo.nsq_consumer.message.duration"
)

const (
	topicAttributeKey       = "scalingo.nsq.topic"
	channelAttributeKey     = "scalingo.nsq.channel"
	messageTypeAttributeKey = "scalingo.nsq.message_type"
	unknownMessageType      = "unknown"
)

func newTelemetry(ctx context.Context) (*telemetry, error) {
	meter := otelsdk.Meter(telemetryInstrumentationName)

	messagesCounter, err := meter.Int64Counter(
		messageCountMetricName,
		metric.WithDescription("Number of NSQ messages handled"),
	)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "create messages counter")
	}

	messageErrorsCounter, err := meter.Int64Counter(
		messageErrorsMetricName,
		metric.WithDescription("Number of NSQ messages handled with errors"),
	)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "create message errors counter")
	}

	messageDuration, err := meter.Float64Histogram(
		messageDurationMetricName,
		metric.WithDescription("NSQ message handling duration in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "create message duration histogram")
	}

	return &telemetry{
		messagesCounter:      messagesCounter,
		messageErrorsCounter: messageErrorsCounter,
		messageDuration:      messageDuration,
	}, nil
}

func (t *telemetry) record(ctx context.Context, startedAt time.Time, topic, channel, messageType string, err error) {
	if messageType == "" {
		messageType = unknownMessageType
	}
	attrs := metric.WithAttributes(
		attribute.String(topicAttributeKey, topic),
		attribute.String(channelAttributeKey, channel),
		attribute.String(messageTypeAttributeKey, messageType),
	)

	t.messagesCounter.Add(ctx, 1, attrs)
	if err != nil {
		t.messageErrorsCounter.Add(ctx, 1, attrs)
	}
	t.messageDuration.Record(ctx, time.Since(startedAt).Seconds(), attrs)
}
