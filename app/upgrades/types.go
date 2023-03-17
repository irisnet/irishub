package upgrades

import (
	"github.com/cosmos/cosmos-sdk/codec"
	store "github.com/cosmos/cosmos-sdk/store/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	htlckeeper "github.com/irisnet/irismod/modules/htlc/keeper"
	servicekeeper "github.com/irisnet/irismod/modules/service/keeper"

	tibckeeper "github.com/bianjieai/tibc-go/modules/tibc/core/keeper"
)

// Upgrade defines a struct containing necessary fields that a SoftwareUpgradeProposal
// must have written, in order for the state migration to go smoothly.
// An upgrade must implement this struct, and then set it in the app.go.
// The app.go will then define the handler.
type Upgrade struct {
	// Upgrade version name, for the upgrade handler, e.g. `v7`
	UpgradeName string

	// UpgradeHandlerConstructor defines the function that creates an upgrade handler
	UpgradeHandlerConstructor func(*module.Manager, module.Configurator, AppKeepers) upgradetypes.UpgradeHandler

	// Store upgrades, should be used for any new modules introduced, new modules deleted, or store names renamed.
	StoreUpgrades store.StoreUpgrades
}

type AppKeepers interface {
	AppCodec() codec.Codec
	HTLCKeeper() htlckeeper.Keeper
	BankKeeper() bankkeeper.Keeper
	ServiceKeeper() servicekeeper.Keeper
	Keys(moduleName string) *storetypes.KVStoreKey
	ModuleManager() *module.Manager
	TIBCkeeper() *tibckeeper.Keeper
}
