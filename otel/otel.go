package otel

import (
	"context"
	"os"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/kelseyhightower/envconfig"
	otelsdk "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"

	"github.com/Scalingo/go-utils/errors/v2"
)

type Config struct {
	ServiceName          string        `required:"true" split_words:"true"`
	ServiceInstanceId    string        `default:"" split_words:"true"`
	HostName             string        `default:"" split_words:"true"`
	Debug                bool          `default:"false"`
	SdkDisabled          bool          `default:"false" split_words:"true"`
	DebugPrettyPrint     bool          `default:"true" split_words:"true"`
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

	if cfg.SdkDisabled {
		return func(ctx context.Context) error {
			return nil
		}, nil
	}

	if cfg.ServiceName == "" {
		return nil, errors.New(ctx, "service name is required")
	}

	serviceInstanceID := cfg.ServiceInstanceId
	if serviceInstanceID == "" {
		serviceInstanceID = setServiceInstanceID()
	}

	hostName := cfg.HostName
	if hostName == "" {
		hostName = setHostname()
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
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceInstanceID(serviceInstanceID),
			semconv.HostName(hostName),
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
		if cfg.DebugPrettyPrint {
			return stdoutmetric.New(stdoutmetric.WithPrettyPrint())
		}
		return stdoutmetric.New()
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

func setHostname() string {
	hostName, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostName
}

func setServiceInstanceID() string {
	// To have "web-1" and other containers reported when is a Scalingo app
	containerInstance, isScalingoApp := os.LookupEnv("CONTAINER")
	if isScalingoApp {
		return containerInstance
	}

	// Otherwise generate a unique UUIDv4
	serviceInstanceID, err := uuid.NewV4()
	if err != nil {
		return "unknown"
	}
	return serviceInstanceID.String()
}
