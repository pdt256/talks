package projection

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/icrowley/fake"

	"github.com/pdt256/talks/code/cqrs/go/bank"
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event"
)

type fakeStream struct {
	timestamp            int64
	maxTimeBetweenEvents int
	maxAccounts          int
}

func NewFakeStream(seed int64, maxAccounts int) *fakeStream {
	rand.Seed(seed)
	uuid.SetRand(rand.New(rand.NewSource(seed)))
	fake.Seed(seed)

	janFirst := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)

	return &fakeStream{
		timestamp:            janFirst.Unix(),
		maxTimeBetweenEvents: int((time.Minute * 30).Seconds()),
		maxAccounts:          maxAccounts,
	}
}

func (f *fakeStream) Load() <-chan event.Event {
	ch := make(chan event.Event)

	go func() {
		var accountGenerators []*accountGenerator
		for i := 0; i < f.maxAccounts; i++ {
			accountGenerators = append(accountGenerators, newAccountGenerator())
		}

		for len(accountGenerators) > 0 {
			i := rand.Intn(len(accountGenerators))

			if !accountGenerators[i].HasNext() {
				accountGenerators = removeAccountGenerator(accountGenerators, i)
				continue
			}

			e := accountGenerators[i].Next(f.nextEventBase())

			if e != nil {
				ch <- e
			}

		}

		close(ch)
	}()

	return ch
}

func removeAccountGenerator(s []*accountGenerator, i int) []*accountGenerator {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func percentOfTime25() bool {
	return rand.Intn(4) == 0
}

func (f *fakeStream) tick() int64 {
	f.timestamp += int64(rand.Intn(f.maxTimeBetweenEvents))
	return f.timestamp
}

func (f *fakeStream) nextEventBase() event.Base {
	return event.Base{Timestamp: f.tick()}
}

type accountGenerator struct {
	isOpen                 bool
	transactionsLeft       int
	maxDepositAmount       int
	maxWithdrawAmount      int
	balance                int
	isDone                 bool
	closeAccountOnNextTurn bool
	accountId              string
}

func newAccountGenerator() *accountGenerator {
	maxTransactionsPerAccount := 100

	return &accountGenerator{
		isOpen:            false,
		transactionsLeft:  rand.Intn(maxTransactionsPerAccount),
		maxDepositAmount:  rand.Intn(15000) + 1,
		maxWithdrawAmount: rand.Intn(1000) + 1,
		accountId:         uuid.New().String(),
	}
}

func (a *accountGenerator) Next(base event.Base) event.Event {

	if !a.isOpen {
		if percentOfTime25() {
			return bank.DepositFailed{
				Base:      base,
				AccountId: a.accountId,
				Amount:    0,
			}
		}

		a.isOpen = true
		return bank.AccountWasOpened{
			Base:      base,
			AccountId: a.accountId,
			FirstName: fake.FirstName(),
			LastName:  fake.LastName(),
		}
	}

	if a.closeAccountOnNextTurn && a.balance == 0 {
		a.closeAccountOnNextTurn = false
		return bank.AccountWasClosed{
			Base:      base,
			AccountId: a.accountId,
		}
	}

	if a.transactionsLeft > 0 {
		a.transactionsLeft--
		switch rand.Intn(3) {
		case 0:
			fallthrough
		case 1:
			depositAmount := rand.Intn(a.maxDepositAmount)
			a.balance += depositAmount
			return bank.MoneyWasDeposited{
				Base:       base,
				AccountId:  a.accountId,
				Amount:     depositAmount,
				NewBalance: a.balance,
			}
		case 2:
			withdrawAmount := rand.Intn(a.maxWithdrawAmount)

			if a.balance < withdrawAmount {
				if percentOfTime25() {
					return bank.WithdrawDenied{
						Base:           base,
						AccountId:      a.accountId,
						Amount:         withdrawAmount,
						CurrentBalance: a.balance,
					}
				}

				return nil
			}

			if (a.balance - withdrawAmount) < 500 {
				return bank.FailedToCloseAccountWithBalance{
					Base:      base,
					AccountId: a.accountId,
					Balance:   a.balance,
				}
			}

			if (a.balance - withdrawAmount) < 2000 {
				a.closeAccountOnNextTurn = true
				withdrawAmount = a.balance
				a.balance = 0
				return bank.MoneyWasWithdrawn{
					Base:       base,
					AccountId:  a.accountId,
					Amount:     withdrawAmount,
					NewBalance: a.balance,
				}
			}

			a.balance -= withdrawAmount
			return bank.MoneyWasWithdrawn{
				Base:       base,
				AccountId:  a.accountId,
				Amount:     withdrawAmount,
				NewBalance: a.balance,
			}
		}
	}

	a.isDone = true
	return nil
}

func (a *accountGenerator) HasNext() bool {
	return !a.isDone
}
