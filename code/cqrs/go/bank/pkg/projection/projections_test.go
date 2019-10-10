package projection_test

import (
	"io"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pdt256/talks/code/cqrs/go/bank"
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event"
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event/provider/inmemorybus"
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/event/provider/jsoniostream"
	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/projection"
)

//go:generate go run ../../gen/bankevents/main.go -seed 100 -accounts 1000 -out events.json

func TestProjections(t *testing.T) {
	// Given
	file := GetEventsFile()
	defer file.Close()
	bus, jsonStream := GetStream()
	count := projection.NewCount()
	bankBalance := projection.NewBankBalance()
	accountBalance := projection.NewAccountBalance()
	bus.Subscribe(
		count,
		bankBalance,
		accountBalance,
	)

	// When
	for e := range jsonStream.Load(file) {
		bus.Publish(e)
	}

	// Then
	// Count the number of events
	assert.Equal(t, 50223, count.EventCount)

	// Count the number of deposits and withdrawals
	assert.Equal(t, 32778, count.DepositCount)
	assert.Equal(t, 15142, count.WithdrawCount)

	// Get the total balance for the bank across all accounts
	assert.Equal(t, 119333376, bankBalance.TotalBalance)

	// Get the total running balance for the bank for each month
	assert.Equal(t, 4309529, bankBalance.TotalBalanceByMonth["2019-01"])
	assert.Equal(t, 10524290, bankBalance.TotalBalanceByMonth["2019-02"])
	assert.Equal(t, 17730801, bankBalance.TotalBalanceByMonth["2019-03"])

	// Find the top 5 accounts with the highest balance. Include account name and id.
	top5Balance := accountBalance.GetTop5AccountsByBalance()
	assert.Equal(t, "bc8e7b50-ef7b-42df-85b3-516b5d7bd0c5 - Eugene Romero: 517267", top5Balance[0].String())
	assert.Equal(t, "ce8b8128-9818-4ee4-a940-10484038cb96 - Gloria Oliver: 495081", top5Balance[1].String())
	assert.Equal(t, "df39503b-4b8b-4659-a7c4-21f1dc27d770 - Virginia Ford: 486708", top5Balance[2].String())
	assert.Equal(t, "a8607124-a2bb-463b-b882-e53874f0cd70 - Wayne Fields: 470417", top5Balance[3].String())
	assert.Equal(t, "75176d7a-8465-44cd-bfc5-de8a817443e5 - Michelle Allen: 452795", top5Balance[4].String())

	// Find the top 5 accounts with the highest balance by month.
	top5ByMonth := accountBalance.GetTop5AccountsByBalanceForMonth("2019-12")
	assert.Equal(t, "8cd15986-a251-46da-a125-cd1d1a8e526e - Alan Howell: 340858", top5ByMonth[0].String())
	assert.Equal(t, "dcfd5d5e-1a62-41da-8c84-ff2ea75f77de - Lori Lawson: 286530", top5ByMonth[1].String())
	assert.Equal(t, "f72aa989-f423-4ac3-a469-72152390937f - Debra Hayes: 284295", top5ByMonth[2].String())
	assert.Equal(t, "112c7fa7-fc84-4c91-8d64-2dc65d59e050 - Paula Mason: 264958", top5ByMonth[3].String())
	assert.Equal(t, "14660e92-d709-4d1f-86b2-75044b3e2ce9 - Evelyn Fox: 260446", top5ByMonth[4].String())
	//accountBalance.PrintTop5AccountsByBalanceByMonth()
}

func GetStream() (event.Bus, event.Stream) {
	bus := inmemorybus.New()
	jsonStream := jsoniostream.New()
	bank.BindEvents(jsonStream)

	return bus, jsonStream
}

func GetEventsFile() io.ReadCloser {
	file, err := os.Open("events.json")
	if err != nil {
		log.Fatalf("unable to open json events file: %v", err)
	}
	return file
}
