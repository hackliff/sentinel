package actuators

import (
	"fmt"

	"github.com/azer/logger"
	"github.com/hackliff/serf/serf"

	"github.com/hackliff/sentinel/plugins/adapters"
)

func init() {
	log.Info("registering actuator: debugger")
	Add("debug", func(adapter_ adapters.Plugin, conf map[string]string) (Plugin, error) {

		return &Debug{
			Adapter: adapter_,
			logger:  logger.New("sentinel.plugins.actuators.debug"),
		}, nil
	})
}

type Debug struct {
	Adapter adapters.Plugin
	logger  *logger.Logger
}

func (d Debug) Description() string {
	return "print out incoming events"
}

func (d Debug) SampleConfig() string {
	return `
# no parameters !
`
}

// TODO cast event type and perform checkers
func (d Debug) Gather(self serf.Member, event serf.Event) error {
	// debug self member
	d.logger.Info("self name: %s", self.Name)
	d.logger.Info("self tags: %s", self.Tags)

	// debug event
	d.logger.Info("event: %v", event)
	d.logger.Info("event string: %s", event.String())
	d.logger.Info("event type: %s", event.EventType().String())
	// debug specific types
	switch e := event.(type) {
	case serf.MemberEvent:
		d.logger.Info("\tmembers: %v", e.Members)
	case serf.UserEvent:
		d.logger.Info("\tname: %s", e.Name)
		d.logger.Info("\tltime: %v", e.LTime)
		d.logger.Info("\tpayload: %s", string(e.Payload))
		d.logger.Info("\tis coalesced: %v", e.Coalesce)
	case *serf.Query:
		d.logger.Info("\tquery name: %s", e.Name)
		d.logger.Info("\tquery ltime: %s", e.LTime)
	}

	// debug adapter
	envelope := adapters.Envelope{
		Title:     "New Event",
		Recipient: "*",
	}
	d.logger.Info("dispatching information")
	d.Adapter.Send(envelope, fmt.Sprintf("received an event: %s", event.String()))

	return nil
}
