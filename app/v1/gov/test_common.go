package gov

import (
	"bytes"
	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v1/asset"
	"github.com/stretchr/testify/require"
	"log"
	"sort"
	"testing"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"

	"fmt"
	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v1/bank"
	"github.com/irisnet/irishub/app/v1/distribution"
	"github.com/irisnet/irishub/app/v1/mock"
	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/app/v1/stake"
	"github.com/irisnet/irishub/modules/guardian"
	sdk "github.com/irisnet/irishub/types"
)

// initialize the mock application for this module
func getMockApp(t *testing.T, numGenAccs int) (*mock.App, Keeper, stake.Keeper, []sdk.AccAddress, []crypto.PubKey, []crypto.PrivKey) {
	mapp := mock.NewApp()

	stake.RegisterCodec(mapp.Cdc)
	RegisterCodec(mapp.Cdc)

	keyGov := sdk.NewKVStoreKey("gov")
	keyDistr := sdk.NewKVStoreKey("distr")

	paramsKeeper := params.NewKeeper(
		mapp.Cdc,
		sdk.NewKVStoreKey("params"),
		sdk.NewTransientStoreKey("transient_params"),
	)
	feeKeeper := auth.NewFeeKeeper(
		mapp.Cdc,
		sdk.NewKVStoreKey("fee"),
		paramsKeeper.Subspace(auth.DefaultParamSpace),
	)

	ck := bank.NewBaseKeeper(mapp.Cdc, mapp.AccountKeeper)
	sk := stake.NewKeeper(
		mapp.Cdc,
		mapp.KeyStake, mapp.TkeyStake,
		mapp.BankKeeper, mapp.ParamsKeeper.Subspace(stake.DefaultParamspace),
		stake.DefaultCodespace,
		stake.NopMetrics())
	dk := distribution.NewKeeper(mapp.Cdc, keyDistr, paramsKeeper.Subspace(distribution.DefaultParamspace), ck, sk, feeKeeper, DefaultCodespace, distribution.NopMetrics())
	guardianKeeper := guardian.NewKeeper(mapp.Cdc, sdk.NewKVStoreKey("guardian"), guardian.DefaultCodespace)
	ak := asset.NewKeeper(mapp.Cdc, protocol.KeyAsset, ck, asset.DefaultCodespace, paramsKeeper.Subspace(asset.DefaultParamSpace))

	gk := NewKeeper(keyGov, mapp.Cdc, paramsKeeper.Subspace(DefaultParamSpace), paramsKeeper, sdk.NewProtocolKeeper(sdk.NewKVStoreKey("main")), ck, dk, guardianKeeper, sk, DefaultCodespace, NopMetrics(), ak)

	mapp.Router().AddRoute("gov", []*sdk.KVStoreKey{keyGov}, NewHandler(gk))

	mapp.SetEndBlocker(getEndBlocker(gk))
	mapp.SetInitChainer(getInitChainer(mapp, gk, sk))

	require.NoError(t, mapp.CompleteSetup(keyGov))

	coin, _ := sdk.IrisCoinType.ConvertToMinDenomCoin(fmt.Sprintf("%d%s", 1042, sdk.Iris))
	genAccs, addrs, pubKeys, privKeys := mock.CreateGenAccounts(numGenAccs, sdk.Coins{coin})

	mock.SetGenesis(mapp, genAccs)

	return mapp, gk, sk, addrs, pubKeys, privKeys
}

// gov and stake endblocker
func getEndBlocker(keeper Keeper) sdk.EndBlocker {
	return func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		tags := EndBlocker(ctx, keeper)
		return abci.ResponseEndBlock{
			Tags: tags,
		}
	}
}

// gov and stake initchainer
func getInitChainer(mapp *mock.App, keeper Keeper, stakeKeeper stake.Keeper) sdk.InitChainer {
	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		mapp.InitChainer(ctx, req)

		stakeGenesis := stake.DefaultGenesisState()

		validators, err := stake.InitGenesis(ctx, stakeKeeper, stakeGenesis)
		if err != nil {
			panic(err)
		}
		InitGenesis(ctx, keeper, GenesisState{
			Params: DefaultParams(),
		})
		return abci.ResponseInitChain{
			Validators: validators,
		}
	}
}

// Sorts Addresses
func SortAddresses(addrs []sdk.AccAddress) {
	var byteAddrs [][]byte
	for _, addr := range addrs {
		byteAddrs = append(byteAddrs, addr.Bytes())
	}
	SortByteArrays(byteAddrs)
	for i, byteAddr := range byteAddrs {
		addrs[i] = byteAddr
	}
}

// implement `Interface` in sort package.
type sortByteArrays [][]byte

func (b sortByteArrays) Len() int {
	return len(b)
}

func (b sortByteArrays) Less(i, j int) bool {
	// bytes package already implements Comparable for []byte.
	switch bytes.Compare(b[i], b[j]) {
	case -1:
		return true
	case 0, 1:
		return false
	default:
		log.Panic("not fail-able with `bytes.Comparable` bounded [-1, 1].")
		return false
	}
}

func (b sortByteArrays) Swap(i, j int) {
	b[j], b[i] = b[i], b[j]
}

// Public
func SortByteArrays(src [][]byte) [][]byte {
	sorted := sortByteArrays(src)
	sort.Sort(sorted)
	return sorted
}
