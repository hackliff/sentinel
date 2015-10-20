package sensors

import (
	"fmt"

	"github.com/hashicorp/serf/client"
	"github.com/levigross/grequests"
)

// TODO paramteric ? useful ?
const COALESCE bool = true

func init() {
	log.Info("\t- registering ping sensor")
	Add("ping", func(RPC *client.RPCClient, conf map[string]string) (SensorPlugin, error) {
		endpoint, ok := conf["endpoint"]
		if !ok {
			return nil, fmt.Errorf("no endpoint provided")
		}
		return &PingSensor{
			Endpoint: endpoint,
			RPC:      RPC,
		}, nil
	})
}

type PingSensor struct {
	// NOTE multiple endpoints ?
	Endpoint string
	RPC      *client.RPCClient
}

func (p PingSensor) Description() string {
	return "sensor monitoring ping responses from http endpoint"
}

func (p PingSensor) SampleConfig() string {
	return `
# endpoint to inspect
endpoint="http://example.com"
`
}

// NOTE should be easier (less boilerplate)
func (p PingSensor) notify(pingErr error) error {
	log.Error("failed to ping endpoint: %s", pingErr)
	// TODO proper protocol
	event := "sensor-failed"
	payload := fmt.Sprintf("endpoint=%s sensor=ping err=%s", p.Endpoint, pingErr)
	if err := p.RPC.UserEvent(event, []byte(payload), COALESCE); err != nil {
		log.Error("Error sending event: %s", err)
		return err
	}
	log.Info("Event '%s' dispatched! Coalescing enabled: %#v", event, COALESCE)
	return pingErr
}

func (p PingSensor) Monitor() error {
	// TODO check http response
	log.Info("inspecting HTTP endpoint: %s", p.Endpoint)
	resp, err := grequests.Get(p.Endpoint, nil)
	if err != nil || !resp.Ok {
		return p.notify(err)
	}

	return nil
}
