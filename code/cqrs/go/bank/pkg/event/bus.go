package event

type Bus interface {
	Subscribe(subscriber Subscriber)
	Publish(data ...Event)
}

type Subscriber interface {
	Accept(event Event)
}

type inMemoryEventBus struct {
	subscribers []Subscriber
}

func NewInMemoryEventBus() Bus {
	return &inMemoryEventBus{}
}

func (b *inMemoryEventBus) Subscribe(subscriber Subscriber) {
	b.subscribers = append(b.subscribers, subscriber)
}

func (b *inMemoryEventBus) Publish(events ...Event) {
	for _, event := range events {
		for _, subscriber := range b.subscribers {
			subscriber.Accept(event)
		}
	}
}

type emptyBus struct{}

func NewEmptyBus() Bus {
	return &emptyBus{}
}

func (s *emptyBus) Subscribe(subscriber Subscriber) {
}

func (s *emptyBus) Publish(data ...Event) {
}
