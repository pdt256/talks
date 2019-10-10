package main

import (
	"flag"
	"log"
	"os"

	"github.com/pdt256/talks/code/cqrs/go/bank"
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event"
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event/provider/inmemorybus"
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event/provider/jsoniostream"
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/projection"
)

func main() {
	seed := flag.Int64("seed", 100, "random number seed")
	accounts := flag.Int("accounts", 1000, "total number of accounts")
	out := flag.String("out", "pkg/projection/events.json", "output file for json events")
	flag.Parse()

	fakeStream := projection.NewFakeStream(*seed, *accounts)
	jsonStream := jsoniostream.New()
	bank.BindEvents(jsonStream)
	bus := inmemorybus.New()

	file, err := os.Create(*out)
	if err != nil {
		log.Fatal(err)
	}

	events := make(chan event.Event)
	errors := jsonStream.Save(file, events)

	for e := range fakeStream.Load() {
		bus.Publish(e)
		events <- e
	}

	close(events)
	for err := range errors {
		log.Fatal(err)
	}
	file.Close()
}
