package asset

import (
	"encoding/json"
	"testing"

	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v1/bank"
	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/modules/guardian"
	"github.com/irisnet/irishub/store"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/assert"
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

	cdc := msgCdc
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	pk := params.NewKeeper(cdc, paramskey, paramsTkey)
	keeper := NewKeeper(cdc, assetKey, bank.BaseKeeper{}, guardian.Keeper{}, DefaultCodespace, pk.Subspace(DefaultParamSpace))

	addr := sdk.AccAddress([]byte("addr1"))

	msg := NewMsgIssueAsset(BaseAsset{0x00, 0x00, "c", "d", "e", 1, "f", 1, 1, true, addr})
	_, err := keeper.IssueAsset(ctx, msg)
	assert.NoError(t, err)

	assert.True(t, keeper.HasAsset(ctx, msg.GetDenom()))

	asset, found := keeper.getAsset(ctx, msg.GetDenom())
	assert.True(t, found)

	assert.Equal(t, msg.GetDenom(), asset.GetDenom())
	assert.Equal(t, msg.Asset.(BaseAsset).Owner, msg.Asset.(BaseAsset).Owner)

	msgJson, _ := json.Marshal(msg.Asset)
	assetJson, _ := json.Marshal(asset)
	assert.Equal(t, msgJson, assetJson)
}
