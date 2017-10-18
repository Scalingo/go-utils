package rollbarplugin

import (
	"os"

	"github.com/Scalingo/go-internal-tools/logger"
	logrus_rollbar "github.com/Scalingo/logrus-rollbar"
	"github.com/sirupsen/logrus"
	"github.com/stvp/rollbar"
)

type RollbarPlugin struct{}

func EnsurePresent() {
	logger.Plugins.EnsurePlugin(RollbarPlugin{})
}

func (p RollbarPlugin) Name() string {
	return "rollbar"
}

func (p RollbarPlugin) AddHook() (bool, logrus.Hook) {
	token := os.Getenv("ROLLBAR_API_KEY")
	if token == "" {
		token = os.Getenv("ROLLBAR_TOKEN")
	}

	if token == "" {
		return false, nil
	}

	rollbar.Token = token
	rollbar.Environment = os.Getenv("GO_ENV")

	return true, logrus_rollbar.New(8)
}
