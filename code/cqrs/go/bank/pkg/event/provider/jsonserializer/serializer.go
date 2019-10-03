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

func (s *serializer) Serialize(event event.Event) ([]byte, error) {
	data, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("failed marshalling event: %v", err)
	}
	return data, nil
}

func (s *serializer) Deserialize(serializedData []byte, eventTypeName string) (event.Event, error) {
	eventType, ok := s.eventTypes[eventTypeName]
	if !ok {
		return nil, fmt.Errorf("unbound event type, %v", eventTypeName)
	}

	e := reflect.New(eventType).Interface()
	err := json.Unmarshal(serializedData, e)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshalling event: %v", err)
	}

	return e, nil
}
