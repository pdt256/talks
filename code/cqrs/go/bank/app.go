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
		a.eventStore.Save(c.AccountId,
			AccountWasOpened{
				AccountId: c.AccountId,
			})

	case CloseAccount:
		account := a.GetAccountAggregate(c.AccountId)
		if account.Balance() > 0 {
			a.eventStore.Save(c.AccountId,
				FailedToCloseAccountWithBalance{
					AccountId: c.AccountId,
					Balance:   account.Balance(),
				})
			return
		}

		a.eventStore.Save(c.AccountId,
			AccountWasClosed{
				AccountId: c.AccountId,
			})

	case DepositMoney:
		account := a.GetAccountAggregate(c.AccountId)
		if account.IsClosed() {
			a.eventStore.Save(c.AccountId,
				DepositFailed{
					AccountId: c.AccountId,
					Amount:    c.Amount,
				})
			return
		}

		a.eventStore.Save(c.AccountId,
			MoneyWasDeposited{
				AccountId: c.AccountId,
				Amount:    c.Amount,
			})

	case WithdrawMoney:
		account := a.GetAccountAggregate(c.AccountId)
		if account.Balance() < c.Amount {
			a.eventStore.Save(c.AccountId,
				WithdrawDenied{
					AccountId:      c.AccountId,
					Amount:         c.Amount,
					CurrentBalance: account.Balance(),
				})
			return
		}

		a.eventStore.Save(c.AccountId,
			MoneyWasWithdrawn{
				AccountId:  c.AccountId,
				Amount:     c.Amount,
				NewBalance: account.Balance() - c.Amount,
			})

	}
}

func (a *App) AcceptEvents(aggregateId string, events ...event.Event) {
	a.numberAcceptedEvents = len(events)
	a.eventStore.Save(aggregateId, events...)
}

func (a *App) AllEmittedEvents() []event.Event {
	return a.eventStore.AllEvents()[a.numberAcceptedEvents:]
}

func (a *App) GetAccountAggregate(aggregateId string) *AccountAggregate {
	events := a.eventStore.Events(aggregateId)
	return NewAccountAggregate(events)
}
