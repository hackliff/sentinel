package config

import (
	"regexp"
	"strings"
)

const PATTERN string = "([a-z0-9]+)=(\"[\\s-:,@a-z0-9]+\")"

type REParser struct {
	RE *regexp.Regexp
}

func NewREParser() *REParser {
	return &REParser{regexp.MustCompile(PATTERN)}
}

// grammar: <plugin type>: kind key1="value" key2="an,array"
// NOTE use plugin-defined spec to cast values and validate ?
func (p REParser) pluginOpts(spec string) map[string]string {
	props := make(map[string]string)
	for _, kvs := range p.RE.FindAllString(spec, -1) {
		parts := strings.Split(kvs, "=")
		props[parts[0]] = strings.Trim(parts[1], "\"")
	}

	return props
}

func (p REParser) pluginName(spec string) string {
	parts := strings.Split(spec, " ")
	return parts[0]
}

func (p REParser) Parse(spec string) *PluginConfig {
	return &PluginConfig{
		Plugin: p.pluginName(spec),
		Opts:   p.pluginOpts(spec),
	}
}
