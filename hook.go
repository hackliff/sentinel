package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/hackliff/serf/command/agent"

	"github.com/hackliff/sentinel/config"
	"github.com/hackliff/sentinel/plugins/actuators"
	"github.com/hackliff/sentinel/plugins/adapters"
	"github.com/hackliff/sentinel/plugins/heartbeats"
)

func composeEventHub(cfg *config.Config, serfAgent *agent.Agent) (*EventHub, error) {
	log.Info("composing new sentinel %#v", cfg)

	// setting up common adapter
	adapterCreator, _ := adapters.Plugins[cfg.Adapter["plugin"]]
	adapter_, _ := adapterCreator(cfg.Adapter)

	// setting up common actuator
	// TODO handle errors
	actuatorCreator, _ := actuators.Plugins[cfg.Actuator["plugin"]]
	actuator_, _ := actuatorCreator(adapter_, cfg.Actuator)

	return &EventHub{
		Actuator: actuator_,
		Self:     serfAgent.Serf().LocalMember(),
		Filters:  cfg.Heartbeat.Filters,
	}, nil
}

func onAgentReady(serfAgent *agent.Agent, serfConfig *agent.Config, shutdownCh <-chan struct{}) {
	cfgPath, err := config.Path()
	if err != nil {
		log.Error("error searching for configuration file: %v", err)
		return
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Error("loading configuration: %v", err)
		return
	}

	hub, _ := composeEventHub(cfg, serfAgent)

	if cfg.Heartbeat.Plugin == "event" {
		log.Info("registering new event handler")
		serfAgent.RegisterEventHandler(hub)
	} else {
		// setting up heartbeat
		heartbeatCreator, _ := heartbeats.Plugins[cfg.Heartbeat.Plugin]
		heartbeat_, _ := heartbeatCreator(cfg.Heartbeat.Properties)
		defer heartbeat_.Stop()

		log.Info("scheduling new sentinel actuator")
		if err := heartbeat_.Schedule(cfg.Name, hub); err != nil {
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
	// FIXME serf shutdowns too damn fast !
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
