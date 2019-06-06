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
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	assetKey := sdk.NewKVStoreKey("assetKey")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(assetKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	return ms, assetKey
}

func TestKeeper_IssueAsset(t *testing.T) {
	ms, assetKey := setupMultiStore()

	cdc := codec.New()
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	keeper := NewKeeper(cdc, assetKey, bank.BaseKeeper{}, guardian.Keeper{}, DefaultCodespace, params.Subspace{})

	addr := sdk.AccAddress([]byte("addr1"))
	addr2 := sdk.AccAddress([]byte("addr2"))
	addr3 := sdk.AccAddress([]byte("addr3"))

	msg := NewMsgIssueAsset("a", "b", "c", "d", 100,
		100, 18, false, addr, []sdk.AccAddress{addr2, addr3})
	_, err := keeper.IssueAsset(ctx, msg)
	assert.NoError(t, err)

	assert.True(t, keeper.HasAsset(ctx, msg.Symbol))

	asset, found := keeper.getAsset(ctx, msg.Symbol)
	assert.True(t, found)

	assert.Equal(t, msg.Symbol, asset.Symbol)
	assert.Equal(t, msg.Owner, asset.Owner)

	msgJson, _ := json.Marshal(msg)
	assetJson, _ := json.Marshal(asset)
	assert.Equal(t, msgJson, assetJson)
}
