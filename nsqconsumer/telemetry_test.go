package nsqconsumer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	otelmock "github.com/Scalingo/go-utils/otel/otelmock"
	oteltest "github.com/Scalingo/go-utils/otel/oteltest"
)

func TestNewTelemetryCreatesInstruments(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	meterProvider := oteltest.InitMockMeterProvider(ctrl)
	mockMeter := otelmock.NewMockMeter(ctrl)

	meterProvider.EXPECT().Meter(telemetryInstrumentationName).Return(mockMeter)

	mockMeter.EXPECT().
		Int64Counter(messageCountMetricName, gomock.Any()).
		Return(otelmock.NewMockInt64Counter(ctrl), nil)
	mockMeter.EXPECT().
		Int64Counter(messageErrorsMetricName, gomock.Any()).
		Return(otelmock.NewMockInt64Counter(ctrl), nil)
	mockMeter.EXPECT().
		Float64Histogram(messageDurationMetricName, gomock.Any()).
		Return(otelmock.NewMockFloat64Histogram(ctrl), nil)

	telemetry, err := newTelemetry(t.Context())
	require.NoError(t, err)
	require.NotNil(t, telemetry)
}

func TestTelemetryRecordRecordsMetrics(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		messageType         string
		err                 error
		expectErrCountCalls int
		expectedType        string
	}{
		{
			name:                "success",
			messageType:         "event",
			expectErrCountCalls: 0,
			expectedType:        "event",
		},
		{
			name:                "error",
			messageType:         "event",
			err:                 errors.New("boom"),
			expectErrCountCalls: 1,
			expectedType:        "event",
		},
		{
			name:                "unknown type",
			messageType:         "",
			expectErrCountCalls: 0,
			expectedType:        "unknown",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			messagesCounter := otelmock.NewMockInt64Counter(ctrl)
			messageErrorsCounter := otelmock.NewMockInt64Counter(ctrl)
			messageDuration := otelmock.NewMockFloat64Histogram(ctrl)

			telemetry := &telemetry{
				messagesCounter:      messagesCounter,
				messageErrorsCounter: messageErrorsCounter,
				messageDuration:      messageDuration,
			}

			startedAt := time.Now().Add(-50 * time.Millisecond)
			topic := "my-topic"
			channel := "my-channel"

			messagesCounter.EXPECT().
				Add(gomock.Any(), int64(1), gomock.Any()).
				Do(func(_ context.Context, _ int64, opts ...metric.AddOption) {
					assertTelemetryAttributes(t, opts, topic, channel, test.expectedType)
				})

			messageErrorsCounter.EXPECT().
				Add(gomock.Any(), int64(1), gomock.Any()).
				Times(test.expectErrCountCalls).
				Do(func(_ context.Context, _ int64, opts ...metric.AddOption) {
					assertTelemetryAttributes(t, opts, topic, channel, test.expectedType)
				})

			messageDuration.EXPECT().
				Record(gomock.Any(), gomock.Any(), gomock.Any()).
				Do(func(_ context.Context, value float64, opts ...metric.RecordOption) {
					require.GreaterOrEqual(t, value, 0.0)
					assertTelemetryAttributesForRecord(t, opts, topic, channel, test.expectedType)
				})

			telemetry.record(t.Context(), startedAt, topic, channel, test.messageType, test.err)
		})
	}
}

func assertTelemetryAttributes(t *testing.T, opts []metric.AddOption, topic, channel, messageType string) {
	t.Helper()

	config := metric.NewAddConfig(opts)
	attrs := config.Attributes()

	assertAttributeValue(t, &attrs, topicAttributeKey, topic)
	assertAttributeValue(t, &attrs, channelAttributeKey, channel)
	assertAttributeValue(t, &attrs, messageTypeAttributeKey, messageType)
}

func assertTelemetryAttributesForRecord(t *testing.T, opts []metric.RecordOption, topic, channel, messageType string) {
	t.Helper()

	config := metric.NewRecordConfig(opts)
	attrs := config.Attributes()

	assertAttributeValue(t, &attrs, topicAttributeKey, topic)
	assertAttributeValue(t, &attrs, channelAttributeKey, channel)
	assertAttributeValue(t, &attrs, messageTypeAttributeKey, messageType)
}

func assertAttributeValue(t *testing.T, attrs *attribute.Set, key, expected string) {
	t.Helper()

	value, ok := attrs.Value(attribute.Key(key))
	require.True(t, ok, "expected %q attribute to be set", key)
	require.Equal(t, expected, value.AsString(), "expected %q attribute to be %q", key, expected)
}
