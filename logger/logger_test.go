package logger

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type testHook struct{ calls int }

func (h *testHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.ErrorLevel}
}

func (h *testHook) Fire(entry *logrus.Entry) error {
	h.calls++
	return nil
}

func TestDefault(t *testing.T) {
	examples := []struct {
		Name   string
		Hooks  func() []logrus.Hook
		Logger func(hooks []logrus.Hook) *logrus.Logger
		Expect func(t *testing.T, logger *logrus.Logger, hooks []logrus.Hook)
	}{
		{
			Name:  "should have no hook by default",
			Hooks: func() []logrus.Hook { return []logrus.Hook{} },
			Logger: func(hooks []logrus.Hook) *logrus.Logger {
				return Default(hooks...)
			},
			Expect: func(t *testing.T, logger *logrus.Logger, hooks []logrus.Hook) {
				assert.Len(t, logger.Hooks, 0)
			},
		}, {
			Name:  "should a one hook if given in argument",
			Hooks: func() []logrus.Hook { return []logrus.Hook{&testHook{}} },
			Logger: func(hooks []logrus.Hook) *logrus.Logger {
				return Default(hooks...)
			},
			Expect: func(t *testing.T, logger *logrus.Logger, hooks []logrus.Hook) {
				assert.Len(t, logger.Hooks, 1)
				logger.Error("test")
				assert.Equal(t, 1, hooks[0].(*testHook).calls)
			},
		},
	}

	for _, example := range examples {
		t.Run(example.Name, func(t *testing.T) {
			hooks := example.Hooks()
			logger := example.Logger(hooks)
			example.Expect(t, logger, hooks)
		})
	}
}
