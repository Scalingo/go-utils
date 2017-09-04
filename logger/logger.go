package logger

import (
	"context"

	"github.com/Sirupsen/logrus"
)

var logLevel logrus.Level = logrus.DebugLevel
var formatter logrus.Formatter = &logrus.TextFormatter{
	TimestampFormat: "2006-01-02T15:04:05.000",
	FullTimestamp:   true,
}
var hooks []logrus.Hook

// SetConfig set the configuration at the level package.
// level: The minimum log level to log.
// formatter: The formatter used to format logs (see logrus formatter)
//
// By default the configuration use the debug level and a text formatter
func SetConfig(level logrus.Level, f logrus.Formatter) {
	logLevel = level
	formatter = f
}

// AddHook add a hook to the default logger stack
func AddHook(hook logrus.Hook) {
	hooks = append(hooks, hook)
}

// Default generate a logrus logger with the configuration set by the SetConfig and AddHook methods
func Default() *logrus.Logger {
	logger := logrus.StandardLogger()
	logger.SetLevel(logLevel)
	logger.Formatter = formatter

	for _, h := range hooks {
		logger.AddHook(h)
	}

	return logger
}

// NewContextWithLogger generate a new context (based on context.Background()) and add a Default() logger on top of it
func NewContextWithLogger() context.Context {
	return AddLoggerToContext(context.Background())
}

// AddLoggerToContext add the Default() logger on top of the current context
func AddLoggerToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, "logger", logrus.NewEntry(Default()))
}

// Get return the logger stored in the context or create a new one if the logger is not set
func Get(ctx context.Context) logrus.FieldLogger {
	if logger, ok := ctx.Value("logger").(logrus.FieldLogger); ok {
		return logger
	}

	return Default().WithField("invalid_context", true)
}

// ToCtx add a logger to a context
func ToCtx(ctx context.Context, logger logrus.FieldLogger) context.Context {
	return context.WithValue(ctx, "logger", logger)
}
