package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/serf/command/agent"
	"github.com/hashicorp/serf/serf"
	"github.com/olebedev/config"
)

type SentinelConfig struct {
	Name    string
	Sensor  map[string]string
	Trigger map[string]string
}

type Config struct {
	Agent     *agent.Config
	Serf      *serf.Config
	Sentinels []SentinelConfig
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
			log.Error("Error determining hostname: %s", err)
			return nil
		}
		config.NodeName = hostname
	}

	return config
}

// TODO merge(opts, config)
func LoadConfiguration(opts *Options) (*Config, error) {
	log.Info("loading configuration file")
	/*
	 *conf, err := readConfigFile(opts.Conf)
	 *if err != nil {
	 *  return nil, err
	 *}
	 */

	log.Info("loading serf agent configuration")
	agentConf := completeAgentConfig()
	if agentConf == nil {
		return nil, fmt.Errorf("failed to load agent configuration")
	}

	log.Info("loading sentinels configuration")
	// TODO loop over conf.List("sentinels")
	hawkeyeConf := SentinelConfig{
		Name: "hawkeye",
		//Sensor:  parsePluginConfig("ping endpoint=http://wzbrzbzbrzb.fr"),
		Sensor:  parsePluginConfig("ping endpoint=whatever.gh"),
		Trigger: parsePluginConfig("clock interval=10s"),
	}

	return &Config{
		Agent:     agentConf,
		Serf:      serf.DefaultConfig(),
		Sentinels: []SentinelConfig{hawkeyeConf},
	}, nil
}
