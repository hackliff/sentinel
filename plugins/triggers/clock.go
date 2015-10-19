package triggers

import (
	// NOTE alternative: https://github.com/carlescere/scheduler
	"github.com/rakyll/ticktock"
	"github.com/rakyll/ticktock/t"

	"github.com/hackliff/sentinel-factory/plugins/sensors"
)

const DEFAULT_INTERVAL string = "1h"

// FIXME as ClockTrigger attribute it stays immutable
var jobs []string

func init() {
	Add("clock", func(conf map[string]string) (TriggerPlugin, error) {
		return &ClockTrigger{}, nil
	})
}

type SentinelJob struct {
	// NOTE could we define here a generic runner and abstract away the sensor ?
	sensor sensors.SensorPlugin
}

// NOTE pass a message with interval or uptime ?
func (s SentinelJob) Run() error {
	return s.sensor.Monitor()
}

type ClockTrigger struct{}

// NOTE use sensors.String() to get name ?
func (c ClockTrigger) Schedule(name string, sensor_ sensors.SensorPlugin) error {
	log.Info("registering new clock-based job: %s", name)
	jobs = append(jobs, name)
	return ticktock.Schedule(
		name,
		SentinelJob{sensor_},
		// TODO parametric
		// NOTE retries
		// FIXME &t.When{Each: "10s"},
		&t.When{Every: t.Every(10).Seconds()},
	)
}

func (c ClockTrigger) Start() {
	log.Info("starting scheduler (%d job(s))", len(jobs))
	ticktock.Start()
}

func (c ClockTrigger) Stop() {
	log.Info("stoping scheduler (%d job(s))", len(jobs))
	for _, jobName := range jobs {
		log.Info("\t- canceling job %s", jobName)
		ticktock.Cancel(jobName)
	}
}
