package cronsetup

import (
	"context"
	"time"

	otelsdk "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/Scalingo/go-utils/cronsetup/internal/cron"
	"github.com/Scalingo/go-utils/errors/v3"
)

type telemetry struct {
	runsCounter      metric.Int64Counter
	runErrorsCounter metric.Int64Counter
	runsDuration     metric.Float64Histogram
}

// jobNameAttributeKey captures the executed job name.
const jobNameAttributeKey = "scalingo.etcd_cron.job_name"

func newTelemetry(ctx context.Context) (*telemetry, error) {
	meter := otelsdk.Meter("scalingo.etcd_cron")

	runsCounter, err := meter.Int64Counter(
		"scalingo.etcd_cron.runs_total",
		metric.WithDescription("Number of cron job runs"),
	)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "create runs counter")
	}

	runErrorsCounter, err := meter.Int64Counter(
		"scalingo.etcd_cron.run_errors_total",
		metric.WithDescription("Number of cron job runs with errors"),
	)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "create run errors counter")
	}

	runsDuration, err := meter.Float64Histogram(
		"scalingo.etcd_cron.runs_duration_milliseconds",
		metric.WithDescription("Cron job execution duration in milliseconds"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "create runs duration histogram")
	}

	return &telemetry{
		runsCounter:      runsCounter,
		runErrorsCounter: runErrorsCounter,
		runsDuration:     runsDuration,
	}, nil
}

func (t *telemetry) wrapJob(job cron.Job) cron.Job {
	jobAttributes := metric.WithAttributes(attribute.String(jobNameAttributeKey, job.Name))
	originalFunc := job.Func

	job.Func = func(ctx context.Context) error {
		startedAt := time.Now()
		err := originalFunc(ctx)

		t.runsCounter.Add(ctx, 1, jobAttributes)
		if err != nil {
			t.runErrorsCounter.Add(ctx, 1, jobAttributes)
		}
		t.runsDuration.Record(ctx, float64(time.Since(startedAt).Milliseconds()), jobAttributes)

		return err
	}

	return job
}
