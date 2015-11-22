package adapters

import (
	"github.com/azer/logger"
	"github.com/hackliff/serf/client"
	"github.com/hackliff/serf/command"
)

// TODO a generic conf structure shared by plugins ?
func get(conf map[string]string, key string, defaultValue interface{}) interface{} {
	res, ok := conf[key]
	if !ok {
		return defaultValue
	}
	return res
}

func init() {
	log.Info("registering adapter: serf")
	Add("serf", func(conf map[string]string) (Plugin, error) {
		// FIXME bad formatted addr (e.g. 127.0.0.1) crashes sentinel
		rpcAddr := get(conf, "addr", "127.0.0.1:7373").(string)
		rpcAuthKey := get(conf, "auth-key", "").(string)
		coalesce := get(conf, "coalesce", true).(bool)

		RPC, err := command.RPCClient(rpcAddr, rpcAuthKey)
		if err != nil {
			return nil, err
		}

		return &Serf{
			RPC:      RPC,
			logger:   logger.New("sentinel.plugins.adapters.serf"),
			coalesce: coalesce,
		}, nil
	})
}

type Serf struct {
	RPC      *client.RPCClient
	logger   *logger.Logger
	coalesce bool
}

func (p Serf) Send(envelope_ Envelope, payload string) error {
	p.logger.Info("emitting event: %s - %s", envelope_.Title, payload)
	return p.RPC.UserEvent(envelope_.Title, []byte(payload), p.coalesce)
}
