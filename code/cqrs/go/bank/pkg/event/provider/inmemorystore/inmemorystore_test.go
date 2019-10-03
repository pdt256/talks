package inmemorystore_test

import (
	"testing"

	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event"
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event/eventtest"
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event/provider/inmemorybus"
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event/provider/inmemorystore"
)

func TestInMemoryEventStore(t *testing.T) {
	eventtest.VerifyStore(t, newInMemoryStore)
}

func newInMemoryStore() (store event.Store, tearDown func()) {
	return inmemorystore.New(inmemorybus.New()), func() {}
}
