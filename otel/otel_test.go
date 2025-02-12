package otel

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name        string
		reinit      bool
		env         map[string]string
		expectError string
	}{
		{
			name:        "initialization without service_name defined should result in error",
			expectError: "required key OTEL_SERVICE_NAME missing value",
		},
		{
			name:        "initialization without exporter endpoint defined should result in error",
			expectError: "otlp endpoint is required",
			env: map[string]string{
				"OTEL_SERVICE_NAME": "test",
			},
		},
		{
			name: "minimal initialization",
			env: map[string]string{
				"OTEL_SERVICE_NAME": "test",
				"OTEL_DEBUG":        "true",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()

			if test.env != nil {
				for k, v := range test.env {
					t.Setenv(k, v)
				}
			}

			shutdown, err := Init(ctx)

			if test.expectError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), test.expectError)

				require.Nil(t, shutdown)
			} else {
				require.NoError(t, err)
				require.NotNil(t, shutdown)

				t.Cleanup(func() {
					require.NoError(t, shutdown(ctx))
				})
			}
		})
	}
}
