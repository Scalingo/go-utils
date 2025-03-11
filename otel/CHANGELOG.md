# Changelog

## To be released

## v0.2.0

- feat: Add `OTEL_DEBUG_PRETTY_PRINT` env var to enable or disable pretty printing of OpenTelemetry payloads (default: enabled)
- fix: Remove `OTEL_COLLECTION_INTERVAL` env var in favor of [OTEL_METRIC_EXPORT_INTERVAL](https://github.com/open-telemetry/opentelemetry-go/blob/a9cbc3d8dec7be22c7d3691ca1755f25c1702a1d/sdk/metric/env.go#L17)
- docs: Updated README.md and added examples with implemented code and unit tests
- test: Add `oteltest` package that contains test helpers for unit testing
- test: Add `otelmock` package that contains mocks required for unit testing

## v0.1.0

- Initial version of the `otel` package to initialize the OpenTelemetry SDK to send metrics
