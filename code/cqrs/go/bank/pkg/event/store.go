package event

type Event interface{}

type Store interface {
	AllEvents() []Event
	Events(aggregateId string) []Event
	Save(aggregateId string, event ...Event)
}

type InMemoryEventStore struct {
	events            []Event
	eventsByAggregate map[string][]Event
	bus               Bus
}

func NewInMemoryEventStore(bus Bus) Store {
	return &InMemoryEventStore{
		eventsByAggregate: make(map[string][]Event),
		bus:               bus,
	}
}

func (s *InMemoryEventStore) AllEvents() []Event {
	return s.events
}

func (s *InMemoryEventStore) Events(aggregateId string) []Event {
	return s.eventsByAggregate[aggregateId]
}

func (s *InMemoryEventStore) Save(aggregateId string, event ...Event) {
	s.events = append(s.events, event...)
	s.eventsByAggregate[aggregateId] = append(s.eventsByAggregate[aggregateId], event...)
	s.bus.Publish(event...)
}
