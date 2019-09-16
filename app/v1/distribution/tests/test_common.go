package tests

import (
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v1/bank"
	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/app/v1/stake"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/store"
	sdk "github.com/irisnet/irishub/types"

	distr "github.com/irisnet/irishub/app/v1/distribution/keeper"
	"github.com/irisnet/irishub/app/v1/distribution/types"
)

var (
	delPk1   = ed25519.GenPrivKey().PubKey()
	delPk2   = ed25519.GenPrivKey().PubKey()
	delPk3   = ed25519.GenPrivKey().PubKey()
	delAddr1 = sdk.AccAddress(delPk1.Address())
	delAddr2 = sdk.AccAddress(delPk2.Address())
	delAddr3 = sdk.AccAddress(delPk3.Address())

	valOpPk1    = ed25519.GenPrivKey().PubKey()
	valOpPk2    = ed25519.GenPrivKey().PubKey()
	valOpPk3    = ed25519.GenPrivKey().PubKey()
	valOpAddr1  = sdk.ValAddress(valOpPk1.Address())
	valOpAddr2  = sdk.ValAddress(valOpPk2.Address())
	valOpAddr3  = sdk.ValAddress(valOpPk3.Address())
	valAccAddr1 = sdk.AccAddress(valOpPk1.Address()) // generate acc addresses for these validator keys too
	valAccAddr2 = sdk.AccAddress(valOpPk2.Address())
	valAccAddr3 = sdk.AccAddress(valOpPk3.Address())

	valConsPk1   = ed25519.GenPrivKey().PubKey()
	valConsPk2   = ed25519.GenPrivKey().PubKey()
	valConsPk3   = ed25519.GenPrivKey().PubKey()
	valConsAddr1 = sdk.ConsAddress(valConsPk1.Address())
	valConsAddr2 = sdk.ConsAddress(valConsPk2.Address())
	valConsAddr3 = sdk.ConsAddress(valConsPk3.Address())

	addrs = []sdk.AccAddress{
		delAddr1, delAddr2, delAddr3,
		valAccAddr1, valAccAddr2, valAccAddr3,
	}

	emptyDelAddr sdk.AccAddress
	emptyValAddr sdk.ValAddress
	emptyPubkey  crypto.PubKey
)

// create a codec used only for testing
func MakeTestCodec() *codec.Codec {
	var cdc = codec.New()
	bank.RegisterCodec(cdc)
	stake.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	types.RegisterCodec(cdc) // distr
	return cdc
}

// test input with default values
func CreateTestInputDefault(t *testing.T, isCheckTx bool, initCoins sdk.Int) (
	sdk.Context, auth.AccountKeeper, distr.Keeper, stake.Keeper, DummyFeeCollectionKeeper) {

	communityTax := sdk.NewDecWithPrec(2, 2)
	return CreateTestInputAdvanced(t, isCheckTx, initCoins, communityTax)
}

// hogpodge of all sorts of input required for testing
func CreateTestInputAdvanced(t *testing.T, isCheckTx bool, initCoins sdk.Int,
	communityTax sdk.Dec) (
	sdk.Context, auth.AccountKeeper, distr.Keeper, stake.Keeper, DummyFeeCollectionKeeper) {

	keyDistr := sdk.NewKVStoreKey("distr")
	keyStake := sdk.NewKVStoreKey("stake")
	tkeyStake := sdk.NewTransientStoreKey("transient_stake")
	keyAcc := sdk.NewKVStoreKey("acc")
	keyFeeCollection := sdk.NewKVStoreKey("fee")
	keyParams := sdk.NewKVStoreKey("params")
	tkeyParams := sdk.NewTransientStoreKey("transient_params")

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)

	ms.MountStoreWithDB(keyDistr, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyStake, sdk.StoreTypeTransient, nil)
	ms.MountStoreWithDB(keyStake, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyFeeCollection, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)

	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	cdc := MakeTestCodec()
	pk := params.NewKeeper(cdc, keyParams, tkeyParams)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "foochainid"}, isCheckTx, log.NewNopLogger())
	accountKeeper := auth.NewAccountKeeper(cdc, keyAcc, auth.ProtoBaseAccount)
	ck := bank.NewBaseKeeper(cdc, accountKeeper)
	sk := stake.NewKeeper(cdc, keyStake, tkeyStake, ck, pk.Subspace(stake.DefaultParamspace), stake.DefaultCodespace, stake.NopMetrics())
	sk.SetPool(ctx, stake.Pool{BondedPool: stake.InitialBondedPool()})
	sk.SetParams(ctx, stake.DefaultParams())

	// fill all the addresses with some coins, set the loose pool tokens simultaneously
	for _, addr := range addrs {
		ck.AddCoins(ctx, addr, sdk.Coins{
			{sk.BondDenom(), initCoins},
		})
		ck.IncreaseLoosenToken(ctx, sdk.Coins{
			{sk.BondDenom(), initCoins},
		})
	}

	fck := DummyFeeCollectionKeeper{}
	keeper := distr.NewKeeper(cdc, keyDistr, pk.Subspace(distr.DefaultParamspace), ck, sk, fck, types.DefaultCodespace, distr.NopMetrics())

	// set the distribution hooks on staking
	sk.SetHooks(keeper.Hooks())

	// set genesis items required for distribution
	keeper.SetFeePool(ctx, types.InitialFeePool())
	params := types.Params{
		CommunityTax:        communityTax,
		BaseProposerReward:  sdk.NewDecWithPrec(1, 2),
		BonusProposerReward: sdk.NewDecWithPrec(4, 2),
	}
	keeper.SetParams(ctx, params)
	return ctx, accountKeeper, keeper, sk, fck
}

//__________________________________________________________________________________
// fee collection keeper used only for testing
type DummyFeeCollectionKeeper struct{}

var heldFees sdk.Coins
var _ types.FeeKeeper = DummyFeeCollectionKeeper{}

// nolint
func (fck DummyFeeCollectionKeeper) GetCollectedFees(_ sdk.Context) sdk.Coins {
	return heldFees
}
func (fck DummyFeeCollectionKeeper) SetCollectedFees(in sdk.Coins) {
	heldFees = in
}
func (fck DummyFeeCollectionKeeper) ClearCollectedFees(_ sdk.Context) {
	heldFees = sdk.Coins{}
}
