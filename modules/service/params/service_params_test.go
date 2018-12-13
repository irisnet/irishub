package serviceparams

import (
	"github.com/irisnet/irishub/store"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/params"
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

	subspace := paramKeeper.Subspace("Service").WithTypeTable(params.NewTypeTable(
		MaxRequestTimeoutParameter.GetStoreKey(), int64(0),
		MinDepositMultipleParameter.GetStoreKey(), int64(0),
	))

	MaxRequestTimeoutParameter.SetReadWriter(subspace)
	find := MaxRequestTimeoutParameter.LoadValue(ctx)
	require.Equal(t, find, false)

	MaxRequestTimeoutParameter.InitGenesis(int64(12345))
	require.Equal(t, int64(12345), MaxRequestTimeoutParameter.Value)

	MaxRequestTimeoutParameter.LoadValue(ctx)
	require.Equal(t, int64(12345), MaxRequestTimeoutParameter.Value)

	MaxRequestTimeoutParameter.Value = 30
	MaxRequestTimeoutParameter.SaveValue(ctx)

	MaxRequestTimeoutParameter.LoadValue(ctx)
	require.Equal(t, int64(30), MaxRequestTimeoutParameter.Value)
}

func TestMinProviderDepositParameter(t *testing.T) {
	skey := sdk.NewKVStoreKey("params")
	tkeyParams := sdk.NewTransientStoreKey("transient_params")

	ctx := defaultContext(skey, tkeyParams)
	cdc := codec.New()

	paramKeeper := params.NewKeeper(
		cdc,
		skey, tkeyParams,
	)

	subspace := paramKeeper.Subspace("Sig").WithTypeTable(params.NewTypeTable(
		MaxRequestTimeoutParameter.GetStoreKey(), int64(0),
		MinDepositMultipleParameter.GetStoreKey(), int64(0),
	))

	MinDepositMultipleParameter.SetReadWriter(subspace)
	find := MinDepositMultipleParameter.LoadValue(ctx)
	require.Equal(t, find, false)

	MinDepositMultipleParameter.InitGenesis(int64(12345))
	require.Equal(t, int64(12345), MinDepositMultipleParameter.Value)

	MinDepositMultipleParameter.LoadValue(ctx)
	require.Equal(t, int64(12345), MinDepositMultipleParameter.Value)

	MinDepositMultipleParameter.Value = 30
	MinDepositMultipleParameter.SaveValue(ctx)

	MinDepositMultipleParameter.LoadValue(ctx)
	require.Equal(t, int64(30), MinDepositMultipleParameter.Value)
}
