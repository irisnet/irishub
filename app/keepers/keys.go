package keepers

import (
	storetypes "cosmossdk.io/store/types"
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/feegrant"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"

	icahosttypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"

	evmtypes "github.com/evmos/ethermint/x/evm/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"

	coinswaptypes "mods.irisnet.org/modules/coinswap/types"
	farmtypes "mods.irisnet.org/modules/farm/types"
	htlctypes "mods.irisnet.org/modules/htlc/types"
	mttypes "mods.irisnet.org/modules/mt/types"
	nfttypes "mods.irisnet.org/modules/nft/types"
	oracletypes "mods.irisnet.org/modules/oracle/types"
	randomtypes "mods.irisnet.org/modules/random/types"
	recordtypes "mods.irisnet.org/modules/record/types"
	servicetypes "mods.irisnet.org/modules/service/types"
	tokentypes "mods.irisnet.org/modules/token/types"

	ibcnfttransfertypes "github.com/bianjieai/nft-transfer/types"
	tibcmttypes "github.com/bianjieai/tibc-go/modules/tibc/apps/mt_transfer/types"
	tibcnfttypes "github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/types"
	tibchost "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"

	guardiantypes "github.com/irisnet/irishub/v4/modules/guardian/types"
	minttypes "github.com/irisnet/irishub/v4/modules/mint/types"
)

func (appKeepers *AppKeepers) genStoreKeys() {
	// Define what keys will be used in the cosmos-sdk key/value store.
	// Cosmos-SDK modules each have a "key" that allows the application to reference what they've stored on the chain.
	appKeepers.keys = storetypes.NewKVStoreKeys(
		authtypes.StoreKey,
		banktypes.StoreKey,
		stakingtypes.StoreKey,
		minttypes.StoreKey,
		distrtypes.StoreKey,
		slashingtypes.StoreKey,
		crisistypes.StoreKey,
		govtypes.StoreKey,
		paramstypes.StoreKey,
		ibcexported.StoreKey,
		upgradetypes.StoreKey,
		consensustypes.StoreKey,
		evidencetypes.StoreKey,
		ibctransfertypes.StoreKey,
		ibcnfttransfertypes.StoreKey,
		icahosttypes.StoreKey,
		capabilitytypes.StoreKey,
		guardiantypes.StoreKey,
		tokentypes.StoreKey,
		nfttypes.StoreKey,
		htlctypes.StoreKey,
		recordtypes.StoreKey,
		coinswaptypes.StoreKey,
		servicetypes.StoreKey,
		oracletypes.StoreKey,
		randomtypes.StoreKey,
		farmtypes.StoreKey,
		feegrant.StoreKey,
		tibchost.StoreKey,
		tibcnfttypes.StoreKey,
		tibcmttypes.StoreKey,
		mttypes.StoreKey,
		authzkeeper.StoreKey,
		// ethermint keys
		evmtypes.StoreKey,
		feemarkettypes.StoreKey,
	)

	// Define transient store keys
	appKeepers.tkeys = storetypes.NewTransientStoreKeys(
		paramstypes.TStoreKey,
		evmtypes.TransientKey,
		feemarkettypes.TransientKey,
	)

	// MemKeys are for information that is stored only in RAM.
	appKeepers.memKeys = storetypes.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

}

// KvStoreKeys returns the map of string to KVStoreKey.
//
// None.
// map[string]*storetypes.KVStoreKey.
func (appKeepers *AppKeepers) KvStoreKeys() map[string]*storetypes.KVStoreKey {
	return appKeepers.keys
}

// TransientStoreKeys returns the map of transient store keys.
//
// No parameters.
// Returns a map[string]*storetypes.TransientStoreKey.
func (appKeepers *AppKeepers) TransientStoreKeys() map[string]*storetypes.TransientStoreKey {
	return appKeepers.tkeys
}

// MemoryStoreKeys returns the map of type map[string]*storetypes.MemoryStoreKey.
//
// No parameters.
// Returns a map of type map[string]*storetypes.MemoryStoreKey.
func (appKeepers *AppKeepers) MemoryStoreKeys() map[string]*storetypes.MemoryStoreKey {
	return appKeepers.memKeys
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (appKeepers *AppKeepers) GetKey(storeKey string) *storetypes.KVStoreKey {
	return appKeepers.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (appKeepers *AppKeepers) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return appKeepers.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (appKeepers *AppKeepers) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return appKeepers.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (appKeepers *AppKeepers) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := appKeepers.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}
