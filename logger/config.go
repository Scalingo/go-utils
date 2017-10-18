package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var defaultLogLevel logrus.Level = logrus.DebugLevel
var defaultFormatter logrus.Formatter = &logrus.TextFormatter{
	TimestampFormat: "2006-01-02T15:04:05.000",
	FullTimestamp:   true,
}

func formatter() logrus.Formatter {
	switch os.Getenv("LOGGER_TYPE") {
	case "json":
		return new(logrus.JSONFormatter)
	case "text":
		return new(logrus.TextFormatter)
	default:
		return defaultFormatter
	}
}

func logLevel() logrus.Level {
	switch os.Getenv("LOGGER_LEVEL") {
	case "panic":
		return logrus.PanicLevel
	case "fatal":
		return logrus.FatalLevel
	case "warn":
		return logrus.WarnLevel
	case "info":
		return logrus.InfoLevel
	case "debug":
		return logrus.DebugLevel
	default:
		return defaultLogLevel
	}
}
