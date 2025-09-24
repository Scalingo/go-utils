# Changelog

## To be released

## v0.8.0

* feat: Add function options to Init() method in order to inject global attributes
* feat: Add WithServiceVersionAttribute helper

## v0.7.0

* build(deps): Update dependencies
  * github.com/Scalingo/go-utils/errors/v3 from `3.0.0` to `3.1.0`
  * go.opentelemetry.io/otel/* from `1.37.0` to `1.38.0`
  * google.golang.org/grpc from `1.74.2` to `1.75.0`
  * github.com/grpc-ecosystem/grpc-gateway/v2 from `2.27.1` to `2.27.2`
  * golang.org/x/net from `0.42.0` to `0.43.0`
  * golang.org/x/sys from `0.34.0` to `0.35.0`
  * golang.org/x/text from `0.27.0` to `0.28.0`
  * google.golang.org/protobuf from `1.36.6` to `1.36.8`
* fix: Switch resource to schemaless to avoid schema/semconv mismatch issues

## v0.6.2

* build: Bump github.com/Scalingo/go-utils/errors from v2.5.1 to v3.0.0

## v0.6.1

* fix: Correctly load TLS certificates again from env vars since it broke in v0.6.0

## v0.6.0

* feat: Switch Exporter Type to gRPC by default, instead of HTTP
* fix: Enforce TLS connections, use allowlisted ciphers and TLS v1.2 at least in production/staging environments
* fix: Return early if SDK is disabled, don't try to load configuration
* fix: Don't load OpenTelemetry SDK if config is incorrect or if anything in initialization fails

## v0.5.0

* fix: Update OpenTelemetry packages from `v1.34.0` to `v1.37.0` and update `semconv` from `v1.26.0` to `v1.34.0`

## v0.4.0

* feat: Add missing resource attributes: `service_instance_id` and `host_name` based on OTEL specs

## v0.3.1

* chore(go): Corrective bump - Go version regression from 1.24.3 to 1.24

## v0.3.0

* chore(go): Upgrade to Go 1.24

## v0.2.1

* feat: Add `OTEL_SDK_DISABLED` env var to completely disable the OpenTelemetry SDK (default: false)

## v0.2.0

* feat: Add `OTEL_DEBUG_PRETTY_PRINT` env var to enable or disable pretty printing of OpenTelemetry payloads (default: enabled)
* fix: Remove `OTEL_COLLECTION_INTERVAL` env var in favor of [OTEL_METRIC_EXPORT_INTERVAL](https://github.com/open-telemetry/opentelemetry-go/blob/a9cbc3d8dec7be22c7d3691ca1755f25c1702a1d/sdk/metric/env.go#L17)
* docs: Updated README.md and added examples with implemented code and unit tests
* test: Add `oteltest` package that contains test helpers for unit testing
* test: Add `otelmock` package that contains mocks required for unit testing

## v0.1.0

* Initial version of the `otel` package to initialize the OpenTelemetry SDK to send metrics
