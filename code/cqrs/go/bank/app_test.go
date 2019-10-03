package bank_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pdt256/talks/code/cqrs/go/bank"
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event"
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event/provider/inmemorybus"
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event/provider/inmemorystore"
)

const accountId = "7A22F482897142ED84CDC15B02C75948"
const aggregateId = accountId

func TestOpenAccount_Emits_AccountWasOpened(t *testing.T) {
	// Given
	app := NewTestApp()

	// When
	app.Execute(
		bank.OpenAccount{
			AccountId: accountId,
		})

	// Then
	ExpectEmittedEvents(t, app,
		bank.AccountWasOpened{
			AccountId: accountId,
		})
}

func TestDepositMoney_WhenNoAccountExists_Emits_DepositFailed(t *testing.T) {
	// Given
	app := NewTestApp()

	// When
	app.Execute(
		bank.DepositMoney{
			AccountId: accountId,
			Amount:    100,
		})

	// Then
	ExpectEmittedEvents(t, app,
		bank.DepositFailed{
			AccountId: accountId,
			Amount:    100,
		})
}

func TestDepositMoney_WhenAccountExists_Emits_MoneyWasDeposited(t *testing.T) {
	// Given
	app := NewTestApp()
	app.AcceptEvents(
		aggregateId,
		bank.AccountWasOpened{
			AccountId: accountId,
		},
		bank.MoneyWasDeposited{
			AccountId: accountId,
			Amount:    50,
		})

	// When
	app.Execute(
		bank.DepositMoney{
			AccountId: accountId,
			Amount:    100,
		})

	// Then
	ExpectEmittedEvents(t, app,
		bank.MoneyWasDeposited{
			AccountId:  accountId,
			Amount:     100,
			NewBalance: 150,
		})
}

func TestWithdrawMoney_WhenFundsAreNotAvailable_Emits_WithdrawDenied(t *testing.T) {
	// Given
	app := NewTestApp()
	app.AcceptEvents(
		aggregateId,
		bank.AccountWasOpened{
			AccountId: accountId,
		},
		bank.MoneyWasDeposited{
			AccountId: accountId,
			Amount:    50,
		})

	// When
	app.Execute(
		bank.WithdrawMoney{
			AccountId: accountId,
			Amount:    100,
		})

	// Then
	ExpectEmittedEvents(t, app,
		bank.WithdrawDenied{
			AccountId:      accountId,
			Amount:         100,
			CurrentBalance: 50,
		})
}

func TestWithdrawMoney_WhenFundsAreNotAvailableAfterPreviousWithdrawal_Emits_WithdrawDenied(t *testing.T) {
	// Given
	app := NewTestApp()
	app.AcceptEvents(
		aggregateId,
		bank.AccountWasOpened{
			AccountId: accountId,
		},
		bank.MoneyWasDeposited{
			AccountId: accountId,
			Amount:    50,
		},
		bank.MoneyWasWithdrawn{
			AccountId: accountId,
			Amount:    40,
		})

	// When
	app.Execute(
		bank.WithdrawMoney{
			AccountId: accountId,
			Amount:    25,
		})

	// Then
	ExpectEmittedEvents(t, app,
		bank.WithdrawDenied{
			AccountId:      accountId,
			Amount:         25,
			CurrentBalance: 10,
		})
}

func TestWithdrawMoney_WhenFundsAreAvailable_Emits_MoneyWasWithdrawn(t *testing.T) {
	// Given
	app := NewTestApp()
	app.AcceptEvents(
		aggregateId,
		bank.AccountWasOpened{
			AccountId: accountId,
		},
		bank.MoneyWasDeposited{
			AccountId: accountId,
			Amount:    100,
		})

	// When
	app.Execute(
		bank.WithdrawMoney{
			AccountId: accountId,
			Amount:    75,
		})

	// Then
	ExpectEmittedEvents(t, app,
		bank.MoneyWasWithdrawn{
			AccountId:  accountId,
			Amount:     75,
			NewBalance: 25,
		})
}

func TestCloseAccount_WhenFundsAreStillInAccount_Emits_FailedToCloseAccountWithBalance(t *testing.T) {
	// Given
	app := NewTestApp()
	app.AcceptEvents(
		aggregateId,
		bank.AccountWasOpened{
			AccountId: accountId,
		},
		bank.MoneyWasDeposited{
			AccountId: accountId,
			Amount:    100,
		})

	// When
	app.Execute(
		bank.CloseAccount{
			AccountId: accountId,
		})

	// Then
	ExpectEmittedEvents(t, app,
		bank.FailedToCloseAccountWithBalance{
			AccountId: accountId,
			Balance:   100,
		})
}

func TestCloseAccount_WithEmptyAccount_Emits_AccountWasClosed(t *testing.T) {
	// Given
	app := NewTestApp()
	app.AcceptEvents(
		aggregateId,
		bank.AccountWasOpened{
			AccountId: accountId,
		})

	// When
	app.Execute(
		bank.CloseAccount{
			AccountId: accountId,
		})

	// Then
	ExpectEmittedEvents(t, app,
		bank.AccountWasClosed{
			AccountId: accountId,
		})
}

func NewTestApp() *bank.App {
	return bank.NewApp(inmemorystore.New(inmemorybus.New()))
}

func ExpectEmittedEvents(t *testing.T, app *bank.App, expectedEvents ...event.Event) {
	t.Helper()

	actualEvents := app.AllEmittedEvents()
	require.Equal(t, len(expectedEvents), len(actualEvents))

	if !assert.ObjectsAreEqualValues(expectedEvents, actualEvents) {
		expectedEventsMsg := getEventMessage(expectedEvents)
		eventsMsg := getEventMessage(actualEvents)

		assert.Fail(t, fmt.Sprintf("Not equal: \n"+
			"expected: %s\n"+
			"actual  : %s\n", expectedEventsMsg, eventsMsg), "msg2")
	}
}

func getEventMessage(events []event.Event) string {
	msg := "\n"
	for _, e := range events {
		msg += fmt.Sprintf("%-30s ", reflect.TypeOf(e).Name())
		b, _ := json.Marshal(e)
		msg += string(b) + "\n"
	}
	return msg
}
