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
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/irisnet/irishub/modules/mock"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/irisnet/irishub/modules/gov/params"
	"github.com/irisnet/irishub/types"
	sdkParams "github.com/cosmos/cosmos-sdk/x/params"
	"github.com/irisnet/irishub/modules/iparam"
	"github.com/irisnet/irishub/modules/upgrade/params"
)

// initialize the mock application for this module
func getMockApp(t *testing.T, numGenAccs int) (*mock.App, Keeper, stake.Keeper, []sdk.AccAddress, []crypto.PubKey, []crypto.PrivKey) {
	mapp := mock.NewApp()

	stake.RegisterWire(mapp.Cdc)
	RegisterWire(mapp.Cdc)

	keyGlobalParams := sdk.NewKVStoreKey("params")
	keyStake := sdk.NewKVStoreKey("stake")
	keyGov := sdk.NewKVStoreKey("gov")

	ck := bank.NewKeeper(mapp.AccountMapper)
	sk := stake.NewKeeper(mapp.Cdc, keyStake, ck, mapp.RegisterCodespace(stake.DefaultCodespace))
	gk := NewKeeper(mapp.Cdc, keyGov, ck, sk, DefaultCodespace)
	pk := sdkParams.NewKeeper(mapp.Cdc, keyGlobalParams)

	mapp.Router().AddRoute("gov", []*sdk.KVStoreKey{keyGov}, NewHandler(gk))

	iparam.SetParamReadWriter(pk.Setter(),
		&govparams.DepositProcedureParameter,
		&govparams.VotingProcedureParameter,
		&govparams.TallyingProcedureParameter,
		&upgradeparams.CurrentUpgradeProposalIdParameter,
		&upgradeparams.ProposalAcceptHeightParameter)

	iparam.RegisterGovParamMapping(&govparams.DepositProcedureParameter,
		&govparams.VotingProcedureParameter,
		&govparams.TallyingProcedureParameter,)

	mapp.SetEndBlocker(getEndBlocker(gk))
	mapp.SetInitChainer(getInitChainer(mapp, gk, sk))

	require.NoError(t, mapp.CompleteSetup([]*sdk.KVStoreKey{keyStake, keyGov, keyGlobalParams}))

	coin, _ := types.NewDefaultCoinType("iris").ConvertToMinCoin(fmt.Sprintf("%d%s", 1042, "iris"))
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
		stakeGenesis.Params.BondDenom = "iris-atto"
		stakeGenesis.Pool.LooseTokens = sdk.NewRat(100000)

		validators, err := stake.InitGenesis(ctx, stakeKeeper, stakeGenesis)
		if err != nil {
			panic(err)
		}
		ct := types.NewDefaultCoinType("iris")
		minDeposit, _ := ct.ConvertToMinCoin(fmt.Sprintf("%d%s", 10, "iris"))
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
				Threshold:         sdk.NewRat(1, 2),
				Veto:              sdk.NewRat(1, 3),
				GovernancePenalty: sdk.NewRat(1, 100),
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
