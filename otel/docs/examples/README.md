# Examples

Examples can be run with these following commands inside each example directory:

- Run the service once and show the OpenTelemetry payload generated:
```bash
OTEL_DEBUG=true OTEL_SERVICE_NAME="test" go run ./main.go
```

- Run unit tests:
```bash
go test -race ./...
```
