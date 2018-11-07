package simulation

import (
	"encoding/json"
	"math/rand"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/irisnet/irishub/simulation/mock"
	"github.com/irisnet/irishub/simulation/mock/simulation"
	abci "github.com/tendermint/tendermint/abci/types"
)

// TestStakeWithRandomMessages
func TestStakeWithRandomMessages(t *testing.T) {
	mapp := mock.NewApp()

	bank.RegisterCodec(mapp.Cdc)
	mapper := mapp.AccountKeeper
	bankKeeper := mapp.BankKeeper

	feeKey := mapp.KeyFeeCollection
	stakeKey := mapp.KeyStake
	stakeTKey := mapp.TkeyStake
	paramsKey := mapp.KeyParams
	paramsTKey := mapp.TkeyParams
	distrKey := sdk.NewKVStoreKey("distr")

	feeCollectionKeeper := auth.NewFeeCollectionKeeper(mapp.Cdc, feeKey)
	paramstore := params.NewKeeper(mapp.Cdc, paramsKey, paramsTKey)
	stakeKeeper := stake.NewKeeper(mapp.Cdc, stakeKey, stakeTKey, bankKeeper, paramstore.Subspace(stake.DefaultParamspace), stake.DefaultCodespace)
	distrKeeper := distribution.NewKeeper(mapp.Cdc, distrKey, paramstore.Subspace(distribution.DefaultParamspace), bankKeeper, stakeKeeper, feeCollectionKeeper, distribution.DefaultCodespace)
	mapp.Router().AddRoute("stake", []*sdk.KVStoreKey{stakeKey, mapp.KeyAccount, distrKey}, stake.NewHandler(stakeKeeper))
	mapp.SetEndBlocker(func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		validatorUpdates := stake.EndBlocker(ctx, stakeKeeper)
		return abci.ResponseEndBlock{
			ValidatorUpdates: validatorUpdates,
		}
	})

	err := mapp.CompleteSetup(stakeKey, stakeTKey, paramsKey, paramsTKey, distrKey)
	if err != nil {
		panic(err)
	}

	appStateFn := func(r *rand.Rand, accs []simulation.Account) json.RawMessage {
		simulation.RandomSetGenesis(r, mapp, accs, []string{"iris-atto"})
		return json.RawMessage("{}")
	}

	setup := func(r *rand.Rand, accs []simulation.Account) {
		ctx := mapp.NewContext(false, abci.Header{})
		distribution.InitGenesis(ctx, distrKeeper, distribution.DefaultGenesisState())
	}

	simulation.Simulate(
		t, mapp.BaseApp, appStateFn,
		[]simulation.WeightedOperation{
			{10, SimulateMsgCreateValidator(mapper, stakeKeeper)},
			{5, SimulateMsgEditValidator(stakeKeeper)},
			{15, SimulateMsgDelegate(mapper, stakeKeeper)},
			{10, SimulateMsgBeginUnbonding(mapper, stakeKeeper)},
			{10, SimulateMsgBeginRedelegate(mapper, stakeKeeper)},
		}, []simulation.RandSetup{
			Setup(mapp, stakeKeeper),
			setup,
		}, []simulation.Invariant{
			//AllInvariants(bankKeeper, stakeKeeper, feeCollectionKeeper, distrKeeper, mapp.AccountKeeper),
		}, 10, 100,
		false,
	)
}
