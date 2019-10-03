package eventtest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event"
)

func VerifyStore(t *testing.T, newStore func() (store event.Store, tearDown func())) {
	type AccountWasOpened struct {
		Name string
	}

	t.Run("save event to store by aggregate id", func(t *testing.T) {
		// Given
		store, tearDown := newStore()
		defer tearDown()
		store.Bind(&AccountWasOpened{})
		event1 := &AccountWasOpened{Name: "John Doe"}
		event2 := &AccountWasOpened{Name: "Jane Doe"}

		// When
		_ = store.Save("1", event1)
		_ = store.Save("2", event2)

		// Then
		events := store.AllEvents()
		require.Equal(t, 2, len(events))
		assert.Equal(t, []event.Event{event1, event2}, events)
	})

	t.Run("get events from store by aggregate id", func(t *testing.T) {
		// Given
		store, tearDown := newStore()
		defer tearDown()
		store.Bind(&AccountWasOpened{})
		event1 := &AccountWasOpened{Name: "John Doe"}
		event2 := &AccountWasOpened{Name: "Jane Doe"}

		// When
		_ = store.Save("3", event1)
		_ = store.Save("4", event2)

		// Then
		events := store.Events("3")
		require.Equal(t, 1, len(events))
		assert.Equal(t, []event.Event{event1}, events)
	})
}
