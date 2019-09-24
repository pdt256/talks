package bank

// Commands
type OpenAccount struct {
	AccountId string
}
type DepositMoney struct {
	AccountId string
	Amount    int
}
type WithdrawMoney struct {
	AccountId string
	Amount    int
}
type CloseAccount struct {
	AccountId string
}

// Events
type AccountWasOpened struct {
	AccountId string
}
type MoneyWasDeposited struct {
	AccountId  string
	Amount     int
	NewBalance int
}
type DepositFailed struct {
	AccountId string
	Amount    int
}
type MoneyWasWithdrawn struct {
	AccountId  string
	Amount     int
	NewBalance int
}
type WithdrawDenied struct {
	AccountId      string
	Amount         int
	CurrentBalance int
}
type AccountWasClosed struct {
	AccountId string
}
type FailedToCloseAccountWithBalance struct {
	AccountId string
	Balance   int
}
