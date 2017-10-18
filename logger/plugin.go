package logger

import (
	"sync"

	"github.com/sirupsen/logrus"
)

var Plugins = PluginManager{}

type Plugin interface {
	AddHook() (bool, logrus.Hook)
}

type PluginManager struct {
	plugins []Plugin
	lock    sync.Mutex
}

func (m PluginManager) AddPlugin(plugin Plugin) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.plugins = append(m.plugins, plugin)
}

func (m PluginManager) Hooks() []logrus.Hook {
	m.lock.Lock()
	defer m.lock.Unlock()

	hooks := []logrus.Hook{}

	for _, plugin := range m.plugins {
		add, hook := plugin.AddHook()
		if add {
			hooks = append(hooks, hook)
		}
	}

	return hooks
}
