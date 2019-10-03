package leveldbstore

import (
	"fmt"
	"log"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"

	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event"
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event/provider/inmemorybus"
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event/provider/jsonserializer"
)

const Separator = "\x00"
const EventPrefix = "e" + Separator

type store struct {
	db         *leveldb.DB
	bus        event.Bus
	serializer event.Serializer
}

type Option func(*store)

func WithBus(bus event.Bus) Option {
	return func(store *store) {
		store.bus = bus
	}
}

func WithSerializer(serializer event.Serializer) Option {
	return func(store *store) {
		store.serializer = serializer
	}
}

func New(dbFilePath string, options ...Option) *store {
	db, err := leveldb.OpenFile(dbFilePath, nil)
	if err != nil {
		log.Fatalf("failed opening db: %v", err)
	}

	store := &store{
		db:         db,
		bus:        inmemorybus.New(),
		serializer: jsonserializer.New(),
	}

	for _, option := range options {
		option(store)
	}

	return store
}

func (s *store) Bind(events ...event.Event) {
	s.serializer.Bind(events...)
}

func (s *store) AllEvents() []event.Event {
	return s.eventsByPrefix(EventPrefix)
}

func (s *store) Events(aggregateId string) []event.Event {
	return s.eventsByPrefix(EventPrefix + aggregateId)
}

func (s *store) eventsByPrefix(keyPrefix string) []event.Event {
	var events []event.Event

	iter := s.db.NewIterator(util.BytesPrefix([]byte(keyPrefix)), nil)
	for iter.Next() {
		tokens := strings.Split(string(iter.Key()), Separator)
		eventTypeName := tokens[2]
		value := iter.Value()

		e, err := s.serializer.Deserialize(value, eventTypeName)
		if err != nil {
			log.Fatalf("failed deserializing event: %v", err)
		}
		events = append(events, e)
	}
	iter.Release()

	_ = iter.Error()

	return events
}

func (s *store) Save(aggregateId string, events ...event.Event) error {
	batch := new(leveldb.Batch)
	for _, e := range events {
		data, err := s.serializer.Serialize(e)
		if err != nil {
			return err
		}
		eventTypeName, _ := event.Type(e)
		batch.Put([]byte(EventPrefix+aggregateId+Separator+eventTypeName), data)
	}

	return s.db.Write(batch, nil)
}

func (s *store) Delete() {
	fmt.Println("delete")
}
