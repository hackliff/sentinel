package main

import (
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/azer/logger"
	"github.com/hashicorp/serf/command"
	"github.com/hashicorp/serf/command/agent"
	"github.com/hashicorp/serf/serf"

	"github.com/hackliff/sentinel/plugins/sensors"
	"github.com/hackliff/sentinel/plugins/triggers"
)

func startAgent(config *agent.Config, serfAgent *agent.Agent) *agent.AgentIPC {
	// Create a log writer, and wrap a logOutput around it
	logWriter := agent.NewLogWriter(512)
	logOutput := io.MultiWriter(agent.LevelFilter(), logWriter)

	eh := &EventAlert{
		SelfFunc: func() serf.Member { return serfAgent.Serf().LocalMember() },
		// TODO use top-level logger
		logger: logger.New("sentinel.handler"),
		// TODO Give it a Radio or a Sensor
	}
	serfAgent.RegisterEventHandler(eh)

	// Start the agent after the handler is registered
	if err := serfAgent.Start(); err != nil {
		log.Error("Failed to start the Serf agent: %v", err)
		return nil
	}

	// Parse the bind address information
	bindIP, bindPort, err := config.AddrParts(config.BindAddr)
	bindAddr := &net.TCPAddr{IP: net.ParseIP(bindIP), Port: bindPort}

	// TODO Start the discovery layer

	// Setup the RPC listener
	rpcListener, err := net.Listen("tcp", config.RPCAddr)
	if err != nil {
		log.Error("Error starting RPC listener: %s", err)
		return nil
	}

	// TODO its own function
	log.Info("Starting Serf agent RPC...")
	ipc := agent.NewAgentIPC(serfAgent, config.RPCAuthKey, rpcListener, logOutput, logWriter)

	log.Info("Serf agent running!")
	log.Info("     Node name: '%s'", config.NodeName)
	log.Info("     Bind addr: '%s'", bindAddr.String())

	if config.AdvertiseAddr != "" {
		advertiseIP, advertisePort, _ := config.AddrParts(config.AdvertiseAddr)
		advertiseAddr := (&net.TCPAddr{IP: net.ParseIP(advertiseIP), Port: advertisePort}).String()
		log.Info("Advertise addr: '%s'", advertiseAddr)
	}

	log.Info("      RPC addr: '%s'", config.RPCAddr)
	log.Info("     Encrypted: %#v", serfAgent.Serf().EncryptionEnabled())
	log.Info("      Snapshot: %v", config.SnapshotPath != "")
	log.Info("       Profile: %s", config.Profile)

	return ipc
}

// TODO use startupJoin
func RunAgent(conf *Config) int {
	serfAgent, _ := agent.Create(conf.Agent, conf.Serf, os.Stderr)

	ipc := startAgent(conf.Agent, serfAgent)
	if ipc == nil {
		return 1
	}

	RPC, err := command.RPCClient(conf.Agent.RPCAddr, conf.Agent.RPCAuthKey)
	if err != nil {
		log.Error("Error connecting to Serf agent: %s", err)
		return 1
	}

	defer RPC.Close()
	defer ipc.Shutdown()
	defer serfAgent.Shutdown()

	// NOTE sentinel object is just a container at this point
	for _, sentinelConf := range conf.Sentinels {
		//pingSensor := sensor.PingSensor{sentinelConf.Sensor["endpoint"]}
		// TODO handle errors
		sensorCreator, _ := sensors.SensorPlugins[sentinelConf.Sensor["plugin"]]
		pingSensor, _ := sensorCreator(RPC, sentinelConf.Sensor)

		triggerCreator, _ := triggers.TriggerPlugins[sentinelConf.Trigger["plugin"]]
		clockTrigger_, _ := triggerCreator(sentinelConf.Trigger)

		sentinel_ := Sentinel{
			Name:    sentinelConf.Name,
			Agent:   serfAgent,
			Trigger: clockTrigger_,
			Sensor:  pingSensor,
		}
		// TODO handle error
		if err := sentinel_.Trigger.Schedule(sentinel_.Name, sentinel_.Sensor); err != nil {
			log.Error("error scheduling sentinel")
			return 1
		}
		go sentinel_.Trigger.Start()
		defer sentinel_.Trigger.Stop()
	}

	log.Info("")
	log.Info("Log data will now stream in as it occurs:\n")

	// Wait for exit
	return handleSignals(serfAgent)
}

// NOTE sentinel_.HandleSignals()
// TODO graceful leave and SIGHUP support
//func handleSignals(agent *agent.Agent) int {
func handleSignals(ag *agent.Agent) int {
	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	// Wait for signal
	select {
	case sig := <-signalCh:
		log.Info("Caught signal: %v\n", sig)
		return 1
	case <-ag.ShutdownCh():
		// Agent is already shutdown!
		return 0
	}
}
