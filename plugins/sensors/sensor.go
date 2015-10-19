package sensors

import "github.com/azer/logger"

var log = logger.New("sentinel.sensor")

func init() {
	log.Info("registering sensors")
}

// TODO make it full compatible with telegraf plugins
// NOTE inspired by telegraf project plugin style
type SensorPlugin interface {
	// NOTE String() method ?
	Description() string
	SampleConfig() string
	// NOTE a serf agent as argument for information and event triggers ?
	Monitor() error
}

// NOTE map[string]interface{} for better type support after casting ?
type SensorCreator func(map[string]string) (SensorPlugin, error)

var SensorPlugins = map[string]SensorCreator{}

func Add(name string, creator SensorCreator) {
	SensorPlugins[name] = creator
}
