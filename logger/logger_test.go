package logger

import (
	"regexp"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefault(t *testing.T) {
	logger := Default()
	assert.NotNil(t, logger)

	log, ok := logger.(*logrus.Logger)
	assert.True(t, ok)
	assert.Equal(t, logrus.InfoLevel, log.Level)
}

func TestWithLogLevel(t *testing.T) {
	// Given
	opt := WithLogLevel(logrus.DebugLevel)

	// When
	logger := Default(opt)

	// Then
	log, ok := logger.(*logrus.Logger)
	assert.True(t, ok)
	assert.Equal(t, logrus.DebugLevel, log.Level)
}

func TestWithLogFormatter(t *testing.T) {
	// Given
	opt := WithLogFormatter(&logrus.JSONFormatter{})

	// When
	logger := Default(opt)

	// Then
	log, ok := logger.(*logrus.Logger)
	assert.True(t, ok)
	assert.IsType(t, &logrus.JSONFormatter{}, log.Formatter)
}

type TestHook struct {
	Fired bool
}

func (h *TestHook) Fire(_ *logrus.Entry) error {
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
	// Given
	hook := TestHook{}
	opt := WithHooks([]logrus.Hook{&hook})

	// When
	logger := Default(opt)
	logger.Info("test")

	// Then
	assert.True(t, hook.HasFired())
}

type TestLastEntryHook struct {
	lastEntry *logrus.Entry
}

func (h *TestLastEntryHook) Fire(entry *logrus.Entry) error {
	h.lastEntry = entry
	return nil
}

func (h *TestLastEntryHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.DebugLevel,
		logrus.InfoLevel,
	}
}

func TestWithRedactedFields(t *testing.T) {
	t.Run("no fields are redacted when no redaction fields are provided", func(t *testing.T) {
		// Given
		hook := TestLastEntryHook{}
		hookOpt := WithHooks([]logrus.Hook{&hook})
		redactedFieldsOpt := WithSetRedactedFields(nil)
		logger := Default(hookOpt, redactedFieldsOpt)

		// When
		logger.WithFields(logrus.Fields{
			"password": "secret",
			"other":    "value",
		}).Info("test")

		// Then
		log, ok := logger.(*logrus.Logger)
		assert.True(t, ok)
		assert.IsType(t, &RedactingFormatter{}, log.Formatter)
		require.Len(t, hook.lastEntry.Data, 2)
		require.Equal(t, "test", hook.lastEntry.Message)
		assert.Equal(t, "secret", hook.lastEntry.Data["password"])
		assert.Equal(t, "value", hook.lastEntry.Data["other"])
	})
	t.Run("nothing is redacted when the redactionOption is nil", func(t *testing.T) {
		// Given
		hook := TestLastEntryHook{}
		hookOpt := WithHooks([]logrus.Hook{&hook})
		redactedFieldsOpt := WithSetRedactedFields([]*RedactionOption{
			nil,
		})
		logger := Default(hookOpt, redactedFieldsOpt)

		// When
		logger.WithFields(logrus.Fields{
			"password": "secret",
			"other":    "value",
		}).Info("test")

		// Then
		log, ok := logger.(*logrus.Logger)
		assert.True(t, ok)
		assert.IsType(t, &RedactingFormatter{}, log.Formatter)
		require.Len(t, hook.lastEntry.Data, 2)
		require.Equal(t, "test", hook.lastEntry.Message)
		assert.Equal(t, "secret", hook.lastEntry.Data["password"])
		assert.Equal(t, "value", hook.lastEntry.Data["other"])
	})
	t.Run("a field is fully redacted when the regexp is nil", func(t *testing.T) {
		// Given
		hook := TestLastEntryHook{}
		hookOpt := WithHooks([]logrus.Hook{&hook})
		redactedFieldsOpt := WithSetRedactedFields([]*RedactionOption{
			{
				Field:  "password",
				Regexp: nil,
			},
		})
		logger := Default(hookOpt, redactedFieldsOpt)

		// When
		logger.WithFields(logrus.Fields{
			"password": "secret",
			"other":    "value",
		}).Info("test")

		// Then
		log, ok := logger.(*logrus.Logger)
		assert.True(t, ok)
		assert.IsType(t, &RedactingFormatter{}, log.Formatter)
		require.Len(t, hook.lastEntry.Data, 2)
		require.Equal(t, "test", hook.lastEntry.Message)
		assert.Equal(t, "[REDACTED]", hook.lastEntry.Data["password"])
		assert.Equal(t, "value", hook.lastEntry.Data["other"])
	})
	t.Run("a field is partially redacted when redactionOption with no replacement is provided", func(t *testing.T) {
		// Given
		hook := TestLastEntryHook{}
		hookOpt := WithHooks([]logrus.Hook{&hook})
		redactedFieldsOpt := WithSetRedactedFields([]*RedactionOption{
			{
				Field:  "path",
				Regexp: regexp.MustCompile(`token=[^&]+`),
			},
		})

		logger := Default(hookOpt, redactedFieldsOpt)

		// When
		logger.WithFields(logrus.Fields{
			"path":  "/apps/66b24069fb0de6002981dd79/logs?timestamp=1727183062&token=verySecretValue&stream=true",
			"other": "value",
		}).Info("test")

		// Then
		log, ok := logger.(*logrus.Logger)
		assert.True(t, ok)
		assert.IsType(t, &RedactingFormatter{}, log.Formatter)
		require.Len(t, hook.lastEntry.Data, 2)
		require.Equal(t, "test", hook.lastEntry.Message)
		assert.Equal(t, "/apps/66b24069fb0de6002981dd79/logs?timestamp=1727183062&[REDACTED]&stream=true", hook.lastEntry.Data["path"])
		assert.Equal(t, "value", hook.lastEntry.Data["other"])
	})

	t.Run("a field is partially redacted when redactionOption with replacement is provided", func(t *testing.T) {
		// Given
		hook := TestLastEntryHook{}
		hookOpt := WithHooks([]logrus.Hook{&hook})
		redactedFieldsOpt := WithSetRedactedFields([]*RedactionOption{
			{
				Field:       "path",
				Regexp:      regexp.MustCompile(`(token=)[^&]+`),
				ReplaceWith: "token=[HIDDEN]",
			},
		})

		logger := Default(hookOpt, redactedFieldsOpt)

		// When
		logger.WithFields(logrus.Fields{
			"path":  "/apps/66b24069fb0de6002981dd79/logs?timestamp=1727183062&token=verySecretValue&stream=true",
			"other": "value",
		}).Info("test")

		// Then
		log, ok := logger.(*logrus.Logger)
		assert.True(t, ok)
		assert.IsType(t, &RedactingFormatter{}, log.Formatter)
		require.Len(t, hook.lastEntry.Data, 2)
		require.Equal(t, "test", hook.lastEntry.Message)
		assert.Equal(t, "/apps/66b24069fb0de6002981dd79/logs?timestamp=1727183062&token=[HIDDEN]&stream=true", hook.lastEntry.Data["path"])
		assert.Equal(t, "value", hook.lastEntry.Data["other"])
	})
}

