package heartbeats

import (
	// NOTE alternative: https://github.com/carlescere/scheduler
	"github.com/azer/logger"
	"github.com/rakyll/ticktock"
	"github.com/rakyll/ticktock/t"
)

const DEFAULT_INTERVAL string = "1h"

// FIXME as Clock attribute it stays immutable
var jobs []string

func init() {
	log.Info("registering heartbeat: clock")
	Add("clock", func(conf map[string]string) (Plugin, error) {
		return &Clock{
			Interval: conf["interval"],
			logger:   logger.New("sentinel.plugins.heartbeats.clock"),
		}, nil
	})
}

type Clock struct {
	Interval string
	logger   *logger.Logger
}

// NOTE use sensors.String() to get name ?
func (c Clock) Schedule(name string, job_ Job) error {
	log.Info("registering new clock-based job: %s", name)
	jobs = append(jobs, name)
	return ticktock.Schedule(
		name,
		job_,
		// TODO parametric
		// NOTE retries
		// FIXME &t.When{Each: c.Interval},
		&t.When{Each: c.Interval},
		//&t.When{Every: t.Every(10).Seconds()},
	)
}

func (c Clock) Start() {
	log.Info("starting scheduler (%d job(s))", len(jobs))
	ticktock.Start()
}

func (c Clock) Stop() {
	log.Info("stoping scheduler (%d job(s))", len(jobs))
	for _, jobName := range jobs {
		log.Info("\t- canceling job %s", jobName)
		ticktock.Cancel(jobName)
	}
}
