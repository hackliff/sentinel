package actuators

import (
	"fmt"
	stdlog "log"
	"os"

	"github.com/hackliff/serf/command/agent"
	"github.com/hackliff/serf/serf"

	"github.com/hackliff/sentinel/plugins/adapters"
)

func init() {
	log.Info("registering actuator: exec")
	Add("exec", func(adapter_ adapters.Plugin, conf map[string]string) (Plugin, error) {
		log.Info("conf: %#v", conf)
		script, ok := conf["script"]
		if !ok {
			return nil, fmt.Errorf("no script provided")
		}

		return &ExecScript{
			Script:  script,
			Adapter: adapter_,
			logger:  stdlog.New(os.Stderr, "", stdlog.LstdFlags),
		}, nil
	})
}

type ExecScript struct {
	Script  string
	Adapter adapters.Plugin
	logger  *stdlog.Logger
}

func (s ExecScript) Description() string {
	return "arbitrary script execution"
}

func (s ExecScript) SampleConfig() string {
	return `
# absolute or agent relative script path
actuator: exec script=./some/where.sh
`
}

func (s ExecScript) Gather(self serf.Member, e serf.Event) error {
	// NOTE use the adapter to alert on error ?
	return agent.InvokeEventScript(s.logger, s.Script, self, e)
}
