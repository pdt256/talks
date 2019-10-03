package inmemorybus_test

import (
	"testing"

	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event/eventtest"
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event/provider/inmemorybus"
)

func TestInMemoryEventBus(t *testing.T) {
	eventtest.VerifyBus(t, inmemorybus.New)
}
