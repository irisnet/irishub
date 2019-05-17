package types

import (
	"github.com/irisnet/irishub/codec"
	"fmt"
)

const (
	AppVersionTag = "app_version"
	MainStore     = "main"
)

var (
	UpgradeConfigKey     = []byte("upgrade_config")
	CurrentVersionKey    = []byte("current_version")
	LastFailedVersionKey = []byte("last_failed_version")
	cdc                  = codec.New()
)

type ProtocolDefinition struct {
	Version   uint64 `json:"version"`
	Software  string `json:"software"`
	Height    uint64 `json:"height"`
	Threshold Dec    `json:"threshold"`
}

type UpgradeConfig struct {
	ProposalID uint64
	Protocol   ProtocolDefinition
}

func (uc UpgradeConfig) String() string {
	return fmt.Sprintf("proposalID: %v, version: %v, software: %s, height: %v, threshold: %s",
		uc.ProposalID, uc.Protocol.Version, uc.Protocol.Software, uc.Protocol.Height, uc.Protocol.Threshold.String(),
	)
}

func NewProtocolDefinition(version uint64, software string, height uint64, threshold Dec) ProtocolDefinition {
	return ProtocolDefinition{
		version,
		software,
		height,
		threshold,
	}
}

func NewUpgradeConfig(proposalID uint64, protocol ProtocolDefinition) UpgradeConfig {
	return UpgradeConfig{
		proposalID,
		protocol,
	}
}

func DefaultUpgradeConfig(software string) UpgradeConfig {
	return UpgradeConfig{
		ProposalID: uint64(0),
		Protocol:   NewProtocolDefinition(uint64(0), software, uint64(1), NewDecWithPrec(9, 1)),
	}
}

type ProtocolKeeper struct {
	storeKey StoreKey
	cdc      *codec.Codec
}

func NewProtocolKeeper(key StoreKey) ProtocolKeeper {
	return ProtocolKeeper{key, cdc}
}

func (pk ProtocolKeeper) GetCurrentVersionByStore(store KVStore) uint64 {
	bz := store.Get(CurrentVersionKey)
	if bz == nil {
		return 0
	}
	var currentVersion uint64
	pk.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &currentVersion)
	return currentVersion
}

func (pk ProtocolKeeper) GetUpgradeConfigByStore(store KVStore) (upgradeConfig UpgradeConfig, found bool) {
	bz := store.Get(UpgradeConfigKey)
	if bz == nil {
		return upgradeConfig, false
	}
	pk.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &upgradeConfig)
	return upgradeConfig, true
}

func (pk ProtocolKeeper) GetCurrentVersion(ctx Context) uint64 {
	store := ctx.KVStore(pk.storeKey)
	bz := store.Get(CurrentVersionKey)
	if bz == nil {
		return 0
	}
	var currentVersion uint64
	pk.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &currentVersion)
	return currentVersion
}

func (pk ProtocolKeeper) SetCurrentVersion(ctx Context, currentVersion uint64) {
	store := ctx.KVStore(pk.storeKey)
	bz := pk.cdc.MustMarshalBinaryLengthPrefixed(currentVersion)
	store.Set(CurrentVersionKey, bz)
}

func (pk ProtocolKeeper) GetLastFailedVersion(ctx Context) uint64 {
	store := ctx.KVStore(pk.storeKey)
	bz := store.Get(LastFailedVersionKey)
	if bz == nil {
		return 0 // default value
	}
	var lastFailedVersion uint64
	pk.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &lastFailedVersion)
	return lastFailedVersion
}

func (pk ProtocolKeeper) SetLastFailedVersion(ctx Context, lastFailedVersion uint64) {
	store := ctx.KVStore(pk.storeKey)
	bz := pk.cdc.MustMarshalBinaryLengthPrefixed(lastFailedVersion)
	store.Set(LastFailedVersionKey, bz)
}

func (pk ProtocolKeeper) GetUpgradeConfig(ctx Context) (upgradeConfig UpgradeConfig, found bool) {
	store := ctx.KVStore(pk.storeKey)
	bz := store.Get(UpgradeConfigKey)
	if bz == nil {
		return upgradeConfig, false
	}
	pk.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &upgradeConfig)
	return upgradeConfig, true
}

func (pk ProtocolKeeper) SetUpgradeConfig(ctx Context, upgradeConfig UpgradeConfig) {
	store := ctx.KVStore(pk.storeKey)
	bz := pk.cdc.MustMarshalBinaryLengthPrefixed(upgradeConfig)
	store.Set(UpgradeConfigKey, bz)
}

func (pk ProtocolKeeper) ClearUpgradeConfig(ctx Context) {
	store := ctx.KVStore(pk.storeKey)
	store.Delete(UpgradeConfigKey)
}

func (pk ProtocolKeeper) IsValidVersion(ctx Context, version uint64) bool {
	currentVersion := pk.GetCurrentVersion(ctx)
	lastFailedVersion := pk.GetLastFailedVersion(ctx)
	return isValidVersion(currentVersion, lastFailedVersion, version)
}

func isValidVersion(currentVersion uint64, lastFailedVersion uint64, version uint64) bool {
	if currentVersion >= lastFailedVersion {
		return currentVersion+1 == version
	} else {
		return lastFailedVersion == version || lastFailedVersion+1 == version
	}
}
