package otel

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	otelsdk "go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func TestInit(t *testing.T) {
	minimalValidEnv := map[string]string{
		"GO_ENV":                      "test",
		"OTEL_SERVICE_NAME":           "test",
		"OTEL_EXPORTER_OTLP_ENDPOINT": "http://localhost:4317",
		// OTEL_EXPORTER_OTLP_METRICS_TIMEOUT is set to avoid to wait 10 seconds in the test
		"OTEL_EXPORTER_OTLP_METRICS_TIMEOUT": "1", // 1 millisecond
	}
	defaultShutdownError := "failed to upload metrics: exporter export timeout"

	tests := []struct {
		name                string
		reinit              bool
		env                 map[string]string
		expectInitSkipped   bool
		expectShutdownError string
		opts                []InitOpt
	}{
		{
			name:              "initialization without service_name and exporter endpoint should skip init",
			expectInitSkipped: true,
			env: map[string]string{
				"GO_ENV": "test",
			},
		},
		{
			name:              "initialization without exporter endpoint should skip init",
			expectInitSkipped: true,
			env: map[string]string{
				"GO_ENV":            "test",
				"OTEL_SERVICE_NAME": "test",
			},
		},
		{
			name:              "initialization with SDK disabled should skip init",
			expectInitSkipped: true,
			env: map[string]string{
				"GO_ENV":            "test",
				"OTEL_SDK_DISABLED": "true",
			},
		},
		{
			name: "minimal initialization",
			// expected error in the case of the unit test, due to endpoint that doesn't respond
			expectShutdownError: defaultShutdownError,
			env:                 minimalValidEnv,
		}, {
			name:                "with additional global attributes",
			expectShutdownError: defaultShutdownError,
			env:                 minimalValidEnv,
			opts: []InitOpt{
				WithServiceVersionAttribute("v1.0.0"),
				func(opts *initDefaultOptions) {
					require.Len(t, opts.defaultAttributes, 4)
					assert.Equal(t, opts.defaultAttributes[3].Key, semconv.ServiceVersionKey)
					assert.Equal(t, opts.defaultAttributes[3].Value, "v1.0.0")
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := t.Context()

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
				meterProvider := otelsdk.GetMeterProvider()
				require.NotSame(t, initialMeterProvider, meterProvider)
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

func TestSetTLSConfig(t *testing.T) {
	dir := t.TempDir()

	// Generate valid CA + client certs for success case
	caPath, clientCertPath, clientKeyPath := generateTestCerts(t, dir)

	tests := []struct {
		name          string
		cfg           Config
		expectErr     string
		expectSuccess bool
	}{
		{
			name: "missing CA path",
			cfg: Config{
				ExporterOtlpCertificate:       "",
				ExporterOtlpClientCertificate: clientCertPath,
				ExporterOtlpClientKey:         clientKeyPath,
			},
			expectErr: "CA certificate must be set",
		},
		{
			name: "missing client cert and key",
			cfg: Config{
				ExporterOtlpCertificate:       caPath,
				ExporterOtlpClientCertificate: "",
				ExporterOtlpClientKey:         "",
			},
			expectErr: "client certificate and client key must be set",
		},
		{
			name: "unreadable CA file",
			cfg: Config{
				ExporterOtlpCertificate:       filepath.Join(dir, "does-not-exist.pem"),
				ExporterOtlpClientCertificate: clientCertPath,
				ExporterOtlpClientKey:         clientKeyPath,
			},
			expectErr: "read CA file",
		},
		{
			name: "invalid CA PEM",
			cfg: Config{
				ExporterOtlpCertificate:       writeFile(t, dir, "bad-ca.pem", []byte("not a pem")),
				ExporterOtlpClientCertificate: clientCertPath,
				ExporterOtlpClientKey:         clientKeyPath,
			},
			expectErr: "append CA PEM to cert pool",
		},
		{
			name: "invalid client key pair",
			cfg: Config{
				ExporterOtlpCertificate:       caPath,
				ExporterOtlpClientCertificate: writeFile(t, dir, "bad-client.crt", []byte("nope")),
				ExporterOtlpClientKey:         writeFile(t, dir, "bad-client.key", []byte("nope")),
			},
			expectErr: "load client key pair",
		},
		{
			name: "success with valid CA and client certs",
			cfg: Config{
				ExporterOtlpCertificate:       caPath,
				ExporterOtlpClientCertificate: clientCertPath,
				ExporterOtlpClientKey:         clientKeyPath,
			},
			expectSuccess: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := t.Context()
			tlsCfg, err := setTLSConfig(ctx, &test.cfg)

			if test.expectErr != "" {
				require.ErrorContains(t, err, test.expectErr)
				require.Nil(t, tlsCfg)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, tlsCfg)
			require.Len(t, tlsCfg.Certificates, 1)
		})
	}
}

// --- helpers ---

// generateTestCerts creates a CA certificate and a client certificate signed by that CA.
// It writes files to dir and returns paths: (caPEMPath, clientCertPath, clientKeyPath).
func generateTestCerts(t *testing.T, dir string) (string, string, string) {
	t.Helper()

	// CA key & cert
	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	caTpl := &x509.Certificate{
		SerialNumber:          bigSerial(t),
		Subject:               pkix.Name{CommonName: "test-ca"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLenZero:        true,
	}
	caDER, err := x509.CreateCertificate(rand.Reader, caTpl, caTpl, &caKey.PublicKey, caKey)
	require.NoError(t, err)

	caPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caDER})
	caPEMPath := filepath.Join(dir, "ca.pem")
	require.NoError(t, os.WriteFile(caPEMPath, caPEM, 0o600))

	// Client key & cert signed by CA
	clientKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	clientTpl := &x509.Certificate{
		SerialNumber: bigSerial(t),
		Subject:      pkix.Name{CommonName: "test-client"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(24 * time.Hour),
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
	}
	caCert, err := x509.ParseCertificate(caDER)
	require.NoError(t, err)
	clientDER, err := x509.CreateCertificate(rand.Reader, clientTpl, caCert, &clientKey.PublicKey, caKey)
	require.NoError(t, err)

	clientCertPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: clientDER})
	clientKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(clientKey)})

	clientCertPath := filepath.Join(dir, "client.crt")
	clientKeyPath := filepath.Join(dir, "client.key")
	require.NoError(t, os.WriteFile(clientCertPath, clientCertPEM, 0o600))
	require.NoError(t, os.WriteFile(clientKeyPath, clientKeyPEM, 0o600))

	return caPEMPath, clientCertPath, clientKeyPath
}

func bigSerial(t *testing.T) *big.Int {
	t.Helper()
	return big.NewInt(time.Now().UnixNano())
}

func writeFile(t *testing.T, dir, name string, content []byte) string {
	t.Helper()
	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, content, 0o600))
	return path
}
