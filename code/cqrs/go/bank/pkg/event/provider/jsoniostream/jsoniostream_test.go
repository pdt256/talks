package jsoniostream_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event"
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event/eventtest"
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event/provider/jsoniostream"
)

func TestLoad_NoJson(t *testing.T) {
	// Given
	reader := strings.NewReader("")
	stream := jsoniostream.New()

	// When
	events := stream.Load(reader)

	// Then
	end := <-events
	assert.Equal(t, nil, end)
}

func TestLoad_EmptyJsonEvents(t *testing.T) {
	// Given
	reader := strings.NewReader("[]")
	stream := jsoniostream.New()

	// When
	events := stream.Load(reader)

	// Then
	assert.Equal(t, nil, <-events)
}

func TestLoad_OneEventInStream(t *testing.T) {
	// Given
	json := `[{"type":"ThingWasDone","payload":{"Timestamp":1570515152,"Name":"John"}}]`
	reader := strings.NewReader(json)
	stream := jsoniostream.New()
	stream.Bind(eventtest.ThingWasDone{})

	// When
	events := stream.Load(reader)

	// Then
	actualEvent := <-events
	assert.Equal(t, &eventtest.ThingWasDone{
		Base: event.Base{Timestamp: 1570515152},
		Name: "John",
	}, actualEvent)
}

func TestLoad_TwoEventsInStream(t *testing.T) {
	// Given
	jsonData := `[{"type":"ThingWasDone","payload":{"Name":"John"}},{"type":"ThingWasDone","payload":{"Name":"Jane"}}]`
	reader := strings.NewReader(jsonData)
	stream := jsoniostream.New()
	stream.Bind(eventtest.ThingWasDone{})

	// When
	events := stream.Load(reader)

	// Then
	actualEvent1 := <-events
	actualEvent2 := <-events
	end := <-events
	assert.Equal(t, &eventtest.ThingWasDone{Name: "John"}, actualEvent1)
	assert.Equal(t, &eventtest.ThingWasDone{Name: "Jane"}, actualEvent2)
	assert.Equal(t, nil, end)
}

func TestSave_NoEvents_ReturnsEmptyJsonList(t *testing.T) {
	// Given
	var writer bytes.Buffer
	stream := jsoniostream.New()
	events := make(chan event.Event)

	// When
	errors := stream.Save(&writer, events)

	// Then
	close(events)
	assert.Equal(t, nil, <-errors)
	assert.Equal(t, "[]", writer.String())
}

func TestSave_OneEventInStream(t *testing.T) {
	// Given
	var writer bytes.Buffer
	stream := jsoniostream.New()
	stream.Bind(eventtest.ThingWasDone{})
	events := make(chan event.Event)

	// When
	errors := stream.Save(&writer, events)

	// Then
	events <- eventtest.ThingWasDone{
		Base: event.Base{Timestamp: 1570515152},
		Name: "John",
	}
	close(events)
	lastError := <-errors
	json := `[{"type":"ThingWasDone","payload":{"Timestamp":1570515152,"Name":"John"}}]`
	assert.Equal(t, json, writer.String())
	assert.Equal(t, nil, lastError)
}

func TestSave_TwoEventsInStream(t *testing.T) {
	// Given
	var writer bytes.Buffer
	stream := jsoniostream.New()
	stream.Bind(eventtest.ThingWasDone{})
	events := make(chan event.Event)

	// When
	errors := stream.Save(&writer, events)

	// Then
	events <- eventtest.ThingWasDone{
		Base: event.Base{Timestamp: 1570515152},
		Name: "John",
	}
	events <- eventtest.ThingWasDone{
		Base: event.Base{Timestamp: 1570515153},
		Name: "Jane",
	}
	close(events)
	lastError := <-errors
	json := `[{"type":"ThingWasDone","payload":{"Timestamp":1570515152,"Name":"John"}},{"type":"ThingWasDone","payload":{"Timestamp":1570515153,"Name":"Jane"}}]`
	assert.Equal(t, json, writer.String())
	assert.Equal(t, nil, lastError)
}
