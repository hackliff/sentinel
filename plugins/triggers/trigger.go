package triggers

import (
	"github.com/azer/logger"
	"github.com/hackliff/sentinel/plugins/sensors"
)

var log = logger.New("sentinel.trigger")

// TriggerPlugin is the interface to implement report triggers. Those run in
// their own gorountines, registering sensor plugins and  waiting for signals
// to execute them.
type TriggerPlugin interface {
	Start()
	Schedule(string, sensors.SensorPlugin) error
	Stop()
}

// NOTE map[string]interface{} for better type support after casting ?
type TriggerCreator func(map[string]string) (TriggerPlugin, error)

var TriggerPlugins = map[string]TriggerCreator{}

func Add(name string, creator TriggerCreator) {
	TriggerPlugins[name] = creator
}
