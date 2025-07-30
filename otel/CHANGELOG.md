# Changelog

## To be released

* fix: Enforce TLS connections, use allowlisted ciphers and TLS v1.2 at least in production/staging environments

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
