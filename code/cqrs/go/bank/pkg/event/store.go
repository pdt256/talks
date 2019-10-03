package event

type Event interface{}

type Store interface {
	AllEvents() []Event
	Events(aggregateId string) []Event
	Save(aggregateId string, events ...Event) error
	Bind(events ...Event)
}
