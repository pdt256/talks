package bank

import (
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event"
)

type AccountAggregate struct {
	state struct {
		isOpen  bool
		balance int
	}
}

func NewAccountAggregate(events []event.Event) *AccountAggregate {
	aggregate := &AccountAggregate{}

	for _, e := range events {
		aggregate.transition(e)
	}

	return aggregate
}

func (a *AccountAggregate) transition(e event.Event) {
	switch e := e.(type) {

	case AccountWasOpened:
		a.state.isOpen = true

	case DepositMoney:
		a.state.balance += e.Amount
	}
}

func (a *AccountAggregate) IsClosed() bool {
	return !a.state.isOpen
}

func (a *AccountAggregate) Balance() int {
	return a.state.balance
}
