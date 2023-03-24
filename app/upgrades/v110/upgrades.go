package v110

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/irisnet/irishub/app/upgrades"
	"github.com/irisnet/irishub/app/upgrades/v110/htlc"
	"github.com/irisnet/irishub/app/upgrades/v110/service"
	htlctypes "github.com/irisnet/irismod/modules/htlc/types"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:               "v1.1",
	UpgradeHandlerConstructor: upgradeHandlerConstructor,
	StoreUpgrades:             &store.StoreUpgrades{},
}

func upgradeHandlerConstructor(m *module.Manager, c module.Configurator, app upgrades.AppKeepers) types.UpgradeHandler {
	return func(ctx sdk.Context, plan types.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// migrate htlc
		if err := htlc.Migrate(ctx, app.AppCodec, app.HTLCKeeper, app.BankKeeper, app.GetKey(htlctypes.StoreKey)); err != nil {
			panic(err)
		}
		// migrate service
		if err := service.Migrate(ctx, app.ServiceKeeper, app.BankKeeper); err != nil {
			panic(err)
		}

		return fromVM, nil
	}
}
