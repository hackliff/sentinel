package heartbeats

import (
	"github.com/azer/logger"
)

var log = logger.New("sentinel.plugins.heartbeats")

type Job interface {
	//Run() error
	Run()
}

// Plugin is the interface to implement report triggers. Those run in
// their own gorountines, registering sensor plugins and  waiting for signals
// to execute them.
type Plugin interface {
	Schedule(string, Job)
	Stop()
}

// NOTE map[string]interface{} for better type support after casting ?
type Creator func(map[string]string) (Plugin, error)

var Plugins = map[string]Creator{}

func Add(name string, creator Creator) {
	Plugins[name] = creator
}
