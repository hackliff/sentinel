package actuators

import (
	"github.com/azer/logger"
	"github.com/hashicorp/serf/serf"
	"github.com/yuin/gopher-lua"

	"github.com/hackliff/sentinel/plugins/adapters"
)

const DEFAULT_LUA_SCRIPT_PATH = "./script.lua"

func init() {
	log.Info("registering actuator: lua")
	Add("lua", func(adapter_ adapters.Plugin, conf map[string]string) (Plugin, error) {
		script, ok := conf["script"]
		if !ok {
			script = DEFAULT_LUA_SCRIPT_PATH
			// TODO check it exists
		}

		vm := lua.NewState()
		// NOTE defer vm.Close() ?

		return &LuaScript{
			Script:  script,
			VM:      vm,
			Adapter: adapter_,
			logger:  logger.New("sentinel.plugins.actuators.ping"),
		}, nil
	})
}

type LuaScript struct {
	VM      *lua.LState
	Script  string
	Adapter adapters.Plugin
	logger  *logger.Logger
}

func (l LuaScript) Description() string {
	return "lua script execution"
}

func (l LuaScript) SampleConfig() string {
	return `
# absolute or agent relative script path
actuator: lua script=./somewhere.lua
`
}

func (l LuaScript) Gather(self serf.Member, e serf.Event) error {
	return l.VM.DoFile(l.Script)
}
