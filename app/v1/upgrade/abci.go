package upgrade

import (
	"fmt"
	"strconv"

	sdk "github.com/irisnet/irishub/types"
)

// do switch
func EndBlocker(ctx sdk.Context, uk Keeper) (tags sdk.Tags) {

	ctx = ctx.WithLogger(ctx.Logger().With("handler", "endBlock").With("module", "iris/upgrade"))

	tags = sdk.NewTags()
	upgradeConfig, ok := uk.protocolKeeper.GetUpgradeConfig(ctx)
	if ok {
		uk.metrics.SetVersion(upgradeConfig.Protocol.Version)

		validator, found := uk.sk.GetValidatorByConsAddr(ctx, (sdk.ConsAddress)(ctx.BlockHeader().ProposerAddress))
		if !found {
			panic(fmt.Sprintf("validator with consensus-address %s not found", (sdk.ConsAddress)(ctx.BlockHeader().ProposerAddress).String()))
		}

		if ctx.BlockHeader().Version.App == upgradeConfig.Protocol.Version {
			uk.SetSignal(ctx, upgradeConfig.Protocol.Version, validator.ConsAddress().String())
			uk.metrics.SetSignal(validator.GetOperator().String(), upgradeConfig.Protocol.Version)

			ctx.Logger().Info("Validator has downloaded the latest software ",
				"validator", validator.GetOperator().String(), "version", upgradeConfig.Protocol.Version)

		} else {

			ok := uk.DeleteSignal(ctx, upgradeConfig.Protocol.Version, validator.ConsAddress().String())
			uk.metrics.DeleteSignal(validator.GetOperator().String(), upgradeConfig.Protocol.Version)

			if ok {
				ctx.Logger().Info("Validator has restarted the old software ",
					"validator", validator.GetOperator().String(), "version", upgradeConfig.Protocol.Version)
			}
		}

		if uint64(ctx.BlockHeight())+1 == upgradeConfig.Protocol.Height {
			success := tally(ctx, upgradeConfig.Protocol.Version, uk, upgradeConfig.Protocol.Threshold)

			if success {
				ctx.Logger().Info("Software Upgrade is successful.", "version", upgradeConfig.Protocol.Version)
				uk.protocolKeeper.SetCurrentVersion(ctx, upgradeConfig.Protocol.Version)
			} else {
				ctx.Logger().Info("Software Upgrade is failure.", "version", upgradeConfig.Protocol.Version)
				uk.protocolKeeper.SetLastFailedVersion(ctx, upgradeConfig.Protocol.Version)
			}

			uk.AddNewVersionInfo(ctx, NewVersionInfo(upgradeConfig, success))
			uk.protocolKeeper.ClearUpgradeConfig(ctx)
		}
	} else {
		uk.metrics.DeleteVersion()
	}

	tags = tags.AppendTag(sdk.AppVersionTag, []byte(strconv.FormatUint(uk.protocolKeeper.GetCurrentVersion(ctx), 10)))

	return tags
}
