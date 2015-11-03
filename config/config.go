package config

import (
	"github.com/azer/logger"
	"github.com/hashicorp/serf/command/agent"
)

var log = logger.New("sentinel.config")

type HeartbeatConf struct {
	// NOTE create Heartbeat plugin here ?
	Plugin     string
	Properties map[string]string
	Filters    []agent.EventFilter
}

type Config struct {
	Name      string
	Actuator  map[string]string
	Heartbeat *HeartbeatConf
	Adapter   map[string]string
}

func Load(confPath string) (*Config, error) {
	log.Info("loading configuration file: %s", confPath)
	cfg, err := readConfigFile(confPath)
	if err != nil {
		return nil, err
	}
	log.Info("parsed configuration: %#v", cfg)

	// TODO handle error
	log.Info("parsing heartbeat configuration")
	adapterCfg, _ := cfg.String("sentinel.adapter")
	actuatorCfg, _ := cfg.String("sentinel.actuator")
	heartbeatCfg, _ := cfg.String("sentinel.heartbeat")
	partials := parsePluginConfig(heartbeatCfg)
	heartbeatConf_ := &HeartbeatConf{
		Plugin: partials["plugin"],
		// NOTE pop out "plugin" and "on" ?
		Properties: partials,
		Filters:    agent.ParseEventFilter(partials["on"]),
	}

	log.Info("loading sentinel configuration")
	return &Config{
		Name:      "hawkeye",
		Adapter:   parsePluginConfig(adapterCfg),
		Actuator:  parsePluginConfig(actuatorCfg),
		Heartbeat: heartbeatConf_,
	}, nil
}
