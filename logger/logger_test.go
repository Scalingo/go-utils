package logger

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestDefault(t *testing.T) {
	logger := Default()
	assert.NotNil(t, logger)
	assert.Equal(t, logrus.InfoLevel, logger.(*logrus.Logger).Level)
}

func TestWithLogLevel(t *testing.T) {
	logger := Default(WithLogLevel(logrus.DebugLevel))
	assert.Equal(t, logrus.DebugLevel, logger.(*logrus.Logger).Level)
}

func TestWithLogFormatter(t *testing.T) {
	logger := Default(WithLogFormatter(&logrus.JSONFormatter{}))
	assert.IsType(t, &logrus.JSONFormatter{}, logger.(*logrus.Logger).Formatter)
}

type TestHook struct {
	Fired bool
}

func (h *TestHook) Fire(entry *logrus.Entry) error {
	h.Fired = true
	return nil
}

func (h *TestHook) HasFired() bool {
	return h.Fired
}

func (h *TestHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.DebugLevel,
		logrus.InfoLevel,
	}
}

func TestWithHooks(t *testing.T) {
	hook := TestHook{}
	logger := Default(WithHooks([]logrus.Hook{&hook}))
	logger.Info("test")
	assert.True(t, hook.HasFired())
}

type TestRedactedHook struct {
	lastEntry *logrus.Entry
}

func (h *TestRedactedHook) Fire(entry *logrus.Entry) error {
	h.lastEntry = entry
	return nil
}

func (h *TestRedactedHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.DebugLevel,
		logrus.InfoLevel,
	}
}

func TestWithRedactedFields(t *testing.T) {
	hook := TestRedactedHook{}
	logger := Default(WithHooks([]logrus.Hook{&hook}), WithSetRedactedFields([]string{"password"}))
	assert.IsType(t, &RedactingFormatter{}, logger.(*logrus.Logger).Formatter)
	logger.Info("test")
	// capture output from logger
	logger.WithFields(logrus.Fields{
		"password": "secret",
		"other":    "value",
	}).Info("test")
	require.Len(t, hook.lastEntry.Data, 2)
	require.Equal(t, "test", hook.lastEntry.Message)
	assert.Equal(t, "REDACTED", hook.lastEntry.Data["password"])
	assert.Equal(t, "value", hook.lastEntry.Data["other"])
}

func TestNewContextWithLogger(t *testing.T) {
	ctx := NewContextWithLogger()
	logger := ctx.Value("logger")
	assert.NotNil(t, logger)
}

func TestAddLoggerToContext(t *testing.T) {
	ctx := context.Background()
	ctx = AddLoggerToContext(ctx)
	logger := ctx.Value("logger")
	assert.NotNil(t, logger)
}

func TestGet(t *testing.T) {
	ctx := context.Background()
	ctx = AddLoggerToContext(ctx)
	logger := Get(ctx)
	assert.NotNil(t, logger)
	_, ok := logger.(*logrus.Logger)
	assert.True(t, ok)
}

func TestWithFieldToCtx(t *testing.T) {
	ctx := context.Background()
	ctx, logger := WithFieldToCtx(ctx, "key", "value")
	assert.NotNil(t, logger)
	assert.Equal(t, "value", logger.(*logrus.Entry).Data["key"])
}

func TestWithFieldsToCtx(t *testing.T) {
	ctx := context.Background()
	fields := logrus.Fields{"key1": "value1", "key2": "value2"}
	ctx, logger := WithFieldsToCtx(ctx, fields)
	assert.NotNil(t, logger)
	assert.Equal(t, "value1", logger.(*logrus.Entry).Data["key1"])
	assert.Equal(t, "value2", logger.(*logrus.Entry).Data["key2"])
}

func TestToCtx(t *testing.T) {
	ctx := context.Background()
	logger := logrus.New()
	ctx = ToCtx(ctx, logger)
	assert.Equal(t, logger, ctx.Value("logger"))
}
