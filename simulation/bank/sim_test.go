package simulation

import (
	"encoding/json"
	"math/rand"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/bank"
	"github.com/irisnet/irishub/simulation/mock"
	"github.com/irisnet/irishub/simulation/mock/simulation"
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
		simulation.RandomSetGenesis(r, mapp, accs, []string{"iris-atto"})
		return json.RawMessage("{}")
	}

	simulation.Simulate(
		t, mapp.BaseApp, appStateFn,
		[]simulation.WeightedOperation{
			{1, SingleInputSendMsg(mapper, bankKeeper)},
		},
		[]simulation.RandSetup{},
		[]simulation.Invariant{
			NonnegativeBalanceInvariant(mapper),
			TotalCoinsInvariant(mapper, func() sdk.Coins { return mapp.TotalCoinsSupply }),
		},
		30, 60,
		false,
	)
}
