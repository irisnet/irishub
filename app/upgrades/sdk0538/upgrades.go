package sdk0538

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/irisnet/irishub/v5/app/upgrades"
)

// Upgrade removes the obsolete IBC capability store before running the SDK and
// IBC module migrations at the coordinated upgrade height.
var Upgrade = upgrades.Upgrade{
	UpgradeName:   "cosmos-sdk-v0.53.8",
	StoreUpgrades: &storetypes.StoreUpgrades{Deleted: []string{"capability"}},
	UpgradeHandlerConstructor: func(m *module.Manager, c module.Configurator, box upgrades.Toolbox) upgradetypes.UpgradeHandler {
		return func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			// SetModuleVersionMap only overwrites present modules; explicitly remove
			// the retired module from the persisted version map as well.
			store := sdk.UnwrapSDKContext(ctx).KVStore(box.GetKey(upgradetypes.StoreKey))
			store.Delete(append([]byte{upgradetypes.VersionMapByte}, []byte("capability")...))
			delete(fromVM, "capability")
			return m.RunMigrations(ctx, c, fromVM)
		}
	},
}
