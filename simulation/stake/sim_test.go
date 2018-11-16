package simulation

import (
	"encoding/json"
	"math/rand"
	"testing"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/irisnet/irishub/simulation/mock"
	"github.com/irisnet/irishub/simulation/mock/simulation"
	"github.com/cosmos/cosmos-sdk/x/stake"
)

// TestStakeWithRandomMessages
func TestStakeWithRandomMessages(t *testing.T) {
	mapp := mock.NewApp()

	bank.RegisterCodec(mapp.Cdc)
	mapper := mapp.AccountKeeper
	coinKeeper := bank.NewKeeper(mapper)
	stakeKey := sdk.NewKVStoreKey("stake")
	stakeKeeper := stake.NewKeeper(mapp.Cdc, stakeKey, coinKeeper, stake.DefaultCodespace)
	mapp.Router().AddRoute("stake", []*sdk.KVStoreKey{mapp.KeyStake, mapp.KeyAccount}, stake.NewHandler(stakeKeeper))
	mapp.SetEndBlocker(func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		validatorUpdates := stake.EndBlocker(ctx, stakeKeeper)
		return abci.ResponseEndBlock{
			ValidatorUpdates: validatorUpdates,
		}
	})

	err := mapp.CompleteSetup([]*sdk.KVStoreKey{stakeKey})
	if err != nil {
		panic(err)
	}

	appStateFn := func(r *rand.Rand, keys []crypto.PrivKey, accs []sdk.AccAddress) json.RawMessage {
		mock.RandomSetGenesis(r, mapp, accs, []string{"iris"})
		return json.RawMessage("{}")
	}

	simulation.Simulate(
		t, mapp.BaseApp, appStateFn,
		[]simulation.TestAndRunTx{
			SimulateMsgCreateValidator(mapper, stakeKeeper),
			SimulateMsgEditValidator(stakeKeeper),
			SimulateMsgDelegate(mapper, stakeKeeper),
			SimulateMsgBeginUnbonding(mapper, stakeKeeper),
			SimulateMsgCompleteUnbonding(stakeKeeper),
			SimulateMsgBeginRedelegate(mapper, stakeKeeper),
			SimulateMsgCompleteRedelegate(stakeKeeper),
		}, []simulation.RandSetup{
			Setup(mapp, stakeKeeper),
		}, []simulation.Invariant{
			AllInvariants(coinKeeper, stakeKeeper, mapp.AccountKeeper),
		}, 10, 100,
		false,
	)
}
