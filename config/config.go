package config

import (
	"github.com/azer/logger"
	"github.com/hackliff/serf/command/agent"
	"github.com/olebedev/config"
)

var log = logger.New("sentinel.config")

const DEFAULT_NAME string = "sentinel"

type PluginConfig struct {
	Plugin string
	Opts   map[string]string
}

type HeartbeatConfig struct {
	// NOTE create Heartbeat plugin here ?
	Plugin  string
	Opts    map[string]string
	Filters []agent.EventFilter
}

type Config struct {
	Name      string
	Actuator  *PluginConfig
	Heartbeat *HeartbeatConfig
	Adapter   *PluginConfig
}

func Load(confPath string) (*Config, error) {
	log.Info("loading configuration file: %s", confPath)
	cfg, err := config.ParseYamlFile(confPath)
	if err != nil {
		return nil, err
	}

	// NOTE use defaults ?
	cfgName, err := cfg.String("sentinel.name")
	if err != nil {
		cfgName = DEFAULT_NAME
	}
	log.Info("parsing plugins configuration")
	adapterCfg, err := cfg.String("sentinel.adapter")
	if err != nil {
		return nil, err
	}
	actuatorCfg, err := cfg.String("sentinel.actuator")
	if err != nil {
		return nil, err
	}
	heartbeatCfg, err := cfg.String("sentinel.heartbeat")
	if err != nil {
		return nil, err
	}

	parser := NewREParser()

	partials := parser.Parse(heartbeatCfg)
	heartbeatConf_ := &HeartbeatConfig{
		Plugin: partials.Plugin,
		// NOTE pop out "on" ?
		Opts: partials.Opts,
		// FIXME what if "on" is not specified ?
		Filters: agent.ParseEventFilter(partials.Opts["on"]),
	}

	log.Info("loading sentinel configuration")
	return &Config{
		Name:      cfgName,
		Adapter:   parser.Parse(adapterCfg),
		Actuator:  parser.Parse(actuatorCfg),
		Heartbeat: heartbeatConf_,
	}, nil
}
