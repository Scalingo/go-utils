# Logger

This package will provide you a generic way to handle logging.

## Configuration

This step is optional, by default the library uses Debug log level and a text formatter.

```go
logger.SetConfig(logrus.DebugLevel, new(logrus.TextFormatter)) // This will set the logger level and type.

log := logger.Default() // Return a default logger
```

## Context

The logger can be passed in a context so it can retain fields.

```go
func main() {
  log := logger.Default().WithFields(logrus.Fields{"caller": "main"})
  add(logger.ToCtx(context.Background(), log), 1, 2)
}

def add(ctx context.Context, a, b int) int {
  log := logger.Get(ctx)
  log.Info("Starting add operation")

  log.WithField("operation", "add")
  do(logger.ToCtx(ctx, log), a,b, func(a,b int)int{return a+b})
}

def do(ctx context.Context, a,b int, op fun(int, int)int) {
  log := logger.Get(ctx)
  log.Info("Doing operation")
  op(a,b)
}
```

```shell

2017-08-27 11:10:10 [INFO] Starting add operation caller=main
2017-08-27 11:10:10 [INFO] Do operation caller=main operation=add
```
