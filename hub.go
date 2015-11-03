package main

import (
	"github.com/hashicorp/serf/command/agent"
	"github.com/hashicorp/serf/serf"

	"github.com/hackliff/sentinel/plugins/actuators"
)

type EventHub struct {
	Actuator actuators.Plugin
	Self     serf.Member
	Filters  []agent.EventFilter
}

var HEARTBEAT_EVENT serf.Event = serf.UserEvent{
	Name:     "heartbeat",
	Payload:  []byte(""),
	Coalesce: true,
}

// NOTE Run is currently a way to implement ticktock job runner
func (h EventHub) Run() error {
	h.HandleEvent(HEARTBEAT_EVENT)
	return nil
}

func (h EventHub) HandleEvent(e serf.Event) {
	for _, filter := range h.Filters {
		if filter.Invoke(e) {
			log.Info("sentinel processing event: %v", e)
			// NOTE parse event e ?
			if err := h.Actuator.Gather(h.Self, e); err != nil {
				log.Error("error invoking sentinel: %s", err)
			}
			// we're done as soon as one filter matched
			return
		} else {
			log.Info("skipping event: %v", e)
		}
	}
}
