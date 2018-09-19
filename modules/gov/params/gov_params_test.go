package govparams

import (
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"testing"
	"fmt"
	"github.com/irisnet/irishub/types"
)

func defaultContext(key sdk.StoreKey) sdk.Context {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	ctx := sdk.NewContext(cms, abci.Header{}, false, log.NewNopLogger())
	return ctx
}

func TestDepositProcedureParam(t *testing.T) {
	skey := sdk.NewKVStoreKey("params")
	ctx := defaultContext(skey)
	paramKeeper := params.NewKeeper(wire.NewCodec(), skey)

	p1deposit, _ := types.NewDefaultCoinType("iris").ConvertToMinCoin(fmt.Sprintf("%d%s", 10, "iris"))
	p2Deposit, _ := types.NewDefaultCoinType("iris").ConvertToMinCoin(fmt.Sprintf("%d%s", 200, "iris"))
	p1 := DepositProcedure{
		MinDeposit:       sdk.Coins{p1deposit},
		MaxDepositPeriod: 1440}

	p2 := DepositProcedure{
		MinDeposit:       sdk.Coins{p2Deposit},
		MaxDepositPeriod: 1440}

	DepositProcedureParameter.SetReadWriter(paramKeeper.Setter())
	find := DepositProcedureParameter.LoadValue(ctx)
	require.Equal(t, find, false)

	DepositProcedureParameter.InitGenesis(nil)
	require.Equal(t, p1, DepositProcedureParameter.Value)

	require.Equal(t, DepositProcedureParameter.ToJson(), "{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"10000000000000000000\"}],\"max_deposit_period\":1440}")

	DepositProcedureParameter.Update(ctx, "{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"200000000000000000000\"}],\"max_deposit_period\":1440}")

	require.NotEqual(t, p1, DepositProcedureParameter.Value)
	require.Equal(t, p2, DepositProcedureParameter.Value)

	result := DepositProcedureParameter.Valid("{\"min_deposit\":[{\"denom\":\"atom\",\"amount\":\"200000000000000000000\"}],\"max_deposit_period\":1440}")
	require.Error(t, result)

	DepositProcedureParameter.InitGenesis(p2)
	require.Equal(t, p2, DepositProcedureParameter.Value)
	DepositProcedureParameter.InitGenesis(p1)
	require.Equal(t, p1, DepositProcedureParameter.Value)

	DepositProcedureParameter.LoadValue(ctx)
	require.Equal(t, p2, DepositProcedureParameter.Value)

}

func TestDepositProcedureParamValid(t *testing.T) {
	skey := sdk.NewKVStoreKey("params")
	ctx := defaultContext(skey)
	paramKeeper := params.NewKeeper(wire.NewCodec(), skey)

	p1deposit, _ := types.NewDefaultCoinType("iris").ConvertToMinCoin(fmt.Sprintf("%d%s", 10, "iris"))
	p1 := DepositProcedure{
		MinDeposit:       sdk.Coins{p1deposit},
		MaxDepositPeriod: 1440}

	DepositProcedureParameter.SetReadWriter(paramKeeper.Setter())
	find := DepositProcedureParameter.LoadValue(ctx)
	require.Equal(t, find, false)

	DepositProcedureParameter.InitGenesis(nil)
	require.Equal(t, p1, DepositProcedureParameter.Value)

	result := DepositProcedureParameter.Valid("{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"2000000000000000000\"}],\"max_deposit_period\":1440}")
	require.NoError(t, result)
	result = DepositProcedureParameter.Valid("{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"2000000000000000000000\"}],\"max_deposit_period\":1440}")
	require.Error(t, result)
	result = DepositProcedureParameter.Valid("{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"200000000000000000\"}],\"max_deposit_period\":1440}")
	require.Error(t, result)

	result = DepositProcedureParameter.Valid("{\"min_deposit\":[{\"denom\":\"iris-att\",\"amount\":\"2000000000000000000\"}],\"max_deposit_period\":1440}")
	require.Error(t, result)

	result = DepositProcedureParameter.Valid("{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"2000000000000000000\"}],\"max_deposit_period\":1440}")
	require.NoError(t, result)

	result = DepositProcedureParameter.Valid("{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"2000000000000000000\"}],\"max_deposit_period\":1}")
	require.Error(t, result)

	result = DepositProcedureParameter.Valid("{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"2000000000000000000\"}],\"max_deposit_period\":1440000}")
	require.Error(t, result)

}