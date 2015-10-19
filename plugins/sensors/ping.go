package sensors

import "fmt"

func init() {
	log.Info("\t- registering ping sensor")
	Add("ping", func(conf map[string]string) (SensorPlugin, error) {
		endpoint, ok := conf["endpoint"]
		if !ok {
			return nil, fmt.Errorf("no endpoint provided")
		}
		return &PingSensor{
			Endpoint: endpoint,
		}, nil
	})
}

type PingSensor struct {
	// NOTE multiple endpoints ?
	Endpoint string
}

func (p PingSensor) Description() string {
	return "sensor monitoring ping responses from http endpoint"
}

func (p PingSensor) SampleConfig() string {
	return `
# endpoint to ping
endpoint = "http://example.com"
`
}

func (p PingSensor) Monitor() error {
	// TODO check http response
	log.Info("inspecting HTTP endpoint: %s", p.Endpoint)
	return nil
}
