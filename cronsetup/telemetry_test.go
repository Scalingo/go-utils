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

	etcdcron "github.com/Scalingo/go-etcd-cron"
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
		Int64Counter("scalingo.etcd_cron.runs_total", gomock.Any()).
		Return(otelmock.NewMockInt64Counter(ctrl), nil)
	mockMeter.EXPECT().
		Int64Counter("scalingo.etcd_cron.run_errors_total", gomock.Any()).
		Return(otelmock.NewMockInt64Counter(ctrl), nil)
	mockMeter.EXPECT().
		Float64Histogram("scalingo.etcd_cron.runs_duration_milliseconds", gomock.Any()).
		Return(otelmock.NewMockFloat64Histogram(ctrl), nil)

	telemetry, err := newTelemetry(context.Background())
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
		minDurationMs  float64
		maxDurationMs  float64
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
			minDurationMs:  80,
			maxDurationMs:  120,
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

			job := etcdcron.Job{
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
					if test.minDurationMs > 0 || test.maxDurationMs > 0 {
						require.GreaterOrEqual(t, value, test.minDurationMs, "expected duration >= %.2fms", test.minDurationMs)
						require.LessOrEqual(t, value, test.maxDurationMs, "expected duration <= %.2fms", test.maxDurationMs)
					}
					assertJobAttributeForRecord(t, opts, job.Name)
				})

			wrapped := telemetry.wrapJob(job)
			if test.expectError {
				require.Error(t, wrapped.Func(context.Background()))
			} else {
				require.NoError(t, wrapped.Func(context.Background()))
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
