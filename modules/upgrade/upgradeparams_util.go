package upgrade

import (
	"github.com/irisnet/irishub/modules/upgrade/params"
	sdk "github.com/irisnet/irishub/types"
)

func GetUpgradeThreshlod(ctx sdk.Context) sdk.Dec {
	upgradeparams.UpgradeParameter.LoadValue(ctx)
	return upgradeparams.UpgradeParameter.Value.Threshold
}