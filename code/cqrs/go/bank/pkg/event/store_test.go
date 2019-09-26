package event_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event"
)

func TestInMemoryEventStore(t *testing.T) {
	VerifyStoreInterface(t, event.NewInMemoryEventStore, event.NewEmptyBus)
}

func VerifyStoreInterface(t *testing.T, newStore func(bus event.Bus) event.Store, newBus func() event.Bus) {
	type AccountWasOpened struct {
		Name string
	}

	t.Run("save event to store by aggregate id", func(t *testing.T) {
		// Given
		store := newStore(newBus())
		event1 := AccountWasOpened{Name: "John Doe"}
		event2 := AccountWasOpened{Name: "Jane Doe"}

		// When
		store.Save("1", event1)
		store.Save("2", event2)

		// Then
		events := store.AllEvents()
		require.Equal(t, 2, len(events))
		assert.Equal(t, []event.Event{event1, event2}, events)
	})

	t.Run("get events from store by aggregate id", func(t *testing.T) {
		// Given
		store := newStore(newBus())
		event1 := AccountWasOpened{Name: "John Doe"}
		event2 := AccountWasOpened{Name: "Jane Doe"}

		// When
		store.Save("1", event1)
		store.Save("2", event2)

		// Then
		events := store.Events("1")
		require.Equal(t, 1, len(events))
		assert.Equal(t, []event.Event{event1}, events)
	})
}
