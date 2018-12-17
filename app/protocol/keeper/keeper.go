package keeper

import (
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

var (
	CurrentProtocolVersionKey = []byte("currentProtocolVersionKey")
	UpgradeConfigkey          = []byte("upgradeConfigkey")
	LastFailureVersionKey     = []byte("lastFailureVersionKey")
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
}

func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		storeKey: storeKey,
		cdc:      cdc,
	}
}

func (keeper Keeper) GetCurrentProtocolVersion(ctx sdk.Context) uint64 {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(CurrentProtocolVersionKey)
	if bz == nil {
		return 0
	}
	var currentProtocolVersion uint64
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &currentProtocolVersion)
	return currentProtocolVersion
}

func (keeper Keeper) GetCurrentProtocolVersionByStore(kvStore sdk.KVStore) uint64 {
	bz := kvStore.Get(CurrentProtocolVersionKey)
	if bz == nil {
		return 0
	}
	var currentProtocolVersion uint64
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &currentProtocolVersion)
	return currentProtocolVersion
}

func (keeper Keeper) SetCurrentProtocolVersion(ctx sdk.Context, currentProtocolVersion uint64) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(currentProtocolVersion)
	store.Set(CurrentProtocolVersionKey, bz)
}

func (keeper Keeper) GetLastFailureVersion(ctx sdk.Context) uint64 {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(LastFailureVersionKey)
	if bz == nil {
		return 0
	}
	var lastFailureVersion uint64
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &lastFailureVersion)
	return lastFailureVersion
}

func (keeper Keeper) SetLastFailureVersion(ctx sdk.Context, lastFailureVersion uint64) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(lastFailureVersion)
	store.Set(LastFailureVersionKey, bz)
}

func (keeper Keeper) GetUpgradeConfig(ctx sdk.Context) (upgradeConfig UpgradeConfig, found bool) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(UpgradeConfigkey)
	if bz == nil {
		return upgradeConfig, false
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &upgradeConfig)
	return upgradeConfig, true
}

func (keeper Keeper) GetUpgradeConfigByStore(kvStore sdk.KVStore) (upgradeConfig UpgradeConfig, found bool) {
	bz := kvStore.Get(UpgradeConfigkey)
	if bz == nil {
		return upgradeConfig, false
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &upgradeConfig)
	return upgradeConfig, true
}

func (keeper Keeper) SetUpgradeConfig(ctx sdk.Context, upgradeConfig UpgradeConfig) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(upgradeConfig)
	store.Set(UpgradeConfigkey, bz)
}

func (keeper Keeper) ClearUpgradeConfig(ctx sdk.Context) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(UpgradeConfigkey)
}

func (Keeper Keeper) IsValidProtocolVersion(ctx sdk.Context, protocolVersion uint64) bool {

	currentProtocolVersion := Keeper.GetCurrentProtocolVersion(ctx)
	lastFailureVersion := Keeper.GetLastFailureVersion(ctx)

	return isValidProtocolVersion(currentProtocolVersion, lastFailureVersion, protocolVersion)
}

func isValidProtocolVersion(currentProtocolVersion uint64, lastFailureVersion uint64, protocolVersion uint64) bool {
	if currentProtocolVersion >= lastFailureVersion {
		if currentProtocolVersion+1 == protocolVersion {
			return true
		} else {
			return false
		}
	} else {
		if lastFailureVersion == protocolVersion || lastFailureVersion+1 == protocolVersion {
			return true
		} else {
			return false
		}
	}
}
