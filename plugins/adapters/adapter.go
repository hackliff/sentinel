// adapters package stores various plugins aimed at bot communication with an
// operator.
package adapters

import "github.com/azer/logger"

var log = logger.New("sentinel.plugins.adapters")

type Envelope struct {
	Title     string
	Recipient string
}

type Plugin interface {
	Send(Envelope, string) error
}

// NOTE map[string]interface{} for better type support after casting ?
type Creator func(map[string]string) (Plugin, error)

var Plugins = map[string]Creator{}

func Add(name string, creator Creator) {
	Plugins[name] = creator
}
