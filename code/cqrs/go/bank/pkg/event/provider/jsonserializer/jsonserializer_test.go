package jsonserializer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event"
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event/eventtest"
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event/provider/jsonserializer"
)

func TestSerialize(t *testing.T) {
	// Given
	serializer := jsonserializer.New()

	// When
	eventJson, err := serializer.Serialize(eventtest.ThingWasDone{
		Base: event.Base{Timestamp: 1570515152},
		Name: "John",
	})

	// Then
	require.NoError(t, err)
	json := `{"type":"ThingWasDone","payload":{"Timestamp":1570515152,"Name":"John"}}`
	assert.Equal(t, json, string(eventJson))
}

func TestDeserialize(t *testing.T) {
	// Given
	json := `{"type":"ThingWasDone","payload":{"Timestamp":1570515152,"Name":"John"}}`
	serializer := jsonserializer.New()
	serializer.Bind(eventtest.ThingWasDone{})

	// When
	actualEvent, err := serializer.Deserialize([]byte(json))

	// Then
	require.NoError(t, err)
	expectedEvent := &eventtest.ThingWasDone{
		Base: event.Base{Timestamp: 1570515152},
		Name: "John",
	}
	assert.Equal(t, expectedEvent, actualEvent)
}
