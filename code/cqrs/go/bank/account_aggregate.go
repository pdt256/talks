package bank

import (
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event"
)

type AccountAggregate struct {
	PendingEvents []event.Event
	state         struct {
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

	case WithdrawMoney:
		a.state.balance -= e.Amount

	}
}

func (a *AccountAggregate) Handle(command interface{}) {
	switch c := command.(type) {

	case OpenAccount:
		a.emitEvent(AccountWasOpened{
			AccountId: c.AccountId,
		})

	case CloseAccount:
		if a.state.balance > 0 {
			a.emitEvent(FailedToCloseAccountWithBalance{
				AccountId: c.AccountId,
				Balance:   a.state.balance,
			})
			return
		}

		a.emitEvent(AccountWasClosed{
			AccountId: c.AccountId,
		})

	case DepositMoney:
		if !a.state.isOpen {
			a.emitEvent(DepositFailed{
				AccountId: c.AccountId,
				Amount:    c.Amount,
			})
			return
		}

		a.emitEvent(MoneyWasDeposited{
			AccountId:  c.AccountId,
			Amount:     c.Amount,
			NewBalance: a.state.balance + c.Amount,
		})

	case WithdrawMoney:
		if a.state.balance < c.Amount {
			a.emitEvent(WithdrawDenied{
				AccountId:      c.AccountId,
				Amount:         c.Amount,
				CurrentBalance: a.state.balance,
			})
			return
		}

		a.emitEvent(MoneyWasWithdrawn{
			AccountId:  c.AccountId,
			Amount:     c.Amount,
			NewBalance: a.state.balance - c.Amount,
		})

	}
}

func (a *AccountAggregate) emitEvent(event ...event.Event) {
	a.PendingEvents = append(a.PendingEvents, event...)
}
