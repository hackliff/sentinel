package main

import (
	"github.com/azer/logger"
	"github.com/hashicorp/serf/command/agent"
	"github.com/hashicorp/serf/serf"

	"github.com/hackliff/sentinel/plugins/sensors"
	"github.com/hackliff/sentinel/plugins/triggers"
)

type Sentinel struct {
	Name string
	// NOTE currently, agent.Serf().Shutdown() would be enough
	Agent   *agent.Agent
	Trigger triggers.TriggerPlugin
	Sensor  sensors.SensorPlugin
}

type EventAlert struct {
	SelfFunc func() serf.Member
	logger   *logger.Logger
}

// TODO cast event type and perform checkers
func (d *EventAlert) HandleEvent(e serf.Event) {
	println("=================================")
	self := d.SelfFunc()
	println(self.Name)
	d.logger.Info("received an event: %v\n", e)
	println(e.String())
	println(e.EventType().String())
	println("=================================")
}
