package eventtest

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event"
)

func VerifyBus(t *testing.T, newBus func() (bus event.Bus)) {
	t.Run("publish events to subscriber", func(t *testing.T) {
		// Given
		subscriber := &loggingSubscriber{}
		bus := newBus()
		bus.Subscribe(subscriber)
		event1 := ThingWasDone{Name: "John Doe"}
		event2 := ThingWasDone{Name: "Jane Doe"}

		// When
		bus.Publish(event1, event2)

		// Then
		assert.Equal(t, []event.Event{event1, event2}, subscriber.eventsSeen)
	})
}

type loggingSubscriber struct {
	eventsSeen []event.Event
}

func (s *loggingSubscriber) Accept(event event.Event) {
	s.eventsSeen = append(s.eventsSeen, event)
}
