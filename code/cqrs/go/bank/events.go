package bank

import "github.com/pdt256/talks/code/cqrs/go/bank/pkg/event"

//go:generate go run gen/eventbinder/main.go -file events.go

type AccountWasOpened struct {
	event.Base
	AccountId string
	FirstName string
	LastName  string
}
type MoneyWasDeposited struct {
	event.Base
	AccountId  string
	Amount     int
	NewBalance int
}
type DepositFailed struct {
	event.Base
	AccountId string
	Amount    int
}
type MoneyWasWithdrawn struct {
	event.Base
	AccountId  string
	Amount     int
	NewBalance int
}
type WithdrawDenied struct {
	event.Base
	AccountId      string
	Amount         int
	CurrentBalance int
}
type AccountWasClosed struct {
	event.Base
	AccountId string
}
type FailedToCloseAccountWithBalance struct {
	event.Base
	AccountId string
	Balance   int
}
