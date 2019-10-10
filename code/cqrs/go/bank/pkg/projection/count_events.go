package projection

import (
	"reflect"

	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event"
)

type countEvents struct {
	TotalEvents int
	Counts      map[string]int
}

func NewCountEvents() *countEvents {
	return &countEvents{
		Counts: make(map[string]int),
	}
}

func (c *countEvents) Accept(e event.Event) {
	c.TotalEvents++
	eventName := reflect.TypeOf(e).Name()
	c.Counts[eventName]++
}
