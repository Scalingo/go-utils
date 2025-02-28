package otel

import (
	"context"
	"time"

	"github.com/kelseyhightower/envconfig"
	otelsdk "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	"github.com/Scalingo/go-utils/errors/v2"
)

type Config struct {
	ServiceName          string        `required:"true" split_words:"true"`
	Debug                bool          `default:"false"`
	ExporterType         string        `default:"http" split_words:"true"`
	ExporterOtlpEndpoint string        `default:"" split_words:"true"`
	MetricExportInterval time.Duration `default:"10s" split_words:"true"`
}

func Init(ctx context.Context) (func(context.Context) error, error) {
	// Get Otel configuration from environment
	var cfg Config
	err := envconfig.Process("OTEL", &cfg)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "load configuration")
	}

	if cfg.ServiceName == "" {
		return nil, errors.New(ctx, "service name is required")
	}
	metricsExporter, err := newMetricsExporter(ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "load exporter")
	}

	metricsReader := sdkmetric.NewPeriodicReader(
		metricsExporter,
		sdkmetric.WithInterval(cfg.MetricExportInterval),
	)

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName),
		),
	)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "create resource")
	}

	// Initialize MeterProvider
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(metricsReader),
		sdkmetric.WithResource(res),
	)

	// Set the MeterProvider in the OTEL SDK global in order to access it globally
	otelsdk.SetMeterProvider(meterProvider)

	return func(ctx context.Context) error {
		if meterProvider != nil {
			return meterProvider.Shutdown(ctx)
		}
		return nil
	}, nil
}

func newMetricsExporter(ctx context.Context, cfg Config) (sdkmetric.Exporter, error) {
	if cfg.Debug {
		return stdoutmetric.New(stdoutmetric.WithPrettyPrint())
	}
	if cfg.ExporterOtlpEndpoint == "" {
		return nil, errors.New(ctx, "otlp endpoint is required")
	}
	switch cfg.ExporterType {
	case "http":
		return otlpmetrichttp.New(ctx)
	case "grpc":
		return otlpmetricgrpc.New(ctx)
	default:
		return nil, errors.New(ctx, "invalid exporter type")
	}
}
