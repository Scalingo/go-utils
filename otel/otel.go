package otel

import (
	"context"
	"sync"
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
	ServiceName        string        `required:"true" split_words:"true"`
	Debug              bool          `default:"false"`
	ExporterType       string        `default:"http" split_words:"true"`
	CollectionInterval time.Duration `default:"10s" split_words:"true"`
}

// Providers encapsulates OpenTelemetry providers and utilities
type Providers struct {
	meterProvider *sdkmetric.MeterProvider
	config        Config
}

var (
	globalProviders *Providers

	// once ensures that initialization happens only once.
	globalOnce sync.Once
)

// Initializes the Providers as a singleton.
func New(ctx context.Context) error {
	var err error
	globalOnce.Do(func() {
		// Get Otel configuration from environment
		var cfg Config
		err = envconfig.Process("OTEL", &cfg)
		if err != nil {
			return
		}
		globalProviders, err = setupProviders(ctx, cfg)
	})
	return err
}

func setupProviders(ctx context.Context, cfg Config) (*Providers, error) {
	// Set up meter provider.
	meterProvider, err := newMeterProvider(ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "create meter provider")
	}

	otelsdk.SetMeterProvider(meterProvider)

	return &Providers{
		meterProvider: meterProvider,
		config:        cfg,
	}, nil
}

func newMeterProvider(ctx context.Context, cfg Config) (*sdkmetric.MeterProvider, error) {
	if cfg.ServiceName == "" {
		return nil, errors.New(ctx, "ServiceName is required")
	}
	metricsExporter, err := newMetricsExporter(ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "load exporter")
	}
	metricsReader := sdkmetric.NewPeriodicReader(metricsExporter, sdkmetric.WithInterval(cfg.CollectionInterval))

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
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(metricsReader),
		sdkmetric.WithResource(res),
	)

	return provider, nil
}

func newMetricsExporter(ctx context.Context, cfg Config) (sdkmetric.Exporter, error) {
	if cfg.Debug {
		return stdoutmetric.New(stdoutmetric.WithPrettyPrint())
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

// Gracefully shuts down the providers
func Shutdown(ctx context.Context) error {
	if globalProviders.meterProvider != nil {
		err := globalProviders.meterProvider.Shutdown(ctx)
		if err != nil {
			return errors.Wrap(ctx, err, "shutdown meter provider")
		}
	}
	return nil
}
