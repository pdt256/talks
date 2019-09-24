package bank

import (
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event"
)

type App struct {
	eventStore           event.Store
	numberAcceptedEvents int
}

func NewApp(eventStore event.Store) *App {
	return &App{eventStore: eventStore}
}

func (a *App) Execute(command interface{}) {
	switch c := command.(type) {

	case OpenAccount:
		a.handleWithAccountAggregate(c.AccountId, command)

	case CloseAccount:
		a.handleWithAccountAggregate(c.AccountId, command)

	case DepositMoney:
		a.handleWithAccountAggregate(c.AccountId, command)

	case WithdrawMoney:
		a.handleWithAccountAggregate(c.AccountId, command)

	}
}

func (a *App) handleWithAccountAggregate(aggregateId string, command interface{}) {
	events := a.eventStore.Events(aggregateId)
	account := NewAccountAggregate(aggregateId, events)
	account.Handle(command)
	a.eventStore.Save(aggregateId, account.PendingEvents...)
}

func (a *App) AcceptEvents(aggregateId string, events ...event.Event) {
	a.numberAcceptedEvents = len(events)
	a.eventStore.Save(aggregateId, events...)
}

func (a *App) AllEmittedEvents() []event.Event {
	return a.eventStore.AllEvents()[a.numberAcceptedEvents:]
}
