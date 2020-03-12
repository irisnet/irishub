package keeper

import (
	dbm "github.com/tendermint/tm-db"

	"github.com/irisnet/irishub/store"
	sdk "github.com/irisnet/irishub/types"
)

func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey, *sdk.KVStoreKey, *sdk.KVStoreKey, *sdk.KVStoreKey, *sdk.TransientStoreKey) {
	db := dbm.NewMemDB()

	accountKey := sdk.NewKVStoreKey("accountKey")
	assetKey := sdk.NewKVStoreKey("assetKey")
	guardianKey := sdk.NewKVStoreKey("guardianKey")
	paramskey := sdk.NewKVStoreKey("params")
	paramsTkey := sdk.NewTransientStoreKey("transient_params")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(accountKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(assetKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(paramskey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(paramsTkey, sdk.StoreTypeIAVL, db)
	_ = ms.LoadLatestVersion()

	return ms, accountKey, assetKey, guardianKey, paramskey, paramsTkey
}
