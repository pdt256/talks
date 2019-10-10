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
	assert.Equal(t, "e5e6e7e8-e9ea-4bec-adee-eff0f1f2f3f4 - Peter Price: 9190207", top5Balance[0].String())
	assert.Equal(t, "45464748-494a-4b4c-8d4e-4f5051525354 - Edward White: 8393523", top5Balance[1].String())
	assert.Equal(t, "75767778-797a-4b7c-bd7e-7f8081828384 - Sarah Gonzales: 8152814", top5Balance[2].String())
	assert.Equal(t, "25262728-292a-4b2c-ad2e-2f3031323334 - Diana Armstrong: 8138200", top5Balance[3].String())
	assert.Equal(t, "f5f6f7f8-f9fa-4bfc-bdfe-ff0001020304 - Jack Owens: 7974992", top5Balance[4].String())

	// Find the top 5 accounts with the highest balance by month.
	top5ByMonth := accountBalance.GetTop5AccountsByBalanceForMonth("2019-01")
	assert.Equal(t, "25262728-292a-4b2c-ad2e-2f3031323334 - Diana Armstrong: 337389", top5ByMonth[0].String())
	assert.Equal(t, "45464748-494a-4b4c-8d4e-4f5051525354 - Edward White: 329289", top5ByMonth[1].String())
	assert.Equal(t, "e5e6e7e8-e9ea-4bec-adee-eff0f1f2f3f4 - Peter Price: 324844", top5ByMonth[2].String())
	assert.Equal(t, "55565758-595a-4b5c-9d5e-5f6061626364 - Gloria Stanley: 293658", top5ByMonth[3].String())
	assert.Equal(t, "b5b6b7b8-b9ba-4bbc-bdbe-bfc0c1c2c3c4 - Ashley Olson: 281789", top5ByMonth[4].String())
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
