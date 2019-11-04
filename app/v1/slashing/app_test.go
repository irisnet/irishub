package slashing

import (
	"testing"

	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v1/bank"
	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/app/v1/stake"
	stakeTypes "github.com/irisnet/irishub/app/v1/stake/types"
	"github.com/irisnet/irishub/server/mock"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

var (
	priv1 = ed25519.GenPrivKey()
	addr1 = sdk.AccAddress(priv1.PubKey().Address())
	coins = sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}
)

// initialize the mock application for this module
func getMockApp(t *testing.T) (*mock.App, stake.Keeper, Keeper) {
	mApp := mock.NewApp()

	RegisterCodec(mApp.Cdc)
	keySlashing := sdk.NewKVStoreKey("slashing")
	bankKeeper := bank.NewBaseKeeper(mApp.Cdc, mApp.AccountKeeper)

	paramsKeeper := params.NewKeeper(mApp.Cdc, mApp.KeyParams, mApp.TkeyParams)
	stakeKeeper := stake.NewKeeper(mApp.Cdc, mApp.KeyStake, mApp.TkeyStake, bankKeeper, paramsKeeper.Subspace(stake.DefaultParamspace), stake.DefaultCodespace, stake.NopMetrics())
	keeper := NewKeeper(mApp.Cdc, keySlashing, stakeKeeper, paramsKeeper.Subspace(DefaultParamspace), DefaultCodespace, NopMetrics())
	mApp.Router().AddRoute("stake", []*sdk.KVStoreKey{mApp.KeyStake, mApp.KeyAccount, mApp.KeyParams}, stake.NewHandler(stakeKeeper))
	mApp.Router().AddRoute("slashing", []*sdk.KVStoreKey{mApp.KeyStake, keySlashing, mApp.KeyParams}, NewHandler(keeper))

	mApp.SetEndBlocker(getEndBlocker(stakeKeeper))
	mApp.SetInitChainer(getInitChainer(mApp, stakeKeeper))

	require.NoError(t, mApp.CompleteSetup(keySlashing))

	return mApp, stakeKeeper, keeper
}

// stake endblocker
func getEndBlocker(keeper stake.Keeper) sdk.EndBlocker {
	return func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		validatorUpdates := stake.EndBlocker(ctx, keeper)
		return abci.ResponseEndBlock{
			ValidatorUpdates: validatorUpdates,
		}
	}
}

// overwrite the mock init chainer
func getInitChainer(mapp *mock.App, keeper stake.Keeper) sdk.InitChainer {
	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		mapp.InitChainer(ctx, req)
		stakeGenesis := stake.DefaultGenesisState()
		validators, err := stake.InitGenesis(ctx, keeper, stakeGenesis)
		if err != nil {
			panic(err)
		}

		return abci.ResponseInitChain{
			Validators: validators,
		}
	}
}

func checkValidator(t *testing.T, mapp *mock.App, keeper stake.Keeper,
	addr sdk.AccAddress, expFound bool) stake.Validator {
	ctxCheck := mapp.BaseApp.NewContext(true, abci.Header{})
	validator, found := keeper.GetValidator(ctxCheck, sdk.ValAddress(addr1))
	require.Equal(t, expFound, found)
	return validator
}

func checkValidatorSigningInfo(t *testing.T, mapp *mock.App, keeper Keeper,
	addr sdk.ConsAddress, expFound bool) ValidatorSigningInfo {
	ctxCheck := mapp.BaseApp.NewContext(true, abci.Header{})
	signingInfo, found := keeper.getValidatorSigningInfo(ctxCheck, addr)
	require.Equal(t, expFound, found)
	return signingInfo
}

func TestSlashingMsgs(t *testing.T) {
	mapp, stakeKeeper, keeper := getMockApp(t)

	genCoin := sdk.NewCoin(stakeTypes.StakeDenom, sdk.NewIntWithDecimal(42, 18))
	bondCoin := sdk.NewCoin(stakeTypes.StakeDenom, sdk.NewIntWithDecimal(10, 18))

	acc1 := &auth.BaseAccount{
		Address: addr1,
		Coins:   sdk.Coins{genCoin},
	}
	accs := []auth.Account{acc1}
	mock.SetGenesis(mapp, accs)

	description := stake.NewDescription("foo_moniker", "", "", "")
	commission := stake.NewCommissionMsg(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())

	createValidatorMsg := stake.NewMsgCreateValidator(
		sdk.ValAddress(addr1), priv1.PubKey(), bondCoin, description, commission,
	)
	mock.SignCheckDeliver(t, mapp.BaseApp, []sdk.Msg{createValidatorMsg}, []uint64{0}, []uint64{0}, true, true, priv1)
	mock.CheckBalance(t, mapp, addr1, sdk.Coins{genCoin.Sub(bondCoin)})
	mapp.BeginBlock(abci.RequestBeginBlock{})

	validator := checkValidator(t, mapp, stakeKeeper, addr1, true)
	require.Equal(t, sdk.ValAddress(addr1), validator.OperatorAddr)
	require.Equal(t, sdk.Bonded, validator.Status)
	require.True(sdk.DecEq(t, sdk.NewDecFromInt(sdk.NewIntWithDecimal(10, 18)), validator.BondedTokens()))
	unjailMsg := MsgUnjail{ValidatorAddr: sdk.ValAddress(validator.ConsPubKey.Address())}

	// no signing info yet
	checkValidatorSigningInfo(t, mapp, keeper, sdk.ConsAddress(addr1), false)

	// unjail should fail with unknown validator
	res := mock.SignCheckDeliver(t, mapp.BaseApp, []sdk.Msg{unjailMsg}, []uint64{0}, []uint64{1}, false, false, priv1)
	require.EqualValues(t, CodeValidatorNotJailed, res.Code)
	require.EqualValues(t, DefaultCodespace, res.Codespace)
}
