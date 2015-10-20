package main

import (
	"fmt"
	"log"

	"github.com/hackliff/sentinel"
	"github.com/hackliff/sentinel/radio"
	"github.com/hackliff/sentinel/sensor"
)

// WebSentinel tests HTTP endpoint and triggers alert on bad HTTP codes
type WebSentinel struct {
	Radio  radio.Radio
	Sensor sensor.Sensor
}

func (d WebSentinel) Run(agent sentinel.SerfAgent) error {
	url := agent.Payload
	log.Printf("inspecting webservice : %s\n", url)
	ok, err := d.Sensor.Check(url)

	if !ok {
		// answer serf queries by printing results on stdout
		agent.Reply("check=fail result=%s", err.Error())
		d.Radio.Alert("Nexus 5", "Sentinel Alert", err.Error())
	} else {
		agent.Reply("check=successful result=%t", ok)
	}

	return nil
}

// BuddySentinel watches for his buddy and notifies operator when one of them
// is going missing
type BuddySentinel struct {
	Radio  radio.Radio
	Sensor sensor.Sensor
}

func (d BuddySentinel) Run(agent sentinel.SerfAgent) error {
	log.Printf("a member has failed : %s\n", agent.Payload)
	// TODO parse payload to extract name and address
	content := fmt.Sprintf("member is down : %s)", agent.Payload)
	return d.Radio.Alert("Nexus 5", "Sentinel Alert", content)
}
