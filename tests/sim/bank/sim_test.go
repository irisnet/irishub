package bank

import (
	"encoding/json"
	"math/rand"
	"testing"

	"github.com/irisnet/irishub/mock"
	"github.com/irisnet/irishub/mock/simulation"
	"github.com/irisnet/irishub/modules/bank"
	stakeTypes "github.com/irisnet/irishub/modules/stake/types"
	sdk "github.com/irisnet/irishub/types"
)

func TestBankWithRandomMessages(t *testing.T) {
	mapp := mock.NewApp()

	bank.RegisterCodec(mapp.Cdc)
	mapper := mapp.AccountKeeper
	bankKeeper := mapp.BankKeeper

	mapp.Router().AddRoute("bank", []*sdk.KVStoreKey{mapp.KeyAccount}, bank.NewHandler(bankKeeper))

	err := mapp.CompleteSetup()
	if err != nil {
		panic(err)
	}

	appStateFn := func(r *rand.Rand, accs []simulation.Account) json.RawMessage {
		simulation.RandomSetGenesis(r, mapp, accs, []string{stakeTypes.StakeDenom})
		return json.RawMessage("{}")
	}

	simulation.Simulate(
		t, mapp.BaseApp, appStateFn,
		[]simulation.WeightedOperation{
			{1, SingleInputSendMsg(mapper, bankKeeper)},
		},
		[]simulation.RandSetup{},
		[]simulation.Invariant{
			bank.NonnegativeBalanceInvariant(mapper),
			bank.TotalCoinsInvariant(mapper, func() sdk.Coins { return mapp.TotalCoinsSupply }),
		},
		30, 60,
		false,
	)
}
