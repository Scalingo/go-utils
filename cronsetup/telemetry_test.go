package cronsetup

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

	meterProvider.EXPECT().Meter("scalingo.etcd_cron").Return(mockMeter)

	mockMeter.EXPECT().
		Int64Counter("scalingo.etcd_cron.run.count", gomock.Any()).
		Return(otelmock.NewMockInt64Counter(ctrl), nil)
	mockMeter.EXPECT().
		Int64Counter("scalingo.etcd_cron.run.errors", gomock.Any()).
		Return(otelmock.NewMockInt64Counter(ctrl), nil)
	mockMeter.EXPECT().
		Float64Histogram("scalingo.etcd_cron.run.duration", gomock.Any()).
		Return(otelmock.NewMockFloat64Histogram(ctrl), nil)

	telemetry, err := newTelemetry(t.Context())
	require.NoError(t, err)
	require.NotNil(t, telemetry)
}

func TestTelemetryWrapJobRecordsMetrics(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		jobName        string
		jobFunc        func(context.Context) error
		expectError    bool
		expectErrCount int
		minDurationSec float64
		maxDurationSec float64
	}{
		{
			name:           "success",
			jobName:        "my job",
			jobFunc:        func(context.Context) error { return nil },
			expectError:    false,
			expectErrCount: 0,
		},
		{
			name:           "error",
			jobName:        "failing job",
			jobFunc:        func(context.Context) error { return errors.New("boom") },
			expectError:    true,
			expectErrCount: 1,
		},
		{
			name:           "duration around 100ms",
			jobName:        "slow job",
			jobFunc:        func(context.Context) error { time.Sleep(100 * time.Millisecond); return nil },
			expectError:    false,
			expectErrCount: 0,
			minDurationSec: 0.08,
			maxDurationSec: 0.12,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			runsCounter := otelmock.NewMockInt64Counter(ctrl)
			runErrorsCounter := otelmock.NewMockInt64Counter(ctrl)
			runsDuration := otelmock.NewMockFloat64Histogram(ctrl)

			telemetry := &telemetry{
				runsCounter:      runsCounter,
				runErrorsCounter: runErrorsCounter,
				runsDuration:     runsDuration,
			}

			job := Job{
				Name: test.jobName,
				Func: test.jobFunc,
			}

			runsCounter.EXPECT().
				Add(gomock.Any(), int64(1), gomock.Any()).
				Do(func(_ context.Context, _ int64, opts ...metric.AddOption) {
					assertJobAttribute(t, opts, job.Name)
				})

			if test.expectErrCount == 0 {
				runErrorsCounter.EXPECT().
					Add(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			} else {
				runErrorsCounter.EXPECT().
					Add(gomock.Any(), int64(test.expectErrCount), gomock.Any()).
					Do(func(_ context.Context, _ int64, opts ...metric.AddOption) {
						assertJobAttribute(t, opts, job.Name)
					})
			}

			runsDuration.EXPECT().
				Record(gomock.Any(), gomock.Any(), gomock.Any()).
				Do(func(_ context.Context, value float64, opts ...metric.RecordOption) {
					require.GreaterOrEqual(t, value, 0.0, "expected non-negative duration")
					if test.minDurationSec > 0 || test.maxDurationSec > 0 {
						require.GreaterOrEqual(t, value, test.minDurationSec, "expected duration >= %.2fs", test.minDurationSec)
						require.LessOrEqual(t, value, test.maxDurationSec, "expected duration <= %.2fs", test.maxDurationSec)
					}
					assertJobAttributeForRecord(t, opts, job.Name)
				})

			wrapped := telemetry.wrapJob(job)
			if test.expectError {
				require.Error(t, wrapped.Func(t.Context()))
			} else {
				require.NoError(t, wrapped.Func(t.Context()))
			}
		})
	}
}

func assertJobAttribute(t *testing.T, opts []metric.AddOption, jobName string) {
	t.Helper()

	config := metric.NewAddConfig(opts)
	attrs := config.Attributes()
	value, ok := (&attrs).Value(attribute.Key(jobNameAttributeKey))
	require.True(t, ok, "expected %q attribute to be set", jobNameAttributeKey)
	require.Equal(t, jobName, value.AsString(), "expected %q attribute to be %q", jobNameAttributeKey, jobName)
}

func assertJobAttributeForRecord(t *testing.T, opts []metric.RecordOption, jobName string) {
	t.Helper()

	config := metric.NewRecordConfig(opts)
	attrs := config.Attributes()
	value, ok := (&attrs).Value(attribute.Key(jobNameAttributeKey))
	require.True(t, ok, "expected %q attribute to be set", jobNameAttributeKey)
	require.Equal(t, jobName, value.AsString(), "expected %q attribute to be %q", jobNameAttributeKey, jobName)
}
