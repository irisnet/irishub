package simulation

import (
	"encoding/json"
	"math/rand"
	"testing"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/irisnet/irishub/modules/gov"
	"github.com/irisnet/irishub/simulation/mock"
	"github.com/irisnet/irishub/simulation/mock/simulation"
	"github.com/cosmos/cosmos-sdk/x/stake"
)

// TestGovWithRandomMessages
func TestGovWithRandomMessages(t *testing.T) {
	mapp := mock.NewApp()

	bank.RegisterCodec(mapp.Cdc)
	gov.RegisterCodec(mapp.Cdc)
	mapper := mapp.AccountKeeper
	coinKeeper := bank.NewKeeper(mapper)
	stakeKey := sdk.NewKVStoreKey("stake")
	stakeKeeper := stake.NewKeeper(mapp.Cdc, stakeKey, coinKeeper, stake.DefaultCodespace)
	paramKey := sdk.NewKVStoreKey("params")
	govKey := sdk.NewKVStoreKey("gov")
	govKeeper := gov.NewKeeper(mapp.Cdc, govKey, coinKeeper, stakeKeeper, gov.DefaultCodespace)
	mapp.Router().AddRoute("gov", []*sdk.KVStoreKey{mapp.KeyGov, mapp.KeyAccount, mapp.KeyStake, mapp.KeyParams}, gov.NewHandler(govKeeper))
	mapp.SetEndBlocker(func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		gov.EndBlocker(ctx, govKeeper)
		return abci.ResponseEndBlock{}
	})

	err := mapp.CompleteSetup([]*sdk.KVStoreKey{stakeKey, paramKey, govKey})
	if err != nil {
		panic(err)
	}

	appStateFn := func(r *rand.Rand, keys []crypto.PrivKey, accs []sdk.AccAddress) json.RawMessage {
		mock.RandomSetGenesis(r, mapp, accs, []string{"iris"})
		return json.RawMessage("{}")
	}

	setup := func(r *rand.Rand, privKeys []crypto.PrivKey) {
		ctx := mapp.NewContext(false, abci.Header{})
		stake.InitGenesis(ctx, stakeKeeper, stake.DefaultGenesisState())
		gov.InitGenesis(ctx, govKeeper, gov.DefaultGenesisState())
	}

	simulation.Simulate(
		t, mapp.BaseApp, appStateFn,
		[]simulation.TestAndRunTx{
			SimulateMsgSubmitProposal(govKeeper, stakeKeeper),
			SimulateMsgDeposit(govKeeper, stakeKeeper),
			SimulateMsgVote(govKeeper, stakeKeeper),
		}, []simulation.RandSetup{
			setup,
		}, []simulation.Invariant{
			AllInvariants(),
		}, 10, 100,
		false,
	)
}
