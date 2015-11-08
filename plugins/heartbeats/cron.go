package heartbeats

import (
	"github.com/azer/logger"
	"github.com/bamzi/jobrunner"
)

const DEFAULT_INTERVAL string = "@every 1h"

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

func (c Clock) Schedule(name string, job_ Job) {
	log.Info("registering new cron job: %s", name)
	jobrunner.Start()
	jobrunner.Schedule(c.Interval, job_)
}

func (c Clock) Stop() {
	log.Info("stoping cron scheduler: %#v", jobrunner.StatusJson())
	jobrunner.Stop()
}
