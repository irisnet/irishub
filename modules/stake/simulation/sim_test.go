package simulation

import (
	"encoding/json"
	"math/rand"
	"testing"

	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/auth"
	"github.com/irisnet/irishub/modules/bank"
	"github.com/irisnet/irishub/modules/distribution"
	"github.com/irisnet/irishub/modules/params"
	"github.com/irisnet/irishub/modules/stake"
	"github.com/irisnet/irishub/modules/mock"
	"github.com/irisnet/irishub/modules/mock/simulation"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/irisnet/irishub/modules/stake/types"
)

// TestStakeWithRandomMessages
func TestStakeWithRandomMessages(t *testing.T) {
	mapp := mock.NewApp()

	bank.RegisterCodec(mapp.Cdc)
	stake.RegisterCodec(mapp.Cdc)

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

	err := mapp.CompleteSetup(distrKey)
	if err != nil {
		panic(err)
	}

	appStateFn := func(r *rand.Rand, accs []simulation.Account) json.RawMessage {
		simulation.RandomSetGenesis(r, mapp, accs, []string{types.StakeDenom})
		return json.RawMessage("{}")
	}

	GenesisSetUp := func(r *rand.Rand, accs []simulation.Account) {
		ctx := mapp.NewContext(false, abci.Header{})
		distribution.InitGenesis(ctx, distrKeeper, distribution.DefaultGenesisState())

		// init stake genesis
		var (
			validators  []stake.Validator
			delegations []stake.Delegation
		)
		stakeGenesis := stake.DefaultGenesisState()

		// XXX Try different numbers of initially bonded validators
		numInitiallyBonded := int64(4)
		valAddrs := make([]sdk.ValAddress, numInitiallyBonded)
		decAmt := sdk.NewDecFromInt(sdk.NewIntWithDecimal(1, 2))
		for i := 0; i < int(numInitiallyBonded); i++ {
			valAddr := sdk.ValAddress(accs[i].Address)
			valAddrs[i] = valAddr

			validator := stake.NewValidator(valAddr, accs[i].PubKey, stake.Description{})
			validator.Tokens = decAmt
			validator.DelegatorShares = decAmt
			delegation := stake.Delegation{accs[i].Address, valAddr, decAmt, 0}
			validators = append(validators, validator)
			delegations = append(delegations, delegation)
		}
		stakeGenesis.Pool.LooseTokens = sdk.NewDecFromInt(sdk.NewIntWithDecimal(1, 10))
		stakeGenesis.Validators = validators
		stakeGenesis.Bonds = delegations

		stake.InitGenesis(ctx, stakeKeeper, stakeGenesis)
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
			//Setup(mapp, stakeKeeper),
			GenesisSetUp,
		}, []simulation.Invariant{}, 10, 100,
		false,
	)
}
