package serviceparams

import (
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/irisnet/irishub/modules/params"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"testing"
	"github.com/irisnet/irishub/types"
	"fmt"
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
		MinProviderDepositParameter.GetStoreKey(), sdk.Coins{},
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
		MinProviderDepositParameter.GetStoreKey(), sdk.Coins{},
	))

	MinProviderDepositParameter.SetReadWriter(subspace)
	find := MinProviderDepositParameter.LoadValue(ctx)
	require.Equal(t, find, false)

	p1deposit, _ := types.NewDefaultCoinType("iris").ConvertToMinCoin(fmt.Sprintf("%d%s", 10, "iris"))
	p2deposit, _ := types.NewDefaultCoinType("iris").ConvertToMinCoin(fmt.Sprintf("%d%s", 1000, "iris"))
	MinProviderDepositParameter.InitGenesis(sdk.Coins{p1deposit})
	require.Equal(t, sdk.Coins{p1deposit}, MinProviderDepositParameter.Value)

	MinProviderDepositParameter.LoadValue(ctx)
	require.Equal(t, sdk.Coins{p1deposit}, MinProviderDepositParameter.Value)

	MinProviderDepositParameter.Value = sdk.Coins{p2deposit}
	MinProviderDepositParameter.SaveValue(ctx)

	MinProviderDepositParameter.LoadValue(ctx)
	require.Equal(t, sdk.Coins{p2deposit}, MinProviderDepositParameter.Value)
}
