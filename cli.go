package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/hackliff/serf/command/agent"

	"github.com/hackliff/sentinel/plugins/actuators"
	"github.com/hackliff/sentinel/plugins/adapters"
	"github.com/hackliff/sentinel/plugins/heartbeats"
)

func onAgentReady(serfAgent *agent.Agent, serfConfig *agent.Config, shutdownCh <-chan struct{}) {
	// TODO parametric path from config
	confPath := "./conf.yml"
	conf, err := LoadConfiguration(confPath)
	if err != nil {
		log.Error("loading configuration: %v", err)
		return
	}

	log.Info("preparing new sentinel %#v", conf.Sentinel)

	// setting up common adapter
	adapterCreator, _ := adapters.Plugins[conf.Sentinel.Adapter["plugin"]]
	adapter_, _ := adapterCreator(conf.Sentinel.Adapter)

	// setting up common actuator
	// TODO handle errors
	actuatorCreator, _ := actuators.Plugins[conf.Sentinel.Actuator["plugin"]]
	actuator_, _ := actuatorCreator(adapter_, conf.Sentinel.Actuator)

	hub := &EventHub{
		Actuator: actuator_,
		Self:     serfAgent.Serf().LocalMember(),
		Filters:  conf.Sentinel.Heartbeat.Filters,
	}

	if conf.Sentinel.Heartbeat.Plugin == "event" {
		log.Info("registering new event handler")
		serfAgent.RegisterEventHandler(hub)
	} else {
		// setting up heartbeat
		heartbeatCreator, _ := heartbeats.Plugins[conf.Sentinel.Heartbeat.Plugin]
		heartbeat_, _ := heartbeatCreator(conf.Sentinel.Heartbeat.Properties)
		defer heartbeat_.Stop()

		log.Info("scheduling new sentinel actuator")
		if err := heartbeat_.Schedule(conf.Sentinel.Name, hub); err != nil {
			log.Info("error scheduling sentinel")
			return
		}

		log.Info("activating heartbeat")
		go heartbeat_.Start()
	}

	handleSignals(serfAgent, shutdownCh)
	return
}

func handleSignals(serfAgent *agent.Agent, shutdownCh <-chan struct{}) {
	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	// Wait for signal
	// FIXME serf shutdown too damn fast !
WAIT:
	select {
	case sig := <-signalCh:
		log.Info("caught signal: %v\n", sig)
		if sig == syscall.SIGHUP {
			// serf catches this signal to reload its configuration
			// NOTE reload sentinel conf ?
			goto WAIT
		} else {
			return
		}
	case <-shutdownCh:
		return
	case <-serfAgent.ShutdownCh():
		// Agent is already shutdown!
		return
	}
}
