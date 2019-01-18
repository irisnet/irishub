package upgrade

import (
	"fmt"
	sdk "github.com/irisnet/irishub/types"
	"strconv"
)

// do switch
func EndBlocker(ctx sdk.Context, uk Keeper) (tags sdk.Tags) {

	ctx = ctx.WithLogger(ctx.Logger().With("handler", "endBlock").With("module", "iris/distribution"))

	tags = sdk.NewTags()
	upgradeConfig, ok := uk.protocolKeeper.GetUpgradeConfig(ctx)
	if ok {
		validator,found := uk.sk.GetValidatorByConsAddr(ctx,(sdk.ConsAddress)(ctx.BlockHeader().ProposerAddress));
		if!found {
			panic(fmt.Sprint("Proposer is not a bonded validator whose consaddress is %s", (sdk.ConsAddress)(ctx.BlockHeader().ProposerAddress).String()))
		}

		if ctx.BlockHeader().Version.App == upgradeConfig.Protocol.Version {
			uk.SetSignal(ctx, upgradeConfig.Protocol.Version, validator.ConsAddress().String())

			ctx.Logger().Info(
				fmt.Sprintf("Validator ConsAddress %s has downloaded the latest software which protocolVersion is %d ",
					validator.GetOperator().String(), ctx, upgradeConfig.Protocol.Version))
		} else {
			uk.DeleteSignal(ctx, upgradeConfig.Protocol.Version, validator.ConsAddress().String())

			ctx.Logger().Info(
				fmt.Sprintf("Validator ConsAddress %s has restarted the old software which protocolVersion is %d ",
					validator.GetOperator().String(), ctx, upgradeConfig.Protocol.Version))
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
	}

	tags = tags.AppendTag(sdk.AppVersionTag, []byte(strconv.FormatUint(uk.protocolKeeper.GetCurrentVersion(ctx), 10)))

	return tags
}
