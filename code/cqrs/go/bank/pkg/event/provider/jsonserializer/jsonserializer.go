package jsonserializer

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event"
)

type serializer struct {
	eventTypes map[string]reflect.Type
}

func New() *serializer {
	return &serializer{
		eventTypes: map[string]reflect.Type{},
	}
}

func (s *serializer) Bind(events ...event.Event) {
	for _, e := range events {
		eventTypeName, eventType := event.Type(e)
		s.eventTypes[eventTypeName] = eventType
	}
}

type jsonEvent struct {
	EventTypeName string      `json:"type"`
	Event         interface{} `json:"payload"`
}

func (s *serializer) Serialize(e event.Event) ([]byte, error) {
	eventTypeName, _ := event.Type(e)
	data, err := json.Marshal(jsonEvent{
		EventTypeName: eventTypeName,
		Event:         e,
	})
	if err != nil {
		return nil, fmt.Errorf("failed marshalling jsonEvent: %v", err)
	}

	return data, nil
}

func (s *serializer) Deserialize(serializedData []byte) (event.Event, error) {
	var rawEvent json.RawMessage
	wrapper := jsonEvent{
		Event: &rawEvent,
	}
	err := json.Unmarshal(serializedData, &wrapper)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshalling jsonEvent [%s]: %v", serializedData, err)
	}

	eventType, ok := s.eventTypes[wrapper.EventTypeName]
	if !ok {
		return nil, fmt.Errorf("unbound event type, %v", wrapper.EventTypeName)
	}

	e := reflect.New(eventType).Interface()
	err = json.Unmarshal(rawEvent, e)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshalling event: %v", err)
	}

	return e.(event.Event), nil
}
