package iservice

import (
	"bytes"
	"log"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/irisnet/irishub/simulation/mock"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/irisnet/irishub/types"
	"fmt"
)

// initialize the mock application for this module
func getMockApp(t *testing.T, numGenAccs int) (*mock.App, Keeper, stake.Keeper, []sdk.AccAddress, []crypto.PubKey, []crypto.PrivKey) {
	mapp := mock.NewApp()

	stake.RegisterCodec(mapp.Cdc)
	RegisterCodec(mapp.Cdc)

	keyService := sdk.NewKVStoreKey("iservice")

	ck := bank.NewBaseKeeper(mapp.AccountKeeper)
	sk := stake.NewKeeper(
		mapp.Cdc,
		mapp.KeyStake, mapp.TkeyStake,
		mapp.BankKeeper, mapp.ParamsKeeper.Subspace(stake.DefaultParamspace),
		mapp.RegisterCodespace(stake.DefaultCodespace))
	ik := NewKeeper(mapp.Cdc, keyService,ck, DefaultCodespace)

	mapp.Router().AddRoute("iservice", []*sdk.KVStoreKey{keyService}, NewHandler(ik))

	mapp.SetEndBlocker(getEndBlocker())
	mapp.SetInitChainer(getInitChainer(mapp, sk))

	require.NoError(t, mapp.CompleteSetup(keyService))

	coin, _ := types.NewDefaultCoinType("iris").ConvertToMinCoin(fmt.Sprintf("%d%s", 1042, "iris"))
	genAccs, addrs, pubKeys, privKeys := mock.CreateGenAccounts(numGenAccs, sdk.Coins{coin})

	mock.SetGenesis(mapp, genAccs)

	return mapp, ik, sk, addrs, pubKeys, privKeys
}

// gov and stake endblocker
func getEndBlocker() sdk.EndBlocker {
	return func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		return abci.ResponseEndBlock{
		}
	}
}

// gov and stake initchainer
func getInitChainer(mapp *mock.App, stakeKeeper stake.Keeper) sdk.InitChainer {
	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		mapp.InitChainer(ctx, req)

		stakeGenesis := stake.DefaultGenesisState()
		stakeGenesis.Pool.LooseTokens = sdk.NewDec(100000)

		validators, err := stake.InitGenesis(ctx, stakeKeeper, stakeGenesis)
		if err != nil {
			panic(err)
		}
		InitGenesis(ctx, DefaultGenesisState())
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
