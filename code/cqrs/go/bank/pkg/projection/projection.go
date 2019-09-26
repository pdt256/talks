package projection

import (
	"github.com/pdt256/talks/code/cqrs/go/bank"
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event"
)

type countEvents struct {
	TotalEvents int
}

func NewCountEvents() *countEvents {
	return &countEvents{}
}

func (c *countEvents) Accept(event event.Event) {
	c.TotalEvents++
}

type accountFunds struct {
	TotalFunds     int
	AccountBalance map[string]int
}

func (a *accountFunds) Accept(event event.Event) {
	switch e := event.(type) {

	case bank.MoneyWasDeposited:
		a.TotalFunds += e.Amount
		a.AccountBalance[e.AccountId] += e.Amount

	case bank.MoneyWasWithdrawn:
		a.TotalFunds -= e.Amount
		a.AccountBalance[e.AccountId] -= e.Amount

	}
}

func NewAccountFunds() *accountFunds {
	return &accountFunds{
		AccountBalance: make(map[string]int),
	}
}
