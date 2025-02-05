# How to use

```go

package main

import (
  "context"

  "github.com/Scalingo/go-utils/otel"
  "github.com/kelseyhightower/envconfig"
)

func main() {
  ctx := context.Background()

  // Get Otel configuration from environment
  var otelConfig otel.Config
  err := envconfig.Process("OTEL", &otelConfig)
  if err != nil {
    return
  }

  // Initialize Otel
  err := otel.New(ctx, otelConfig)
  if err != nil {
    return
  }

  // Handle shutdown properly so nothing leaks.
  defer otel.Shutdown(ctx)
}
```
