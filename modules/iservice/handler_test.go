package iservice

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk"
	"github.com/cosmos/cosmos-sdk/state"
	"github.com/stretchr/testify/assert"
)

//______________________________________________________________________

func initAccounts(n int, amount int64) ([]sdk.Actor, map[string]int64) {
	accStore := map[string]int64{}
	senders := newActors(n)
	for _, sender := range senders {
		accStore[string(sender.Address)] = amount
	}
	return senders, accStore
}

func newDeliver(sender sdk.Actor, accStore map[string]int64) deliver {
	store := state.NewMemKVStore()
	return deliver{
		store:  store,
		sender: sender,
		params: loadParams(store),
	}
}

// newTxDefineService - new TxDefineService
func newTxDefineService(name string, description string) TxDefineService {
	return TxDefineService{
		Name:        name,
		Description: description,
	}
}

func TestTxDefineService(t *testing.T) {
	require := assert.New(t)
	accounts, accStore := initAccounts(3, 1000)
	sender := accounts[0]
	deliverer := newDeliver(sender, accStore)

	//first make a candidate
	txDefineService := newTxDefineService("testname", "testdesc")
	deliverer.sender = sender
	got := deliverer.defineService(txDefineService)
	require.NoError(got, "expected tx to be ok, got %v", got)
}
