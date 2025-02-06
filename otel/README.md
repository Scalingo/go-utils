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

  // Initialize Otel
  err := otel.New(ctx)
  if err != nil {
    return
  }

  // Handle shutdown properly so nothing leaks.
  defer otel.Shutdown(ctx)
}
```
