package otel

import (
	"context"
	"strconv"
	"time"

	"github.com/Scalingo/go-utils/env"
)

func ConfigFromEnv() Config {
	E := env.InitMapFromEnv(map[string]string{
		"OTEL_SERVICE_NAME":                     "",
		"OTEL_DEBUG":                            "false",
		"OTEL_EXPORTER_TYPE":                    "http",
		"OTEL_EXPORTER_OTLP_ENDPOINT":           "",
		"OTEL_EXPORTER_OTLP_TIMEOUT":            "",
		"OTEL_EXPORTER_OTLP_CERTIFICATE":        "",
		"OTEL_EXPORTER_OTLP_CLIENT_CERTIFICATE": "",
		"OTEL_EXPORTER_OTLP_CLIENT_KEY":         "",
		"OTEL_COLLECTION_INTERVAL":              "",
	})

	interval := 10 * time.Second
	if E["OTEL_COLLECTION_INTERVAL"] != "" {
		if secs, err := strconv.Atoi(E["OTEL_COLLECTION_INTERVAL"]); err == nil {
			interval = time.Duration(secs) * time.Second
		}
	}

	return Config{
		ServiceName:               E["OTEL_SERVICE_NAME"],
		Debug:                     E["OTEL_DEBUG"] == "true",
		ExporterType:              E["OTEL_EXPORTER_TYPE"],
		ExporterEndpoint:          E["OTEL_EXPORTER_OTLP_ENDPOINT"],
		ExporterTimeout:           E["OTEL_EXPORTER_OTLP_TIMEOUT"],
		ExporterCertificate:       E["OTEL_EXPORTER_OTLP_CERTIFICATE"],
		ExporterClientCertificate: E["OTEL_EXPORTER_OTLP_CLIENT_CERTIFICATE"],
		ExporterClientKey:         E["OTEL_EXPORTER_OTLP_CLIENT_KEY"],
		CollectionInterval:        interval,
	}
}

func InitFromEnv(ctx context.Context) (func(context.Context) error, error) {
	return InitGlobalWrapper(ctx, ConfigFromEnv())
}
