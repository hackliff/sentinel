package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hackliff/serf/command/agent"

	"github.com/hackliff/sentinel/config"
	"github.com/hackliff/sentinel/plugins/actuators"
	"github.com/hackliff/sentinel/plugins/adapters"
	"github.com/hackliff/sentinel/plugins/heartbeats"
)

//func composeEventHub(cfg *config.Config, serfAgent *agent.Agent) (*EventHub, error) {
func composeEventHub(cfg *config.Config, serfAgent *agent.Agent) (agent.EventHandler, error) {
	log.Info("composing new sentinel %#v", cfg)

	log.Info("setting up sentinel adapter: %s", cfg.Adapter.Plugin)
	adapterCreator, ok := adapters.Plugins[cfg.Adapter.Plugin]
	if !ok {
		return nil, fmt.Errorf("invalid adapter: %s", cfg.Adapter.Plugin)
	}
	adapter_, err := adapterCreator(cfg.Adapter.Opts)
	if err != nil {
		return nil, err
	}

	// setting up common actuator
	log.Info("setting up sentinel actuator: %s", cfg.Actuator.Plugin)
	actuatorCreator, ok := actuators.Plugins[cfg.Actuator.Plugin]
	if !ok {
		return nil, fmt.Errorf("invalid actuator: %s", cfg.Actuator.Plugin)
	}
	actuator_, err := actuatorCreator(adapter_, cfg.Actuator.Opts)
	if err != nil {
		return nil, err
	}

	return EventHub{
		Actuator: actuator_,
		Self:     serfAgent.Serf().LocalMember(),
		Filters:  cfg.Heartbeat.Filters,
	}, nil
}

// TODO remove serfConif
func onAgentReady(serfAgent *agent.Agent, serfConfig *agent.Config, shutdownCh <-chan struct{}) {
	// search for the path in environment, fallback to default
	cfgPath, err := config.Path()
	if err != nil {
		log.Error("error searching configuration file: %v", err)
		return
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Error("failed to load configuration: %v", err)
		return
	}

	hub, err := composeEventHub(cfg, serfAgent)
	if err != nil {
		log.Error("failed to compose sentinel event hub: %s", err)
		return
	}

	if cfg.Heartbeat.Plugin == "event" {
		log.Info("registering new event handler")
		serfAgent.RegisterEventHandler(hub)
	} else {
		// setting up heartbeat
		heartbeatCreator, ok := heartbeats.Plugins[cfg.Heartbeat.Plugin]
		if !ok {
			log.Error("invalid heartbeat: %s", cfg.Heartbeat.Plugin)
			return
		}
		heartbeat_, err := heartbeatCreator(cfg.Heartbeat.Opts)
		if err != nil {
			log.Error("failed to initialize heartbeat: %s", err)
			return
		}
		defer heartbeat_.Stop()

		log.Info("scheduling new sentinel actuator")
		heartbeat_.Schedule(cfg.Name, hub)
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
