# Package `otel` v0.6.2

## Usage

### Collect a metric synchronously (with a counter instrument)

See the directory [docs/examples/int64-sync-counter](docs/examples/int64-sync-counter) for a complete example.

### Collect a metric asynchronously (with a gauge instrument)

See the directory [docs/examples/int64-async-gauge](docs/examples/int64-async-gauge) for a complete example.

## Development of this package

### Generate mocks

Run the following command to generate mocks for the `otel` package inside the `otelmock` package:
```bash
gomock_generator
```

> [!IMPORTANT]
> And then, please make sure that on each mock file, that each struct embed the original interface.

For example, the `gomock` library underneath generates the following code for the `metric.Int64Gauge` interface:
```go
// MockInt64Gauge is a mock of Int64Gauge interface.
type MockInt64Gauge struct {
	ctrl     *gomock.Controller
	recorder *MockInt64GaugeMockRecorder
}
```

Please make sure that the `MockInt64Gauge` struct embeds the `metric.Int64Gauge` interface like this, otherwise the generated mocks will not be usable:
```go
type MockInt64Gauge struct {
	metric.Int64Gauge

	ctrl     *gomock.Controller
	recorder *MockInt64GaugeMockRecorder
}
```

This is actually a missing feature of `gomock` that is not yet implemented, see [this issue](https://github.com/uber-go/mock/issues/64).

## Updating this package

> [!IMPORTANT]
> The version of `semconv` (aka semantic conventions) (that was imported in the `otel.go` file) should always match with the version that is released inside OpenTelemetry packages.

Otherwise this error will be generated when initializing the package:
```plaintext
fail to init otel sdk: create resource: conflicting Schema URL: https://opentelemetry.io/schemas/1.34.0 and https://opentelemetry.io/schemas/1.26.0
```

Currently there is no matrix of compatibility between OpenTelemetry Go SDK versions and the `semconv` version,
so you should always use the latest `semconv` version that is compatible with the OpenTelemetry Go SDK version you are using.

Tips: Take the one from the latest OpenTelemetry Go SDK release changelog.
Example: For OTEL [v1.37.0](https://github.com/open-telemetry/opentelemetry-go/releases/tag/v1.37.0), the `semconv` version is `v1.34.0` according to the changelog (https://github.com/open-telemetry/opentelemetry-go/pull/6832/files).
