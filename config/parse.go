package config

import (
	"strings"

	"github.com/olebedev/config"
)

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
