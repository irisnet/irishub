package rand

import (
	"bytes"
	"fmt"
	"log"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/irisnet/irishub/app/v1/stake"
	"github.com/irisnet/irishub/server/mock"
	sdk "github.com/irisnet/irishub/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
)

// initialize the mock application for this module
func getMockApp(t *testing.T, numGenAccs int) (*mock.App, Keeper, stake.Keeper, []sdk.AccAddress, []crypto.PubKey, []crypto.PrivKey) {
	mapp := mock.NewApp()

	stake.RegisterCodec(mapp.Cdc)
	RegisterCodec(mapp.Cdc)

	keyRand := sdk.NewKVStoreKey("rand")

	sk := stake.NewKeeper(
		mapp.Cdc,
		mapp.KeyStake, mapp.TkeyStake,
		mapp.BankKeeper, mapp.ParamsKeeper.Subspace(stake.DefaultParamspace),
		stake.DefaultCodespace,
		stake.NopMetrics())
	rk := NewKeeper(mapp.Cdc, keyRand, DefaultCodespace)

	mapp.Router().AddRoute("rand", []*sdk.KVStoreKey{keyRand}, NewHandler(rk))

	mapp.SetBeginBlocker(getBeginBlocker(rk))
	mapp.SetEndBlocker(getEndBlocker())
	mapp.SetInitChainer(getInitChainer(mapp, rk, sk))

	require.NoError(t, mapp.CompleteSetup(keyRand))

	coin, _ := sdk.IrisCoinType.ConvertToMinDenomCoin(fmt.Sprintf("%d%s", 1042, sdk.Iris))
	genAccs, addrs, pubKeys, privKeys := mock.CreateGenAccounts(numGenAccs, sdk.Coins{coin})

	mock.SetGenesis(mapp, genAccs)

	return mapp, rk, sk, addrs, pubKeys, privKeys
}

// rand endblocker
func getEndBlocker() sdk.EndBlocker {
	return func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		return abci.ResponseEndBlock{}
	}
}

// rand beginblocker
func getBeginBlocker(randKeeper Keeper) sdk.BeginBlocker {
	return func(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
		tags := BeginBlocker(ctx, req, randKeeper)

		return abci.ResponseBeginBlock{
			Tags: tags,
		}
	}
}

// rand initchainer
func getInitChainer(mapp *mock.App, randKeeper Keeper, stakeKeeper stake.Keeper) sdk.InitChainer {
	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		mapp.InitChainer(ctx, req)

		stakeGenesis := stake.DefaultGenesisState()

		validators, err := stake.InitGenesis(ctx, stakeKeeper, stakeGenesis)
		if err != nil {
			panic(err)
		}

		InitGenesis(ctx, randKeeper, DefaultGenesisState())
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
