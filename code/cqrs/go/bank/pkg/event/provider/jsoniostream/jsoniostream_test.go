package jsoniostream_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event/provider/jsoniostream"
)

func TestLoad_NoEvents(t *testing.T) {
	// Given
	stream := jsoniostream.New(strings.NewReader(""))

	// When
	eventChannel := stream.Load()

	// Then
	end := <-eventChannel
	assert.Equal(t, nil, end)
}

func TestLoad_OneEventInStream(t *testing.T) {
	// Given
	type ThingWasDone struct{ Name string }
	jsonData := `[{"type":"ThingWasDone","payload":{"Name":"John"}}]`
	reader := strings.NewReader(jsonData)
	stream := jsoniostream.New(reader)
	stream.Bind(ThingWasDone{})

	// When
	eventChannel := stream.Load()

	// Then
	actualEvent := <-eventChannel
	assert.Equal(t, &ThingWasDone{Name: "John"}, actualEvent)
}

func TestLoad_TwoEventsInStream(t *testing.T) {
	// Given
	type ThingWasDone struct{ Name string }
	jsonData := `[{"type":"ThingWasDone","payload":{"Name":"John"}},{"type":"ThingWasDone","payload":{"Name":"Jane"}}]`
	reader := strings.NewReader(jsonData)
	stream := jsoniostream.New(reader)
	stream.Bind(ThingWasDone{})

	// When
	eventChannel := stream.Load()

	// Then
	actualEvent1 := <-eventChannel
	actualEvent2 := <-eventChannel
	end := <-eventChannel
	assert.Equal(t, &ThingWasDone{Name: "John"}, actualEvent1)
	assert.Equal(t, &ThingWasDone{Name: "Jane"}, actualEvent2)
	assert.Equal(t, nil, end)
}
