package gov

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"log"
	"sort"
	"testing"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"

	"fmt"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/bank"
	"github.com/irisnet/irishub/modules/mock"
	"github.com/irisnet/irishub/modules/stake"
	"github.com/irisnet/irishub/modules/gov/params"
	"github.com/irisnet/irishub/types"
	stakeTypes "github.com/irisnet/irishub/modules/stake/types"
	"github.com/irisnet/irishub/modules/distribution"
	"github.com/irisnet/irishub/modules/params"
	"github.com/irisnet/irishub/modules/auth"
	"github.com/irisnet/irishub/modules/guardian"
	protocolKeeper "github.com/irisnet/irishub/app/protocol/keeper"
)

// initialize the mock application for this module
func getMockApp(t *testing.T, numGenAccs int) (*mock.App, Keeper, stake.Keeper, []sdk.AccAddress, []crypto.PubKey, []crypto.PrivKey) {
	mapp := mock.NewApp()

	stake.RegisterCodec(mapp.Cdc)
	RegisterCodec(mapp.Cdc)

	keyGov := sdk.NewKVStoreKey("gov")
	keyDistr := sdk.NewKVStoreKey("distr")
    keyProtocol := sdk.NewKVStoreKey("protocol")

    pk := protocolKeeper.NewKeeper(mapp.Cdc,keyProtocol)
	paramsKeeper := params.NewKeeper(
		mapp.Cdc,
		sdk.NewKVStoreKey("params"),
		sdk.NewTransientStoreKey("transient_params"),
	)
	feeCollectionKeeper := auth.NewFeeCollectionKeeper(
		mapp.Cdc,
		sdk.NewKVStoreKey("fee"),
	)

	ck := bank.NewBaseKeeper(mapp.AccountKeeper)
	sk := stake.NewKeeper(
		mapp.Cdc,
		mapp.KeyStake, mapp.TkeyStake,
		mapp.BankKeeper, mapp.ParamsKeeper.Subspace(stake.DefaultParamspace),
		stake.DefaultCodespace)
	dk := distribution.NewKeeper(mapp.Cdc, keyDistr, paramsKeeper.Subspace(distribution.DefaultParamspace), ck, sk, feeCollectionKeeper, DefaultCodespace)
	guardianKeeper := guardian.NewKeeper(mapp.Cdc, sdk.NewKVStoreKey("guardian"), guardian.DefaultCodespace)
	gk := NewKeeper(mapp.Cdc, keyGov, dk, ck, guardianKeeper, sk, pk, DefaultCodespace)

	mapp.Router().AddRoute("gov", []*sdk.KVStoreKey{keyGov}, NewHandler(gk))

	mapp.SetEndBlocker(getEndBlocker(gk))
	mapp.SetInitChainer(getInitChainer(mapp, gk, sk))

	require.NoError(t, mapp.CompleteSetup(keyGov))

	coin, _ := types.NewDefaultCoinType(stakeTypes.StakeDenomName).ConvertToMinCoin(fmt.Sprintf("%d%s", 1042, stakeTypes.StakeDenomName))
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
		stakeGenesis.Params.BondDenom = stakeTypes.StakeDenom
		stakeGenesis.Pool.LooseTokens = sdk.NewDecFromInt(sdk.NewInt(100000))

		validators, err := stake.InitGenesis(ctx, stakeKeeper, stakeGenesis)
		if err != nil {
			panic(err)
		}
		ct := types.NewDefaultCoinType(stakeTypes.StakeDenomName)
		minDeposit, _ := ct.ConvertToMinCoin(fmt.Sprintf("%d%s", 10, stakeTypes.StakeDenomName))
		InitGenesis(ctx, keeper, GenesisState{
			StartingProposalID: 1,
			DepositProcedure: govparams.DepositProcedure{
				MinDeposit:       sdk.Coins{minDeposit},
				MaxDepositPeriod: 1440,
			},
			VotingProcedure: govparams.VotingProcedure{
				VotingPeriod: 30,
			},
			TallyingProcedure: govparams.TallyingProcedure{
				Threshold:     sdk.NewDecWithPrec(5, 1),
				Veto:          sdk.NewDecWithPrec(334, 3),
				Participation: sdk.NewDecWithPrec(667, 3),
			},
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
