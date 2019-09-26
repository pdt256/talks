package projection_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pdt256/talks/code/cqrs/go/bank"
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event"
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/projection"
)

const aggregateId = "E659E9A52BD9438A8200EDC521194DAA"

func TestCountEvents_CalculatesTotalEvents(t *testing.T) {
	// Given
	countEvents := projection.NewCountEvents()
	app := NewTestApp(countEvents)

	// When
	app.AcceptEvents(aggregateId,
		bank.AccountWasOpened{AccountId: "A"},
		bank.MoneyWasDeposited{AccountId: "A", Amount: 100},
		bank.MoneyWasWithdrawn{AccountId: "A", Amount: 75},
		bank.MoneyWasWithdrawn{AccountId: "A", Amount: 25},
		bank.AccountWasClosed{AccountId: "A"},
	)

	// Then
	assert.Equal(t, 5, countEvents.TotalEvents)
}

func TestAccountFunds_CalculatesTotalFundsAndAccountBalances(t *testing.T) {
	// Given
	accountFunds := projection.NewAccountFunds()
	app := NewTestApp(accountFunds)

	// When
	app.AcceptEvents(aggregateId,
		bank.AccountWasOpened{AccountId: "A"},
		bank.MoneyWasDeposited{AccountId: "A", Amount: 100},
		bank.MoneyWasWithdrawn{AccountId: "A", Amount: 50},
		bank.AccountWasOpened{AccountId: "B"},
		bank.MoneyWasDeposited{AccountId: "B", Amount: 25},
	)

	// Then
	assert.Equal(t, 75, accountFunds.TotalFunds)
	assert.Equal(t, 50, accountFunds.AccountBalance["A"])
	assert.Equal(t, 25, accountFunds.AccountBalance["B"])
}

func NewTestApp(subscriber event.Subscriber) *bank.App {
	bus := event.NewInMemoryEventBus()
	bus.Subscribe(subscriber)
	eventStore := event.NewInMemoryEventStore(bus)
	return bank.NewApp(eventStore)
}
