package types

import (
	"github.com/irisnet/irishub/codec"
)

const (
	AppVersionTag = "app_version"
	MainStore     = "main"
)

var (
	KeyMain              = NewKVStoreKey(MainStore)
	UpgradeConfigKey     = []byte("upgrade_config")
	CurrentVersionKey    = []byte("current_version")
	LastFailedVersionKey = []byte("last_failed_version")
)

type ProtocolDefinition struct {
	Version 	uint64	`json:"version"`
	Software	string	`json:"software"`
	Height		uint64	`json:"height"`
}

type UpgradeConfig struct {
	ProposalID uint64
	Protocol   ProtocolDefinition
}

func NewProtocolDefinition(version uint64, software string, height uint64) ProtocolDefinition {
	return ProtocolDefinition{
		version,
		software,
		height,
	}
}

func NewUpgradeConfig(proposalID uint64, protocol ProtocolDefinition) UpgradeConfig {
	return UpgradeConfig {
		proposalID,
		protocol,
	}
}

type ProtocolKeeper struct {
	cdc *codec.Codec
}

func NewProtocolKeeper(cdc *codec.Codec) ProtocolKeeper {
	keeper := ProtocolKeeper{cdc }
	return keeper
}

func (pk ProtocolKeeper) GetCurrentVersionByStore(kvStore KVStore) uint64 {
	bz := kvStore.Get(CurrentVersionKey)
	if bz == nil {
		return 0
	}
	var currentVersion uint64
	pk.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &currentVersion)
	return currentVersion
}

func (pk ProtocolKeeper) GetUpgradeConfigByStore(kvStore KVStore) (upgradeConfig UpgradeConfig, found bool) {
	bz := kvStore.Get(UpgradeConfigKey)
	if bz == nil {
		return upgradeConfig, false
	}
	pk.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &upgradeConfig)
	return upgradeConfig, true
}

func (pk ProtocolKeeper) GetCurrentVersion(ctx Context) uint64 {
	store := ctx.KVStore(KeyMain)
	bz := store.Get(CurrentVersionKey)
	if bz == nil {
		return 0
	}
	var currentVersion uint64
	pk.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &currentVersion)
	return currentVersion
}

func (pk ProtocolKeeper) SetCurrentVersion(ctx Context, currentVersion uint64) {
	store := ctx.KVStore(KeyMain)
	bz := pk.cdc.MustMarshalBinaryLengthPrefixed(currentVersion)
	store.Set(CurrentVersionKey, bz)
}

func (pk ProtocolKeeper) GetLastFailedVersion(ctx Context) uint64 {
	store := ctx.KVStore(KeyMain)
	bz := store.Get(LastFailedVersionKey)
	if bz == nil {
		return 0
	}
	var lastFailedVersion uint64
	pk.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &lastFailedVersion)
	return lastFailedVersion
}

func (pk ProtocolKeeper) SetLastFailedVersion(ctx Context, lastFailedVersion uint64) {
	store := ctx.KVStore(KeyMain)
	bz := pk.cdc.MustMarshalBinaryLengthPrefixed(lastFailedVersion)
	store.Set(LastFailedVersionKey, bz)
}

func (pk ProtocolKeeper) GetUpgradeConfig(ctx Context) (upgradeConfig UpgradeConfig, found bool) {
	store := ctx.KVStore(KeyMain)
	bz := store.Get(UpgradeConfigKey)
	if bz == nil {
		return upgradeConfig, false
	}
	pk.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &upgradeConfig)
	return upgradeConfig, true
}

func (pk ProtocolKeeper) SetUpgradeConfig(ctx Context, upgradeConfig UpgradeConfig) {
	store := ctx.KVStore(KeyMain)
	bz := pk.cdc.MustMarshalBinaryLengthPrefixed(upgradeConfig)
	store.Set(UpgradeConfigKey, bz)
}

func (pk ProtocolKeeper) ClearUpgradeConfig(ctx Context) {
	store := ctx.KVStore(KeyMain)
	store.Delete(UpgradeConfigKey)
}

func (pk ProtocolKeeper) IsValidVersion(ctx Context, version uint64) bool {
	currentVersion := pk.GetCurrentVersion(ctx)
	lastFailedVersion := pk.GetLastFailedVersion(ctx)

	if currentVersion >= lastFailedVersion {
		return currentVersion+1 == version
	} else {
		return lastFailedVersion == version || lastFailedVersion+1 == version
	}
}

/*
func TestKeeper(t *testing.T) {
	require.Equal(t, false, isValidVersion(1, 1, 1))
	require.Equal(t, true, isValidVersion(1, 4, 4))
	require.Equal(t, true, isValidVersion(1, 4, 5))
	require.Equal(t, true, isValidVersion(1, 1, 2))
	require.Equal(t, true, isValidVersion(2, 1, 3))
}
*/