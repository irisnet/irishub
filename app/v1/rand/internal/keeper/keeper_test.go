package keeper

import (
	"testing"

	"github.com/irisnet/irishub/store"
	sdk "github.com/irisnet/irishub/types"
	dbm "github.com/tendermint/tendermint/libs/db"
)

func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey, *sdk.KVStoreKey, *sdk.KVStoreKey, *sdk.TransientStoreKey) {
	db := dbm.NewMemDB()
	accountKey := sdk.NewKVStoreKey("accountKey")
	assetKey := sdk.NewKVStoreKey("assetKey")
	paramskey := sdk.NewKVStoreKey("params")
	paramsTkey := sdk.NewTransientStoreKey("transient_params")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(accountKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(assetKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(paramskey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(paramsTkey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	return ms, accountKey, assetKey, paramskey, paramsTkey
}

func TestRandKeeper(t *testing.T) {
	// ms, accountKey, assetKey, paramskey, paramsTkey := setupMultiStore()

	// TODO
}
