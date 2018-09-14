package parameter

import (
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/irisnet/irishub/modules/gov"
	govParam "github.com/irisnet/irishub/modules/gov/params"
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

	p1 := gov.DepositProcedure{
		MinDeposit:       sdk.Coins{sdk.NewInt64Coin("iris", 10)},
		MaxDepositPeriod: 1440}

	p2 := gov.DepositProcedure{
		MinDeposit:       sdk.Coins{sdk.NewInt64Coin("iris", 20)},
		MaxDepositPeriod: 1440}

	govParam.DepositProcedureParameter.SetReadWriter(paramKeeper.Setter())

	InitGenesisParameter(&govParam.DepositProcedureParameter, ctx,nil)

	require.Equal(t, p1, govParam.DepositProcedureParameter.Value)

	require.Equal(t, govParam.DepositProcedureParameter.ToJson(), "{\"min_deposit\":[{\"denom\":\"iris\",\"amount\":\"10\"}],\"max_deposit_period\":1440}")

	InitGenesisParameter(&govParam.DepositProcedureParameter, ctx, p2)

	require.Equal(t, p1, govParam.DepositProcedureParameter.Value)
}