func TestNewContextWithLogger(t *testing.T) {
	ctx := NewContextWithLogger()
	logger := ctx.Value(loggerContextKey)
	assert.NotNil(t, logger)
}

func TestAddLoggerToContext(t *testing.T) {
	ctx := t.Context()
	ctx = AddLoggerToContext(ctx)
	logger := ctx.Value(loggerContextKey)
	assert.NotNil(t, logger)
}

func TestGet(t *testing.T) {
	ctx := t.Context()
	ctx = AddLoggerToContext(ctx)
	logger := Get(ctx)
	assert.NotNil(t, logger)
	_, ok := logger.(*logrus.Logger)
	assert.True(t, ok)
}

func TestWithFieldToCtx(t *testing.T) {
	ctx := t.Context()
	_, logger := WithFieldToCtx(ctx, "key", "value")
	assert.NotNil(t, logger)

	entry, ok := logger.(*logrus.Entry)
	assert.True(t, ok)
	assert.Equal(t, "value", entry.Data["key"])
}

func TestWithFieldsToCtx(t *testing.T) {
	ctx := t.Context()
	fields := logrus.Fields{"key1": "value1", "key2": "value2"}
	_, logger := WithFieldsToCtx(ctx, fields)
	assert.NotNil(t, logger)

	entry, ok := logger.(*logrus.Entry)
	assert.True(t, ok)
	assert.Equal(t, "value1", entry.Data["key1"])
	assert.Equal(t, "value2", entry.Data["key2"])
}

func TestToCtx(t *testing.T) {
	ctx := t.Context()
	logger := logrus.New()
	ctx = ToCtx(ctx, logger)
	assert.Equal(t, logger, ctx.Value(loggerContextKey))
}
