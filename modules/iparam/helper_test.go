package iparam

import (
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/irisnet/irishub/modules/gov/params"
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

func TestInitGenesisParameter(t *testing.T) {

	skey := sdk.NewKVStoreKey("params")
	ctx := defaultContext(skey)
	paramKeeper := params.NewKeeper(wire.NewCodec(), skey)

	p1 := govparams.DepositProcedure{
		MinDeposit:       sdk.Coins{sdk.NewInt64Coin("iris", 10)},
		MaxDepositPeriod: 1440}

	p2 := govparams.DepositProcedure{
		MinDeposit:       sdk.Coins{sdk.NewInt64Coin("iris", 20)},
		MaxDepositPeriod: 1440}

	SetParamReadWriter(paramKeeper.Setter(),&govparams.DepositProcedureParameter,&govparams.DepositProcedureParameter)

	InitGenesisParameter(&govparams.DepositProcedureParameter, ctx, nil)

	require.Equal(t, p1, govparams.DepositProcedureParameter.Value)

	require.Equal(t, govparams.DepositProcedureParameter.ToJson(), "{\"min_deposit\":[{\"denom\":\"iris\",\"amount\":\"10\"}],\"max_deposit_period\":1440}")

	InitGenesisParameter(&govparams.DepositProcedureParameter, ctx, p2)

	require.Equal(t, p1, govparams.DepositProcedureParameter.Value)
}

func TestRegisterParamMapping(t *testing.T) {

	skey := sdk.NewKVStoreKey("params")
	ctx := defaultContext(skey)
	paramKeeper := params.NewKeeper(wire.NewCodec(), skey)

	p1 := govparams.DepositProcedure{
		MinDeposit:       sdk.Coins{sdk.NewInt64Coin("iris", 10)},
		MaxDepositPeriod: 1440}

	p2 := govparams.DepositProcedure{
		MinDeposit:       sdk.Coins{sdk.NewInt64Coin("iris", 30)},
		MaxDepositPeriod: 1440}

	SetParamReadWriter(paramKeeper.Setter(),&govparams.DepositProcedureParameter, &govparams.DepositProcedureParameter)
	RegisterGovParamMapping(&govparams.DepositProcedureParameter)
	InitGenesisParameter(&govparams.DepositProcedureParameter, ctx, nil)

	require.Equal(t, ParamMapping[govparams.DepositProcedureParameter.GetStoreKey()].ToJson(), "{\"min_deposit\":[{\"denom\":\"iris\",\"amount\":\"10\"}],\"max_deposit_period\":1440}")
	require.Equal(t, p1, govparams.DepositProcedureParameter.Value)

	ParamMapping[govparams.DepositProcedureParameter.GetStoreKey()].Update(ctx, "{\"min_deposit\":[{\"denom\":\"iris\",\"amount\":\"30\"}],\"max_deposit_period\":1440}")
	govparams.DepositProcedureParameter.LoadValue(ctx)
	require.Equal(t, p2, govparams.DepositProcedureParameter.Value)
}
