package upgrade
import (
sdk "github.com/cosmos/cosmos-sdk/types"
)

var(
	KVStoreKeyListKey = []byte("k/")
)

// Get Proposal from store by ProposalID

func (keeper Keeper) GetKVStoreKeylist(ctx sdk.Context) string {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KVStoreKeyListKey)
	if bz == nil {
		return " "
	}
	KVstoreKeylist := string(bz)
	return KVstoreKeylist
}

// run when do switch
func (keeper Keeper) SetKVStoreKeylist(ctx sdk.Context) {
	store := ctx.KVStore(keeper.storeKey)

	version :=keeper.GetCurrentVersion(ctx)

	storeSet := make(map[string]bool)
	for _,x :=range version.ModuleList{
		storeSet[x.Store] = true
	}

	var KVStoreKeyList string
	for key,_ :=range storeSet{
		if KVStoreKeyList == "" {
			KVStoreKeyList+=key
		}else{
			KVStoreKeyList+=(":"+key)
		}
	}

	bz := []byte(KVStoreKeyList)
	store.Set(KVStoreKeyListKey, bz)
}

//"main:acc:ibc:stake:slashing:fee:gov:upgrade"
func InitGenesis_commitID(ctx sdk.Context, k Keeper) {
	k.SetKVStoreKeylist(ctx)
}
