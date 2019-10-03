package bank

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
