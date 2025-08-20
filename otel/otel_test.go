package otel

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	otelsdk "go.opentelemetry.io/otel"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name                string
		reinit              bool
		env                 map[string]string
		expectInitSkipped   bool
		expectShutdownError string
	}{
		{
			name:              "initialization without service_name and exporter endpoint should skip init",
			expectInitSkipped: true,
		},
		{
			name:              "initialization without exporter endpoint should skip init",
			expectInitSkipped: true,
			env: map[string]string{
				"OTEL_SERVICE_NAME": "test",
			},
		},
		{
			name:              "initialization with SDK disabled should skip init",
			expectInitSkipped: true,
			env: map[string]string{
				"OTEL_SDK_DISABLED": "true",
			},
		},
		{
			name: "minimal initialization",
			// expected error in the case of the unit test, due to endpoint that doesn't respond
			expectShutdownError: "failed to upload metrics: exporter export timeout",
			env: map[string]string{
				"OTEL_SERVICE_NAME":           "test",
				"OTEL_EXPORTER_OTLP_ENDPOINT": "http://localhost:4317",
				// OTEL_EXPORTER_OTLP_METRICS_TIMEOUT is set to avoid to wait 10 seconds in the test
				"OTEL_EXPORTER_OTLP_METRICS_TIMEOUT": "1", // 1 millisecond
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()

			initialMeterProvider := otelsdk.GetMeterProvider()

			if test.env != nil {
				for k, v := range test.env {
					t.Setenv(k, v)
				}
			}

			shutdown := Init(ctx)
			require.NotNil(t, shutdown)

			if test.expectInitSkipped {
				// Should be the same object as before initialization
				require.Same(t, initialMeterProvider, otelsdk.GetMeterProvider())
			} else {
				require.NotSame(t, initialMeterProvider, otelsdk.GetMeterProvider())
			}

			t.Cleanup(func() {
				err := shutdown()

				if test.expectShutdownError != "" {
					require.Error(t, err, test.expectShutdownError)
					return
				}
				require.NoError(t, err)
			})
		})
	}
}
