package sentinel

import (
	"bufio"
	"fmt"
	"os"
)

type replyFn func(string, ...interface{})

// SerfAgent maps environment variables set by serf to a convenient structure
// https://www.serfdom.io/docs/agent/event-handlers.html
type SerfAgent struct {
	// Event is the event type that is occurring. This will be one of
	// member-join, member-leave, member-failed, member-update, member-reap,
	// user custom, or query custom
	Event string

	// Payload is data sent along the event
	// TODO parse json
	Payload string

	// Node is the name of the node that is executing the event handler
	Node string

	// Role is the role of the node that is executing the event handler
	Role string

	// Tags is a set for each tag the agent has
	Tags []string

	// LamportTime is the LamportTime of the user event if SERF_EVENT is "user" or "query"
	LamportTime string
}

func readEventPayload() string {
	stdin := bufio.NewReader(os.Stdin)
	line, _, _ := stdin.ReadLine()
	return string(line)
}

func NewSerfAgent() *SerfAgent {
	// TODO search for SERF_TAGS in os.Environ()
	var lt string
	var tags []string

	eventName := os.Getenv("SERF_EVENT")
	if eventName == "user" {
		eventName = os.Getenv("SERF_USER_EVENT")
		lt = os.Getenv("SERF_USER_LTIME")
	} else if eventName == "query" {
		eventName = os.Getenv("SERF_QUERY_NAME")
		lt = os.Getenv("SERF_QUERY_LTIME")
	}
	return &SerfAgent{
		Event:       eventName,
		Payload:     readEventPayload(),
		Node:        os.Getenv("SERF_SELF_NAME"),
		Role:        os.Getenv("SERF_SELF_ROLE"),
		Tags:        tags,
		LamportTime: lt,
	}
}

func (s SerfAgent) Reply(format string, a ...interface{}) {
	// TODO switch writer depending on s.isQuery == true
	fmt.Printf(format, a...)
}
