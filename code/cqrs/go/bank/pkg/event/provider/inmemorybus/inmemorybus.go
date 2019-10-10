package inmemorybus

import "github.com/pdt256/talks/code/cqrs/go/bank/pkg/event"

type bus struct {
	subscribers []event.Subscriber
}

func New() event.Bus {
	return &bus{}
}

func (b *bus) Subscribe(subscriber ...event.Subscriber) {
	b.subscribers = append(b.subscribers, subscriber...)
}

func (b *bus) Publish(events ...event.Event) {
	for _, e := range events {
		for _, subscriber := range b.subscribers {
			subscriber.Accept(e)
		}
	}
}
