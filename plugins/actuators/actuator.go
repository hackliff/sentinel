package actuators

import (
	"github.com/azer/logger"
	"github.com/hackliff/serf/serf"

	"github.com/hackliff/sentinel/plugins/adapters"
)

var log = logger.New("sentinel.plugins.actuators")

// TODO make it full compatible with telegraf plugins
// NOTE inspired by telegraf project plugin style
type Plugin interface {
	// NOTE String() method ?
	Description() string
	SampleConfig() string
	// NOTE a serf agent as argument for information and event triggers ?
	Gather(serf.Member, serf.Event) error
}

// NOTE map[string]interface{} for better type support after casting ?
type Creator func(adapters.Plugin, map[string]string) (Plugin, error)

var Plugins = map[string]Creator{}

func Add(name string, creator Creator) {
	Plugins[name] = creator
}
