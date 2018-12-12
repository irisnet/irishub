package keeper
import (
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/codec"
)


var (
	CurrentProtocolVersionKey = []byte("currentProtocolVersionKey")
	UpgradeConfigkey  = []byte("upgradeConfigkey")
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
}

func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		storeKey:storeKey,
		cdc:cdc,
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

func (keeper Keeper) GetUpgradeConfig(ctx sdk.Context) UpgradeConfig {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(UpgradeConfigkey)
	if bz == nil {
		return UpgradeConfig{}
	}
	var upgradeConfig UpgradeConfig
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &upgradeConfig)
	return upgradeConfig
}

func (keeper Keeper) SetUpgradeConfig(ctx sdk.Context, upgradeConfig UpgradeConfig) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(upgradeConfig)
	store.Set(UpgradeConfigkey, bz)
}

func (keeper Keeper) ClearUpgradeConfig(ctx sdk.Context) {
	store := ctx.KVStore(keeper.storeKey)
	upgradeConfig := UpgradeConfig{}
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(upgradeConfig)
	store.Set(UpgradeConfigkey, bz)
}