package jsoniostream

import (
	"encoding/json"
	"io"
	"log"
	"reflect"

	"github.com/mitchellh/mapstructure"

	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event"
)

type stream struct {
	reader     io.Reader
	eventTypes map[string]reflect.Type
}

func New(reader io.Reader) *stream {
	return &stream{
		reader:     reader,
		eventTypes: map[string]reflect.Type{},
	}
}

func (s *stream) Bind(events ...event.Event) {
	for _, e := range events {
		eventTypeName, eventType := event.Type(e)
		s.eventTypes[eventTypeName] = eventType
	}
}

type jsonEvent struct {
	EventTypeName string      `json:"type"`
	Event         event.Event `json:"payload"`
}

func (s *stream) Load() <-chan event.Event {
	ch := make(chan event.Event)

	go func() {
		defer close(ch)

		decoder := json.NewDecoder(s.reader)
		_, err := decoder.Token()
		if err != nil {
			log.Print(err)
			return
		}

		for decoder.More() {
			wrapper := jsonEvent{}
			err := decoder.Decode(&wrapper)
			if err != nil {
				log.Print(err)
				return
			}

			eventType, ok := s.eventTypes[wrapper.EventTypeName]
			if !ok {
				log.Printf("unbound event type, %v", wrapper.EventTypeName)
				continue
			}

			e := reflect.New(eventType).Interface()
			err = mapstructure.Decode(wrapper.Event, e)
			if err != nil {
				log.Printf("failed unmarshalling event: %v", err)
				continue
			}

			ch <- e
		}
	}()

	return ch
}
