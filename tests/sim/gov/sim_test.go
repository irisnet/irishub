package gov

import (
	"encoding/json"
	"math/rand"
	"testing"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/irisnet/irishub/mock"
	"github.com/irisnet/irishub/mock/simulation"
	"github.com/irisnet/irishub/modules/bank"
	distr "github.com/irisnet/irishub/modules/distribution"
	"github.com/irisnet/irishub/modules/gov"
	"github.com/irisnet/irishub/modules/guardian"
	"github.com/irisnet/irishub/modules/stake"
	sdk "github.com/irisnet/irishub/types"
)

// TestGovWithRandomMessages
func TestGovWithRandomMessages(t *testing.T) {
	mapp := mock.NewApp()

	bank.RegisterCodec(mapp.Cdc)
	gov.RegisterCodec(mapp.Cdc)

	bankKeeper := mapp.BankKeeper
	stakeKey := mapp.KeyStake
	stakeTKey := mapp.TkeyStake
	paramKey := mapp.KeyParams
	govKey := sdk.NewKVStoreKey("gov")
	distrKey := sdk.NewKVStoreKey("distr")
	guardianKey := sdk.NewKVStoreKey("guardian")

	paramKeeper := mapp.ParamsKeeper
	stakeKeeper := stake.NewKeeper(
		mapp.Cdc, stakeKey,
		stakeTKey, bankKeeper,
		paramKeeper.Subspace(stake.DefaultParamspace),
		stake.DefaultCodespace,
		stake.NopMetrics(),
	)
	distrKeeper := distr.NewKeeper(
		mapp.Cdc,
		distrKey,
		mapp.ParamsKeeper.Subspace(distr.DefaultParamspace),
		mapp.BankKeeper, &stakeKeeper, mapp.FeeKeeper,
		distr.DefaultCodespace,
		distr.NopMetrics(),
	)
	guardianKeeper := guardian.NewKeeper(
		mapp.Cdc,
		guardianKey,
		guardian.DefaultCodespace,
	)
	govKeeper := gov.NewKeeper(
		govKey,
		mapp.Cdc,
		mapp.ParamsKeeper.Subspace(gov.DefaultParamSpace),
		paramKeeper,
		sdk.NewProtocolKeeper(mapp.KeyMain),
		bankKeeper,
		distrKeeper,
		guardianKeeper,
		stakeKeeper,
		gov.DefaultCodespace,
		gov.NopMetrics(),
	)

	mapp.Router().AddRoute("gov", []*sdk.KVStoreKey{govKey, mapp.KeyAccount, stakeKey, paramKey}, gov.NewHandler(govKeeper))
	mapp.SetEndBlocker(func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		gov.EndBlocker(ctx, govKeeper)
		return abci.ResponseEndBlock{}
	})

	err := mapp.CompleteSetup(govKey)
	if err != nil {
		panic(err)
	}

	appStateFn := func(r *rand.Rand, accs []simulation.Account) json.RawMessage {
		simulation.RandomSetGenesis(r, mapp, accs, []string{"stake"})
		return json.RawMessage("{}")
	}

	setup := func(r *rand.Rand, accs []simulation.Account) {
		ctx := mapp.NewContext(false, abci.Header{})
		stake.InitGenesis(ctx, stakeKeeper, stake.DefaultGenesisState())

		gov.InitGenesis(ctx, govKeeper, gov.DefaultGenesisState())
	}

	// Test with unscheduled votes
	simulation.Simulate(
		t, mapp.BaseApp, appStateFn,
		[]simulation.WeightedOperation{
			{2, SimulateMsgSubmitProposal(govKeeper, stakeKeeper)},
			{3, SimulateMsgDeposit(govKeeper, stakeKeeper)},
			{20, SimulateMsgVote(govKeeper, stakeKeeper)},
		}, []simulation.RandSetup{
			setup,
		}, []simulation.Invariant{
			//AllInvariants(),
		}, 10, 100,
		false,
	)

	// Test with scheduled votes
	simulation.Simulate(
		t, mapp.BaseApp, appStateFn,
		[]simulation.WeightedOperation{
			{10, SimulateSubmittingVotingAndSlashingForProposal(govKeeper, stakeKeeper)},
			{5, SimulateMsgDeposit(govKeeper, stakeKeeper)},
		}, []simulation.RandSetup{
			setup,
		}, []simulation.Invariant{
			gov.AllInvariants(),
		}, 10, 100,
		false,
	)
}
