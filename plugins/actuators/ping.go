package actuators

import (
	"fmt"

	"github.com/azer/logger"
	"github.com/hashicorp/serf/serf"
	"github.com/levigross/grequests"

	"github.com/hackliff/sentinel/plugins/adapters"
)

func init() {
	log.Info("registering actuator: ping")
	Add("ping", func(adapter_ adapters.Plugin, conf map[string]string) (Plugin, error) {
		endpoint, ok := conf["endpoint"]
		if !ok {
			return nil, fmt.Errorf("no endpoint provided")
		}
		return &Ping{
			Endpoint: endpoint,
			Adapter:  adapter_,
			logger:   logger.New("sentinel.plugins.actuators.ping"),
		}, nil
	})
}

type Ping struct {
	// NOTE multiple endpoints ?
	Endpoint string
	Adapter  adapters.Plugin
	logger   *logger.Logger
}

func (p Ping) Description() string {
	return "actuator monitoring ping responses from http endpoint"
}

func (p Ping) SampleConfig() string {
	return `
# endpoint to inspect
actuator: ping endpoint="http://example.com"
`
}

// NOTE should be easier (less boilerplate)
func (p Ping) notify(pingErr error) error {
	envelope_ := adapters.Envelope{
		Title:     "actuator-failed",
		Recipient: "*",
	}
	// TODO proper protocol ?
	payload := fmt.Sprintf("endpoint=%s actuator=ping err=%s", p.Endpoint, pingErr)
	if err := p.Adapter.Send(envelope_, payload); err != nil {
		p.logger.Error("Error sending event: %s", err)
		return err
	}
	p.logger.Info("Event '%s' dispatched", envelope_.Title)
	return pingErr
}

func (p Ping) Gather(self serf.Member, e serf.Event) error {
	// TODO check http response
	p.logger.Info("inspecting HTTP endpoint: %s", p.Endpoint)
	resp, err := grequests.Get(p.Endpoint, nil)
	if err != nil || !resp.Ok {
		p.logger.Error("failed to ping endpoint: %s", err)
		return p.notify(err)
	}

	return nil
}
