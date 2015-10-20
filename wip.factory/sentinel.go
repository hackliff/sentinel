package sentinel

import (
	"log"
)

type Sentinel interface {
	//Run() (bool, error)
	Run(SerfAgent) error
}

type Squad struct {
	Agent    SerfAgent
	handlers map[string]Sentinel
}

func NewSquad() *Squad {
	return &Squad{
		Agent:    *NewSerfAgent(),
		handlers: make(map[string]Sentinel),
	}
}

func (s Squad) Register(event string, sentinel Sentinel) {
	// NOTE allow an array of handlers for one event ?
	// NOTE what about one handler for multiple events ?
	log.Printf("registering handler for event '%s'\n", event)
	s.handlers[event] = sentinel
}

func (s *Squad) Dispatch() error {
	// lookup sentinels registered
	handler, ok := s.handlers[s.Agent.Event]
	if ok {
		log.Printf("processing event %v\n", s.Agent.Event)
		return handler.Run(s.Agent)
	}

	log.Printf("no handler found for event '%s'\n", s.Agent.Event)
	return nil
}
