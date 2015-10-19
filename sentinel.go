package main

import (
	"fmt"

	"github.com/azer/logger"
	"github.com/hashicorp/serf/command/agent"
	"github.com/hashicorp/serf/serf"

	"github.com/hackliff/sentinel-factory/plugins/sensors"
	"github.com/hackliff/sentinel-factory/plugins/triggers"
)

type Sentinel struct {
	// NOTE currently, agent.Serf().Shutdown() would be enough
	Name    string
	Agent   *agent.Agent
	Trigger triggers.TriggerPlugin
	Sensor  sensors.SensorPlugin
}

func (s *Sentinel) Monitor() {
	// TODO handle error
	s.Trigger.Schedule(s.Name, s.Sensor)
	go s.Trigger.Start()
	serfShutdownCh := s.Agent.Serf().ShutdownCh()

	for {
		select {
		case <-serfShutdownCh:
			fmt.Printf("[WARN] sentinel: Serf shutdown detected, quitting\n")
			// NOTE redundant with defer agent.Shutdown() in cli.go
			s.Trigger.Stop()
			//s.agent.Shutdown()
			return

		case <-s.Agent.ShutdownCh():
			fmt.Printf("[WARN] sentinel: Agent shutdown detected, quitting\n")
			return
		}
	}
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
