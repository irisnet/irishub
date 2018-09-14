package GovParams

import (
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/irisnet/irishub/modules/gov"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"testing"
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

	p1 := gov.DepositProcedure{
		MinDeposit:       sdk.Coins{sdk.NewInt64Coin("iris", 10)},
		MaxDepositPeriod: 1440}

	p2 := gov.DepositProcedure{
		MinDeposit:       sdk.Coins{sdk.NewInt64Coin("iris", 30)},
		MaxDepositPeriod: 1440}

	DepositProcedureParameter.SetReadWriter(paramKeeper.Setter())
	find := DepositProcedureParameter.LoadValue(ctx)
	require.Equal(t, find, false)

	DepositProcedureParameter.InitGenesis(nil)
	require.Equal(t, p1, DepositProcedureParameter.Value)

	require.Equal(t, DepositProcedureParameter.ToJson(), "{\"min_deposit\":[{\"denom\":\"iris\",\"amount\":\"10\"}],\"max_deposit_period\":1440}")
	DepositProcedureParameter.Update(ctx, "{\"min_deposit\":[{\"denom\":\"iris\",\"amount\":\"30\"}],\"max_deposit_period\":1440}")
	require.NotEqual(t, p1, DepositProcedureParameter.Value)
	require.Equal(t, p2, DepositProcedureParameter.Value)

	result := DepositProcedureParameter.Valid("{\"min_deposit\":[{\"denom\":\"atom\",\"amount\":\"30\"}],\"max_deposit_period\":1440}")
	require.Error(t, result)

	DepositProcedureParameter.InitGenesis(p2)
	require.Equal(t, p2, DepositProcedureParameter.Value)
	DepositProcedureParameter.InitGenesis(p1)
	require.Equal(t, p1, DepositProcedureParameter.Value)

	DepositProcedureParameter.LoadValue(ctx)
	require.Equal(t, p2, DepositProcedureParameter.Value)

}
