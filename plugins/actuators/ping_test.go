package actuators

import "testing"

const testEndpoint = "http://google.com"
const invalidTestEndpoint = testEndpoint + "/whatever"

func TestPingActuator(t *testing.T) {
	var _ Plugin = &Ping{}
}
