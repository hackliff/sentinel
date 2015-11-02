package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/serf/command/agent"
	"github.com/hashicorp/serf/serf"
	"github.com/olebedev/config"
)

type HeartbeatConf struct {
	// NOTE create Heartbeat plugin here ?
	Plugin     string
	Properties map[string]string
	Filters    []agent.EventFilter
}

type SentinelConfig struct {
	Name      string
	Actuator  map[string]string
	Heartbeat *HeartbeatConf
	Adapter   map[string]string
}

type Config struct {
	Agent    *agent.Config
	Serf     *serf.Config
	Sentinel *SentinelConfig
}

// grammar: <plugin type>: name key1=value key2=an,array
// NOTE use plugin-defined spec to cast values and validate ?
func parsePluginConfig(descr string) map[string]string {
	props := make(map[string]string)
	parts := strings.Split(descr, " ")
	props["plugin"] = parts[0]
	for _, part := range parts[1:] {
		kv := strings.Split(part, "=")
		props[kv[0]] = kv[1]
	}

	return props
}

// TODO render Go template (ex: {{ hostname }}) or support consul-template ?
func readConfigFile(confPath string) (*config.Config, error) {
	return config.ParseYamlFile(confPath)
}

func completeAgentConfig() *agent.Config {
	config := agent.DefaultConfig()
	// TODO config = agent.MergeConfig(config, &cmdConfig)

	if config.NodeName == "" {
		hostname, err := os.Hostname()
		if err != nil {
			log.Error("error determining hostname: %s", err)
			return nil
		}
		config.NodeName = hostname
	}

	return config
}

func LoadConfiguration(confPath string) (*Config, error) {
	// TODO load config file
	log.Info("loading configuration file")
	cfg, err := readConfigFile(confPath)
	if err != nil {
		return nil, err
	}
	log.Info("parsed configuration: %#v", cfg)

	log.Info("loading serf agent configuration")
	agentConf := completeAgentConfig()
	if agentConf == nil {
		return nil, fmt.Errorf("failed to load agent configuration")
	}

	log.Info("parsing heartbeat configuration")
	adapterCfg, _ := cfg.String("sentinel.adapter")
	actuatorCfg, _ := cfg.String("sentinel.actuator")
	heartbeatCfg, _ := cfg.String("sentinel.heartbeat")
	//partials := parsePluginConfig("clock interval=10s")
	//partials := parsePluginConfig("event on=member-join,user:actuator-failed,query")
	// TODO handle error
	partials := parsePluginConfig(heartbeatCfg)
	heartbeatConf_ := &HeartbeatConf{
		Plugin: partials["plugin"],
		// NOTE pop out "plugin" and "on" ?
		Properties: partials,
		Filters:    agent.ParseEventFilter(partials["on"]),
	}

	log.Info("loading sentinel configuration")
	hawkeyeConf := &SentinelConfig{
		Name: "hawkeye",
		//Actuator:  parsePluginConfig("ping endpoint=example.com"),
		//Adapter: parsePluginConfig("shell level=debug"),
		Adapter:   parsePluginConfig(adapterCfg),
		Actuator:  parsePluginConfig(actuatorCfg),
		Heartbeat: heartbeatConf_,
	}

	return &Config{
		Agent:    agentConf,
		Serf:     serf.DefaultConfig(),
		Sentinel: hawkeyeConf,
	}, nil
}
