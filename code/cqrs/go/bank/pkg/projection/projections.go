package projection

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/pdt256/talks/code/cqrs/go/bank"
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event"
)

type Count struct {
	EventCount    int
	DepositCount  int
	WithdrawCount int
}

func NewCount() *Count {
	return &Count{}
}

func (c *Count) Accept(e event.Event) {
	c.EventCount++

	switch e.(type) {
	case *bank.MoneyWasDeposited:
		c.DepositCount++

	case *bank.MoneyWasWithdrawn:
		c.WithdrawCount++
	}
}

type BankBalance struct {
	TotalBalance        int
	TotalBalanceByMonth map[string]int
}

func NewBankBalance() *BankBalance {
	return &BankBalance{
		TotalBalanceByMonth: make(map[string]int),
	}
}

func (b *BankBalance) Accept(e event.Event) {
	monthString := GetMonth(e.EventTimestamp())

	if _, ok := b.TotalBalanceByMonth[monthString]; !ok {
		b.TotalBalanceByMonth[monthString] = b.TotalBalance
	}

	switch e := e.(type) {
	case *bank.MoneyWasDeposited:
		b.TotalBalance += e.Amount
		b.TotalBalanceByMonth[monthString] += e.Amount

	case *bank.MoneyWasWithdrawn:
		b.TotalBalance -= e.Amount
		b.TotalBalanceByMonth[monthString] -= e.Amount
	}
}

type AccountBalance struct {
	TotalBalanceByAccount         map[string]int
	AccountNames                  map[string]string
	TotalBalanceByMonthAndAccount map[string]map[string]int
}

func NewAccountBalance() *AccountBalance {
	return &AccountBalance{
		TotalBalanceByAccount:         make(map[string]int),
		TotalBalanceByMonthAndAccount: make(map[string]map[string]int),
		AccountNames:                  make(map[string]string),
	}
}

func (a *AccountBalance) Accept(e event.Event) {
	monthString := GetMonth(e.EventTimestamp())

	if _, ok := a.TotalBalanceByMonthAndAccount[monthString]; !ok {
		a.TotalBalanceByMonthAndAccount[monthString] = make(map[string]int)

		for accountId, balance := range a.TotalBalanceByAccount {
			a.TotalBalanceByMonthAndAccount[monthString][accountId] = balance
		}
	}

	switch e := e.(type) {
	case *bank.AccountWasOpened:
		a.AccountNames[e.AccountId] = e.FirstName + " " + e.LastName

	case *bank.MoneyWasDeposited:
		a.TotalBalanceByAccount[e.AccountId] += e.Amount
		a.TotalBalanceByMonthAndAccount[monthString][e.AccountId] += e.Amount

	case *bank.MoneyWasWithdrawn:
		a.TotalBalanceByAccount[e.AccountId] -= e.Amount
		a.TotalBalanceByMonthAndAccount[monthString][e.AccountId] -= e.Amount
	}
}

func GetMonth(timestamp int64) string {
	return time.Unix(timestamp, 0).UTC().Format("2006-01")
}

func (a AccountAndBalance) String() string {
	return fmt.Sprintf("%s - %s: %d", a.AccountId, a.Name, a.Balance)
}

func (a *AccountBalance) GetTop5AccountsByBalance() []*AccountAndBalance {
	return GetTop5AccountsByBalance(a.TotalBalanceByAccount, a.AccountNames)
}

func (a *AccountBalance) GetTop5AccountsByBalanceForMonth(month string) []*AccountAndBalance {
	return GetTop5AccountsByBalance(a.TotalBalanceByMonthAndAccount[month], a.AccountNames)
}

func (a *AccountBalance) PrintTop5AccountsByBalanceByMonth() {
	months := make([]string, 0, len(a.TotalBalanceByMonthAndAccount))
	for k := range a.TotalBalanceByMonthAndAccount {
		months = append(months, k)
	}

	sort.Strings(months)

	for _, month := range months {
		fmt.Printf("[%s]\n", month)
		accountBalances := a.TotalBalanceByMonthAndAccount[month]

		for _, accountAndBalance := range GetTop5AccountsByBalance(accountBalances, a.AccountNames) {
			fmt.Printf("  %s - %s: %d\n",
				accountAndBalance.AccountId,
				accountAndBalance.Name,
				accountAndBalance.Balance,
			)
		}
	}
}

type AccountAndBalance struct {
	Balance   int
	AccountId string
	Name      string
}

func GetTop5AccountsByBalance(balances map[string]int, names map[string]string) []*AccountAndBalance {
	var accountBalances []*AccountAndBalance
	for accountId, balance := range balances {
		accountBalances = append(accountBalances, &AccountAndBalance{
			Balance:   balance,
			AccountId: accountId,
		})
	}

	sort.Slice(accountBalances[:], func(i, j int) bool {
		return accountBalances[i].Balance > accountBalances[j].Balance
	})

	for i := range accountBalances {
		accountId := accountBalances[i].AccountId
		accountBalances[i].Name = names[accountId]
	}

	five := int(math.Min(5, float64(len(accountBalances))))

	return accountBalances[:five]
}
