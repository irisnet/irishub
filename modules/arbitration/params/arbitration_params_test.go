package arbitrationparams

import (
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"testing"
	"time"
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
		ComplaintRetrospectParameter.GetStoreKey(), time.Duration(0),
		ArbitrationTimelimitParameter.GetStoreKey(), time.Duration(0),
	))

	ComplaintRetrospectParameter.SetReadWriter(subspace)
	find := ComplaintRetrospectParameter.LoadValue(ctx)
	require.Equal(t, find, false)

	ComplaintRetrospectParameter.InitGenesis(5 * time.Second)
	require.Equal(t, 5 * time.Second, ComplaintRetrospectParameter.Value)

	ComplaintRetrospectParameter.LoadValue(ctx)
	require.Equal(t, 5 * time.Second, ComplaintRetrospectParameter.Value)

	ComplaintRetrospectParameter.Value = 10 * time.Second
	ComplaintRetrospectParameter.SaveValue(ctx)

	ComplaintRetrospectParameter.LoadValue(ctx)
	require.Equal(t, 10 * time.Second, ComplaintRetrospectParameter.Value)
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

	subspace := paramKeeper.Subspace("Service").WithTypeTable(params.NewTypeTable(
		ComplaintRetrospectParameter.GetStoreKey(), time.Duration(0),
		ArbitrationTimelimitParameter.GetStoreKey(), time.Duration(0),
	))

	ArbitrationTimelimitParameter.SetReadWriter(subspace)
	find := ArbitrationTimelimitParameter.LoadValue(ctx)
	require.Equal(t, find, false)

	ArbitrationTimelimitParameter.InitGenesis(5 * time.Second)
	require.Equal(t, 5 * time.Second, ArbitrationTimelimitParameter.Value)

	ArbitrationTimelimitParameter.LoadValue(ctx)
	require.Equal(t, 5 * time.Second, ArbitrationTimelimitParameter.Value)

	ArbitrationTimelimitParameter.Value = 10 * time.Second
	ArbitrationTimelimitParameter.SaveValue(ctx)

	ArbitrationTimelimitParameter.LoadValue(ctx)
	require.Equal(t, 10 * time.Second, ArbitrationTimelimitParameter.Value)
}
