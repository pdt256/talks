package jsoniostream

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"reflect"

	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event"
)

type stream struct {
	eventTypes map[string]reflect.Type
}

func New() *stream {
	return &stream{
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
	Event         interface{} `json:"payload"`
}

func (s *stream) Load(reader io.Reader) <-chan event.Event {
	ch := make(chan event.Event)

	go func() {
		defer close(ch)

		decoder := json.NewDecoder(reader)
		decoder.UseNumber()
		_, err := decoder.Token()
		if err != nil {
			log.Print(err)
			return
		}

		for decoder.More() {
			var rawEvent json.RawMessage
			wrapper := jsonEvent{
				Event: &rawEvent,
			}
			err := decoder.Decode(&wrapper)
			if err != nil {
				log.Printf("failed decoding: %v, %#v", err, wrapper)
				return
			}

			eventType, ok := s.eventTypes[wrapper.EventTypeName]
			if !ok {
				log.Printf("unbound event type, %v", wrapper.EventTypeName)
				continue
			}

			e := reflect.New(eventType).Interface()
			err = json.Unmarshal(rawEvent, e)
			if err != nil {
				log.Printf("failed unmarshalling event: %v", err)
				continue
			}

			ch <- e.(event.Event)
		}
	}()

	return ch
}

func (s *stream) Save(writer io.Writer, events <-chan event.Event) <-chan error {
	errors := make(chan error)
	go func() {
		totalSaved := 0

		_, _ = fmt.Fprint(writer, "[")

		for e := range events {
			eventTypeName, _ := event.Type(e)
			data, err := json.Marshal(jsonEvent{
				EventTypeName: eventTypeName,
				Event:         e,
			})
			if err != nil {
				errors <- fmt.Errorf("failed marshalling jsonEvent: %v", err)
			}

			if totalSaved > 0 {
				_, _ = fmt.Fprint(writer, ",")
			}

			_, _ = fmt.Fprintf(writer, "%s", data)

			totalSaved++
		}
		_, _ = fmt.Fprint(writer, "]")

		close(errors)
	}()

	return errors
}
