package keeper

import (
	"testing"

	"github.com/irisnet/irishub/store"
	sdk "github.com/irisnet/irishub/types"
	dbm "github.com/tendermint/tendermint/libs/db"
)

func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	htlcKey := sdk.NewKVStoreKey("htlckey")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(htlcKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return ms, htlcKey
}

func TestKeeper_CreateHTLC(t *testing.T) {
	// TODO
}
