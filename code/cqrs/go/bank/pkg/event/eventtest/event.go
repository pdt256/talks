package eventtest

import "github.com/pdt256/talks/code/cqrs/go/bank/pkg/event"

type ThingWasDone struct {
	event.Base
	Name string
}
