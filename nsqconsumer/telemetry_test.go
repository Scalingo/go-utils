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
		Float64Histogram(messageDurationMetricName, gomock.Any()).
		Return(otelmock.NewMockFloat64Histogram(ctrl), nil)

	telemetry, err := newTelemetry(t.Context())
	require.NoError(t, err)
	require.NotNil(t, telemetry)
}

func TestTelemetryRecordRecordsMetrics(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		messageType    string
		err            error
		expectedType   string
		expectedStatus string
	}{
		{
			name:           "success",
			messageType:    "event",
			expectedType:   "event",
			expectedStatus: statusSuccess,
		},
		{
			name:           "error",
			messageType:    "event",
			err:            errors.New("boom"),
			expectedType:   "event",
			expectedStatus: statusError,
		},
		{
			name:           "unknown type",
			messageType:    "",
			expectedType:   "unknown",
			expectedStatus: statusSuccess,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			messageDuration := otelmock.NewMockFloat64Histogram(ctrl)

			telemetry := &telemetry{
				messageDuration: messageDuration,
			}

			startedAt := time.Now().Add(-50 * time.Millisecond)
			topic := "my-topic"
			channel := "my-channel"

			messageDuration.EXPECT().
				Record(gomock.Any(), gomock.Any(), gomock.Any()).
				Do(func(_ context.Context, value float64, opts ...metric.RecordOption) {
					require.GreaterOrEqual(t, value, 0.0)
					assertTelemetryAttributesForRecord(t, opts, topic, channel, test.expectedType, test.expectedStatus)
				})

			telemetry.record(t.Context(), startedAt, topic, channel, test.messageType, test.err)
		})
	}
}

func assertTelemetryAttributesForRecord(t *testing.T, opts []metric.RecordOption, topic, channel, messageType, status string) {
	t.Helper()

	config := metric.NewRecordConfig(opts)
	attrs := config.Attributes()

	assertAttributeValue(t, &attrs, topicAttributeKey, topic)
	assertAttributeValue(t, &attrs, channelAttributeKey, channel)
	assertAttributeValue(t, &attrs, messageTypeAttributeKey, messageType)
	assertAttributeValue(t, &attrs, statusAttributeKey, status)
}

func assertAttributeValue(t *testing.T, attrs *attribute.Set, key, expected string) {
	t.Helper()

	value, ok := attrs.Value(attribute.Key(key))
	require.True(t, ok, "expected %q attribute to be set", key)
	require.Equal(t, expected, value.AsString(), "expected %q attribute to be %q", key, expected)
}
