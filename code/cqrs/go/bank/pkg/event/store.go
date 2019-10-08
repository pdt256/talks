package event

type Event interface {
	EventTimestamp() int64
}

type Base struct {
	Timestamp int64
}

func (b Base) EventTimestamp() int64 {
	return b.Timestamp
}

type Store interface {
	AllEvents() []Event
	Events(aggregateId string) []Event
	Save(aggregateId string, events ...Event) error
	Bind(events ...Event)
}
