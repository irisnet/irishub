package asset

import (
	"encoding/json"
	"testing"

	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v1/bank"
	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/guardian"
	"github.com/irisnet/irishub/store"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey, *sdk.KVStoreKey, *sdk.TransientStoreKey) {
	db := dbm.NewMemDB()
	assetKey := sdk.NewKVStoreKey("assetKey")
	paramskey := sdk.NewKVStoreKey("params")
	paramsTkey := sdk.NewTransientStoreKey("transient_params")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(assetKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(paramskey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(paramsTkey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	return ms, assetKey, paramskey, paramsTkey
}

func TestKeeper_IssueAsset(t *testing.T) {
	ms, assetKey, paramskey, paramsTkey := setupMultiStore()

	cdc := codec.New()
	RegisterCodec(cdc)
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	pk := params.NewKeeper(cdc, paramskey, paramsTkey)
	keeper := NewKeeper(cdc, assetKey, bank.BaseKeeper{}, guardian.Keeper{}, DefaultCodespace, pk.Subspace(DefaultParamSpace))

	addr := sdk.AccAddress([]byte("addr1"))

	ft := NewFungibleToken(0x00, "c", "d", "e", 1, "f", sdk.NewIntWithDecimal(1, 0), sdk.NewIntWithDecimal(1, 0), sdk.NewIntWithDecimal(1, 0), true, addr)
	_, err := keeper.IssueAsset(ctx, ft)
	assert.NoError(t, err)

	assert.True(t, keeper.HasAsset(ctx, "i.d"))

	asset, found := keeper.getAsset(ctx, "i.d")
	assert.True(t, found)

	assert.Equal(t, ft.GetDenom(), asset.GetDenom())
	assert.Equal(t, ft.Owner, ft.Owner)

	msgJson, _ := json.Marshal(ft)
	assetJson, _ := json.Marshal(asset)
	assert.Equal(t, msgJson, assetJson)
}

func TestKeeper_IssueGatewayAsset(t *testing.T) {
	ms, assetKey, paramskey, paramsTkey := setupMultiStore()

	cdc := codec.New()
	RegisterCodec(cdc)
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	pk := params.NewKeeper(cdc, paramskey, paramsTkey)
	keeper := NewKeeper(cdc, assetKey, bank.BaseKeeper{}, guardian.Keeper{}, DefaultCodespace, pk.Subspace(DefaultParamSpace))

	owner := sdk.AccAddress([]byte("owner"))
	gatewayOwner := sdk.AccAddress([]byte("gatewayOwner"))
	moniker := "moniker"
	identity := "identity"
	details := "details"
	website := "website"

	// construct a test gateway
	gateway := Gateway{
		Owner:    gatewayOwner,
		Moniker:  moniker,
		Identity: identity,
		Details:  details,
		Website:  website,
	}
	gatewayAsset := BaseAsset{FUNGIBLE, GATEWAY, "moniker", "d", "e", 1, "f", sdk.NewIntWithDecimal(1, 0), sdk.NewIntWithDecimal(1, 0), sdk.NewIntWithDecimal(1, 0), true, owner}
	gatewayAsset1 := BaseAsset{FUNGIBLE, GATEWAY, "moniker", "d", "e", 1, "f", sdk.NewIntWithDecimal(1, 0), sdk.NewIntWithDecimal(1, 0), sdk.NewIntWithDecimal(1, 0), true, gatewayOwner}

	// unknown gateway moniker
	_, err := keeper.IssueAsset(ctx, gatewayAsset1)
	assert.Error(t, err)
	asset, found := keeper.getAsset(ctx, "moniker.d")
	assert.False(t, found)
	assert.Nil(t, asset)

	// no unauthorized creator
	keeper.SetGateway(ctx, gateway)
	_, err = keeper.IssueAsset(ctx, gatewayAsset)
	assert.Error(t, err)
	asset, found = keeper.getAsset(ctx, "moniker.d")
	assert.False(t, found)
	assert.Nil(t, asset)

	_, err = keeper.IssueAsset(ctx, gatewayAsset1)
	assert.NoError(t, err)
	asset, found = keeper.getAsset(ctx, "moniker.d")
	assert.True(t, found)
	assert.Equal(t, "moniker.d", asset.GetUniqueID())
}

func TestCreateKeeper(t *testing.T) {
	ms, assetKey, paramskey, paramsTkey := setupMultiStore()

	cdc := codec.New()
	RegisterCodec(cdc)
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	bankKeeper := bank.BaseKeeper{}
	guardianKeeper := guardian.Keeper{}
	paramsKeeper := params.NewKeeper(cdc, paramskey, paramsTkey)

	createKeeper := NewKeeper(cdc, assetKey, bankKeeper, guardianKeeper, DefaultCodespace, paramsKeeper.Subspace(DefaultParamSpace))

	// define variables
	owner := sdk.AccAddress([]byte("owner"))
	moniker := "moniker"
	identity := "identity"
	details := "details"
	website := "website"

	// construct a test gateway
	gateway := Gateway{
		Owner:    owner,
		Moniker:  moniker,
		Identity: identity,
		Details:  details,
		Website:  website,
	}

	// assert the gateway of the given moniker does not exist at the beginning
	require.False(t, createKeeper.HasGateway(ctx, moniker))

	// create a gateway and asset that the gateway exists now
	createKeeper.SetGateway(ctx, gateway)
	require.True(t, createKeeper.HasGateway(ctx, moniker))

	// asset GetGateway will return the previous gateway
	newGateway, _ := createKeeper.GetGateway(ctx, moniker)
	require.Equal(t, gateway, newGateway)
}
