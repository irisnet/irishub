package types

import (
	"testing"

	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/auth"
	"github.com/irisnet/irishub/modules/bank"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
)

func TestPoolEqual(t *testing.T) {
	p1 := InitialBondedPool()
	p2 := InitialBondedPool()
	require.True(t, p1.Equal(p2))
	p2.BondedTokens = sdk.NewDec(3)
	require.False(t, p1.Equal(p2))
}

func TestAddBondedTokens(t *testing.T) {

	ms, authKey := setupMultiStore()
	cdc := codec.New()
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	accountKeeper := auth.NewAccountKeeper(cdc, authKey, auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper)

	bondedPool := InitialBondedPool()
	poolA := Pool{
		BondedPool: bondedPool,
		BankKeeper: bankKeeper,
	}
	poolA.BondedPool.BondedTokens = sdk.NewDec(10)
	poolA.BankKeeper.IncreaseLoosenToken(ctx, sdk.Coins{sdk.NewCoin(StakeDenom, sdk.NewInt(10))})

	poolA = poolA.loosenTokenToBonded(ctx, sdk.NewDec(10))

	require.True(sdk.DecEq(t, sdk.NewDec(20), poolA.BondedPool.BondedTokens))
	require.True(sdk.DecEq(t, sdk.NewDec(0), poolA.GetLoosenTokenAmount(ctx)))
}

func TestRemoveBondedTokens(t *testing.T) {
	ms, authKey := setupMultiStore()
	cdc := codec.New()
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	accountKeeper := auth.NewAccountKeeper(cdc, authKey, auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper)

	bondedPool := InitialBondedPool()
	poolA := Pool{
		BondedPool: bondedPool,
		BankKeeper: bankKeeper,
	}
	poolA.BondedPool.BondedTokens = sdk.NewDec(10)
	poolA.BankKeeper.IncreaseLoosenToken(ctx, sdk.Coins{sdk.NewCoin(StakeDenom, sdk.NewInt(10))})

	poolA = poolA.bondedTokenToLoosen(ctx, sdk.NewDec(5))

	require.True(sdk.DecEq(t, sdk.NewDec(5), poolA.BondedPool.BondedTokens))
	require.True(sdk.DecEq(t, sdk.NewDec(15), poolA.GetLoosenTokenAmount(ctx)))
}
