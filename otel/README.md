# Package `otel` v0.7.0

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
