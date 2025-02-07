package otel

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
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
			name: "minimal initialization",
			env: map[string]string{
				"OTEL_SERVICE_NAME": "test",
			},
		},
		{
			name:   "re-initialization to check singleton usage",
			reinit: true,
			env: map[string]string{
				"OTEL_SERVICE_NAME": "test",
			},
		},
		{
			name: "initialization with debug enabled",
			env: map[string]string{
				"OTEL_SERVICE_NAME": "test",
				"OTEL_DEBUG":        "true",
			},
		},
		{
			name: "initialization with grpc exporter",
			env: map[string]string{
				"OTEL_SERVICE_NAME":  "test",
				"OTEL_EXPORTER_TYPE": "grpc",
			},
		},
		{
			name: "initialization with collection interval at 30s",
			env: map[string]string{
				"OTEL_SERVICE_NAME":        "test",
				"OTEL_COLLECTION_INTERVAL": "30s",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Cleanup(func() {
				globalProviders = nil
				globalOnce = sync.Once{}
			})

			ctx := context.Background()

			if test.env != nil {
				for k, v := range test.env {
					t.Setenv(k, v)
				}
			}

			err := New(ctx)
			if test.expectError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), test.expectError)

				require.Nil(t, globalProviders)
			} else {
				require.NoError(t, err)

				require.NotNil(t, globalProviders)
				assert.NotNil(t, globalProviders.meterProvider)

				// Check when reinitializing the SDK
				if test.reinit {
					previousGlobalProviders := *globalProviders

					err = New(ctx)
					require.NoError(t, err)
					// Check that pointer are the same
					assert.Equal(t, previousGlobalProviders, *globalProviders)
				}
			}
		})
	}
}
