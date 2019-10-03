package leveldbstore_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event"
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event/eventtest"
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event/provider/inmemorybus"
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event/provider/jsonserializer"
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event/provider/leveldbstore"
)

func TestLeveldbEventStore(t *testing.T) {
	eventtest.VerifyStore(t, newLevelDbStore)
}

var dbCount int

func newLevelDbStore() (event.Store, func()) {
	dbCount++
	dbPath := filepath.Join(os.TempDir(), fmt.Sprintf("testevents-%d-%d", os.Getuid(), dbCount))
	store := leveldbstore.New(dbPath,
		leveldbstore.WithBus(inmemorybus.New()),
		leveldbstore.WithSerializer(jsonserializer.New()),
	)

	teardown := func() {
		err := os.RemoveAll(dbPath)
		if err != nil {
			log.Fatalf("unable to teardown db: %v", err)
		}
	}

	return store, teardown
}
