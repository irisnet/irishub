package serviceparams

import (
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/params"
	"github.com/irisnet/irishub/store"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"testing"
)

func defaultContext(key sdk.StoreKey, tkeyParams *sdk.TransientStoreKey) sdk.Context {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)

	cms.LoadLatestVersion()
	ctx := sdk.NewContext(cms, abci.Header{}, false, log.NewNopLogger())
	return ctx
}

func TestMaxRequestTimeoutParameter(t *testing.T) {
	skey := sdk.NewKVStoreKey("params")
	tkeyParams := sdk.NewTransientStoreKey("transient_params")

	ctx := defaultContext(skey, tkeyParams)
	cdc := codec.New()

	paramKeeper := params.NewKeeper(
		cdc,
		skey, tkeyParams,
	)

	subspace := paramKeeper.Subspace("Gov").WithTypeTable(params.NewTypeTable(
		ServiceParameter.GetStoreKey(), Params{},
	))

	ServiceParameter.SetReadWriter(subspace)
	find := ServiceParameter.LoadValue(ctx)
	require.Equal(t, find, false)

	ServiceParameter.InitGenesis(NewSericeParams())
	require.Equal(t, int64(100), GetMaxRequestTimeout(ctx))

	SetMaxRequestTimeout(ctx, 1000)
	require.Equal(t, int64(1000), GetMaxRequestTimeout(ctx))
}
