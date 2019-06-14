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

func TestCreateGatewayKeeper(t *testing.T) {
	ms, assetKey, paramskey, paramsTkey := setupMultiStore()

	cdc := codec.New()
	RegisterCodec(cdc)
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	bankKeeper := bank.BaseKeeper{}
	guardianKeeper := guardian.Keeper{}
	paramsKeeper := params.NewKeeper(cdc, paramskey, paramsTkey)

	keeper := NewKeeper(cdc, assetKey, bankKeeper, guardianKeeper, DefaultCodespace, paramsKeeper.Subspace(DefaultParamSpace))

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
	require.False(t, keeper.HasGateway(ctx, moniker))

	// create a gateway and assert that the gateway exists now
	keeper.SetGateway(ctx, gateway)
	require.True(t, keeper.HasGateway(ctx, moniker))

	// assert GetGateway will return the previous gateway
	res, _ := keeper.GetGateway(ctx, moniker)
	require.Equal(t, gateway, res)
}

func TestQueryGatewayKeeper(t *testing.T) {
	ms, assetKey, paramskey, paramsTkey := setupMultiStore()

	cdc := codec.New()
	RegisterCodec(cdc)
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	bankKeeper := bank.BaseKeeper{}
	guardianKeeper := guardian.Keeper{}
	paramsKeeper := params.NewKeeper(cdc, paramskey, paramsTkey)

	keeper := NewKeeper(cdc, assetKey, bankKeeper, guardianKeeper, DefaultCodespace, paramsKeeper.Subspace(DefaultParamSpace))

	// define variables
	var (
		owners     = []sdk.AccAddress{sdk.AccAddress([]byte("owner1")), sdk.AccAddress([]byte("owner2"))}
		monikers   = []string{"moni", "ker"}
		identities = []string{"id1", "id2"}
		details    = []string{"details1", "details2"}
		websites   = []string{"website1", "website2"}
	)

	// construct gateways
	gateway1 := Gateway{
		Owner:    owners[0],
		Moniker:  monikers[0],
		Identity: identities[0],
		Details:  details[0],
		Website:  websites[0],
	}

	gateway2 := Gateway{
		Owner:    owners[1],
		Moniker:  monikers[1],
		Identity: identities[1],
		Details:  details[1],
		Website:  websites[1],
	}

	// create gateways
	keeper.SetGateway(ctx, gateway1)
	keeper.SetOwnerGateway(ctx, gateway1.Owner, gateway1.Moniker)

	keeper.SetGateway(ctx, gateway2)
	keeper.SetOwnerGateway(ctx, gateway2.Owner, gateway2.Moniker)

	// query gateway
	res1, _ := keeper.GetGateway(ctx, gateway1.Moniker)
	require.Equal(t, gateway1, res1)

	res2, _ := keeper.GetGateway(ctx, gateway2.Moniker)
	require.Equal(t, gateway2, res2)

	// query gateways with a specified owner
	var gateways1 []Gateway
	iter1 := keeper.GetGateways(ctx, gateway1.Owner)
	defer iter1.Close()

	for ; iter1.Valid(); iter1.Next() {
		var moniker string
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iter1.Value(), &moniker)

		gateway, err := keeper.GetGateway(ctx, moniker)
		if err != nil {
			continue
		}

		gateways1 = append(gateways1, gateway)
	}

	require.Equal(t, []Gateway{gateway1}, gateways1)

	var gateways2 []Gateway
	iter2 := keeper.GetGateways(ctx, gateway2.Owner)
	defer iter2.Close()

	for ; iter2.Valid(); iter2.Next() {
		var moniker string
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iter2.Value(), &moniker)

		gateway, err := keeper.GetGateway(ctx, moniker)
		if err != nil {
			continue
		}

		gateways2 = append(gateways2, gateway)
	}

	require.Equal(t, []Gateway{gateway2}, gateways2)

	// query all gateways
	var gateways3 []Gateway
	iter3 := keeper.GetAllGateways(ctx)
	defer iter3.Close()

	for ; iter3.Valid(); iter3.Next() {
		var gateway Gateway
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iter3.Value(), &gateway)
		gateways3 = append(gateways3, gateway)
	}

	require.Equal(t, []Gateway{gateway2, gateway1}, gateways3)
}
