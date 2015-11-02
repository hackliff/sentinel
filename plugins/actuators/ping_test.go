package actuators

import "testing"

var testEndpoint string = "http://google.com"
var invalidTestEndpoint string = testEndpoint + "/whatever"

func TestPingSensor_implements(t *testing.T) {
	var _ SensorPlugin = &PingSensor{}
}

func TestPingSensorCheck(t *testing.T) {
	sensor := &PingSensor{Endpoint: testEndpoint}
	_, err := sensor.Monitor()
	if err != nil {
		t.Errorf("Ping sensor check errored: %v\n", err)
	}
}

func TestPingCheckValidEndpoint(t *testing.T) {
	sensor := &PingSensor{Endpoint: testEndpoint}
	ok, _ := sensor.Monitor()
	if !ok {
		t.Errorf("Ping sensor check failed\n")
	}
}

func TestPingCheckUnavailableEndpoint(t *testing.T) {
	sensor := &PingSensor{Endpoint: invalidTestEndpoint}
	ok, _ := sensor.Monitor()
	if ok {
		t.Errorf("Ping sensor should have failed\n")
	}
}
