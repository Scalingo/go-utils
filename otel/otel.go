package otel

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	otelsdk "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// Config defines the Otel Wrapper configuration for metrics
type Config struct {
	ServiceName               string
	Debug                     bool
	ExporterType              string
	ExporterEndpoint          string
	ExporterTimeout           string
	ExporterCertificate       string
	ExporterClientCertificate string
	ExporterClientKey         string
	CollectionInterval        time.Duration
	MetricsReader             sdkmetric.Reader
	MetricsExporter           sdkmetric.Exporter
}

// OtelWrapper encapsulates OpenTelemetry MetricProvider and utilities
type OtelWrapper struct {
	meterProvider *sdkmetric.MeterProvider
	globalMeter   metric.Meter
	config        Config
}

var (
	globalWrapper *OtelWrapper

	// singletonShutdown holds the shutdown function once the meter provider is initialized.
	singletonShutdown func(context.Context) error

	// once ensures that initialization happens only once.
	globalOnce sync.Once

	// initErr captures any error encountered during initialization.
	initErr error
)

// InitSingleton initializes the MeterProvider as a singleton.
// Subsequent calls to this function will return the same shutdown function and error.
// The configuration is used only on the first call.
func InitGlobalWrapper(ctx context.Context, cfg Config) (func(context.Context) error, error) {
	globalOnce.Do(func() {
		singletonShutdown, initErr = New(ctx, cfg)
	})
	return singletonShutdown, initErr
}

func New(ctx context.Context, cfg Config) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up meter provider.
	meterProvider, err := newMeterProvider(ctx, cfg)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otelsdk.SetMeterProvider(meterProvider)

	globalWrapper = &OtelWrapper{
		meterProvider: meterProvider,
		globalMeter:   meterProvider.Meter(cfg.ServiceName),
		config:        cfg,
	}

	return
}

func newMeterProvider(ctx context.Context, cfg Config) (*sdkmetric.MeterProvider, error) {
	if cfg.ServiceName == "" {
		return nil, errors.New("ServiceName is required")
	}
	if cfg.MetricsExporter == nil {
		exporter, err := newMetricsExporter(ctx, cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to load exporter", err)
		}
		cfg.MetricsExporter = exporter
	}
	if cfg.MetricsReader == nil {
		cfg.MetricsReader = sdkmetric.NewPeriodicReader(cfg.MetricsExporter, sdkmetric.WithInterval(cfg.CollectionInterval))
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Initialize MeterProvider
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(cfg.MetricsReader),
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
		return nil, errors.New("invalid exporter type")
	}
}

// GetGlobalMeter returns the global Meter
func GetGlobalMeter() metric.Meter {
	if globalWrapper == nil {
		panic("Global OtelWrapper is not initialized")
	}
	return globalWrapper.globalMeter
}

// SendMetric wraps the meter to send a metric
func SendMetric(ctx context.Context, name string, value float64, description string) error {
	counter, err := GetGlobalMeter().Float64Counter(
		name,
		metric.WithDescription(description),
	)
	if err != nil {
		return fmt.Errorf("failed to create counter: %w", err)
	}
	counter.Add(ctx, value)
	return nil
}
