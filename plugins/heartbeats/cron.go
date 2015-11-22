package heartbeats

import (
	"github.com/azer/logger"
	"github.com/bamzi/jobrunner"
	"github.com/hackliff/serf/command/agent"
	"github.com/hackliff/serf/serf"
)

const DEFAULT_INTERVAL string = "@every 1h"

// NOTE build this structure for each call, with useful info ?
var CRON_EVENT serf.Event = serf.UserEvent{
	Name:     "heartbeat-cron",
	Payload:  []byte(""),
	Coalesce: true,
}

type Job struct {
	EventHandler agent.EventHandler
}

func (j Job) Run() {
	j.EventHandler.HandleEvent(CRON_EVENT)
}

func init() {
	log.Info("registering heartbeat: clock")
	Add("cron", func(conf map[string]string) (Plugin, error) {
		interval, ok := conf["interval"]
		if !ok {
			interval = DEFAULT_INTERVAL
		}
		return &Clock{
			Interval: interval,
			logger:   logger.New("sentinel.plugins.heartbeats.clock"),
		}, nil
	})
}

type Clock struct {
	Interval string
	logger   *logger.Logger
}

func (c Clock) Schedule(name string, eventHandler agent.EventHandler) {
	log.Info("registering new cron job: %s", name)
	job := Job{eventHandler}

	jobrunner.Start()
	jobrunner.Schedule(c.Interval, job)
}

func (c Clock) Stop() {
	log.Info("stoping cron scheduler: %#v", jobrunner.StatusJson())
	jobrunner.Stop()
}
