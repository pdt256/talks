package projection_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pdt256/talks/code/cqrs/go/bank"
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event/provider/inmemorybus"
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/projection"
)

func TestCountEvents_CalculatesTotals(t *testing.T) {
	// Given
	countEvents := projection.NewCountEvents()
	bus := inmemorybus.New()
	bus.Subscribe(countEvents)

	// When
	bus.Publish(
		bank.AccountWasOpened{AccountId: "A"},
		bank.MoneyWasDeposited{AccountId: "A", Amount: 100},
		bank.MoneyWasWithdrawn{AccountId: "A", Amount: 75},
		bank.MoneyWasWithdrawn{AccountId: "A", Amount: 25},
		bank.AccountWasClosed{AccountId: "A"},
	)

	// Then
	assert.Equal(t, 5, countEvents.TotalEvents)
	expectedCounts := map[string]int{
		"AccountWasClosed":  1,
		"AccountWasOpened":  1,
		"MoneyWasDeposited": 1,
		"MoneyWasWithdrawn": 2,
	}
	assert.Equal(t, expectedCounts, countEvents.Counts)
}
