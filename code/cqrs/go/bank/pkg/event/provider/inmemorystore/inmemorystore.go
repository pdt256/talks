package inmemorystore

import (
	"sync"

	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event"
)

type store struct {
	bus event.Bus

	mux               *sync.Mutex
	events            []event.Event
	eventsByAggregate map[string][]event.Event
}

func New(bus event.Bus) event.Store {
	return &store{
		bus:               bus,
		mux:               &sync.Mutex{},
		eventsByAggregate: make(map[string][]event.Event),
	}
}

func (s *store) Bind(events ...event.Event) {}

func (s *store) AllEvents() []event.Event {
	return s.events
}

func (s *store) Events(aggregateId string) []event.Event {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.eventsByAggregate[aggregateId]
}

func (s *store) Save(aggregateId string, events ...event.Event) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.events = append(s.events, events...)
	s.eventsByAggregate[aggregateId] = append(s.eventsByAggregate[aggregateId], events...)
	s.bus.Publish(events...)
	return nil
}
