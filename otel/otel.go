package otel

import (
	"context"
	"sync"
	"time"

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

// OtelProviders encapsulates OpenTelemetry providers and utilities
type OtelProviders struct {
	meterProvider *sdkmetric.MeterProvider
	shutdownFunc  func(context.Context) error
	config        Config
}

var (
	globalProvider *OtelProviders

	// once ensures that initialization happens only once.
	globalOnce sync.Once
)

// Initializes the OtelWrapper as a singleton.
// The configuration is used only on the first call.
func New(ctx context.Context, cfg Config) error {
	var err error
	globalOnce.Do(func() {
		globalProvider, err = setupProviders(ctx, cfg)
	})
	return err
}

func setupProviders(ctx context.Context, cfg Config) (*OtelProviders, error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown := func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			shutdownErr := fn(ctx)
			if shutdownErr != nil {
				err = errors.Wrapf(ctx, err, "failed to shutdown provider: %v", shutdownErr)
			}
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) error {
		err := shutdown(ctx)
		if err != nil {
			return errors.Wrapf(ctx, inErr, "failed to shutdown otel providers: %v", err)
		}
		return inErr
	}

	// Set up meter provider.
	meterProvider, err := newMeterProvider(ctx, cfg)
	if err != nil {
		return nil, handleErr(err)
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otelsdk.SetMeterProvider(meterProvider)

	return &OtelProviders{
		meterProvider: meterProvider,
		shutdownFunc:  shutdown,
		config:        cfg,
	}, nil
}

func newMeterProvider(ctx context.Context, cfg Config) (*sdkmetric.MeterProvider, error) {
	if cfg.ServiceName == "" {
		return nil, errors.New(ctx, "ServiceName is required")
	}
	metricsExporter, err := newMetricsExporter(ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "failed to load exporter")
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
		return nil, errors.Wrap(ctx, err, "failed to create resource")
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
	if err := globalProvider.shutdownFunc(ctx); err != nil {
		return errors.Wrap(ctx, err, "failed to shutdown otel providers")
	}
	return nil
}
