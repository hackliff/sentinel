package main

import (
	"github.com/hackliff/serf/command/agent"
	"github.com/hackliff/serf/serf"

	"github.com/hackliff/sentinel/plugins/actuators"
)

type EventHub struct {
	Actuator actuators.Plugin
	Self     serf.Member
	Filters  []agent.EventFilter
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
