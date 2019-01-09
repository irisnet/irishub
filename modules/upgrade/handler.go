package upgrade

import (
	sdk "github.com/irisnet/irishub/types"
	"strconv"
)

// do switch
func EndBlocker(ctx sdk.Context, uk Keeper) (tags sdk.Tags) {
	tags = sdk.NewTags()
	upgradeConfig,ok := uk.protocolKeeper.GetUpgradeConfig(ctx)
	if ok {
		if ctx.BlockHeader().Version.App == upgradeConfig.Protocol.Version {
			uk.SetSignal(ctx, upgradeConfig.Protocol.Version, (sdk.ConsAddress)(ctx.BlockHeader().ProposerAddress).String())
		} else {
			uk.DeleteSignal(ctx, upgradeConfig.Protocol.Version, (sdk.ConsAddress)(ctx.BlockHeader().ProposerAddress).String())
		}

		x := upgradeConfig.Threshold.String()
		_ = x

		if uint64(ctx.BlockHeight())+1 == upgradeConfig.Protocol.Height {
			success := tally(ctx, upgradeConfig.Protocol.Version, uk, upgradeConfig.Threshold)

			if success {
				uk.protocolKeeper.SetCurrentVersion(ctx, upgradeConfig.Protocol.Version)
			} else {
				uk.protocolKeeper.SetLastFailedVersion(ctx, upgradeConfig.Protocol.Version)
			}

			uk.AddNewVersionInfo(ctx, NewVersionInfo(upgradeConfig, success))
			uk.protocolKeeper.ClearUpgradeConfig(ctx)
		}
	}

	// TODO: const CurrentVersionTagKey CurrentSoftwareTagKey
	tags = tags.AppendTag(sdk.AppVersionTag, []byte(strconv.FormatUint(uk.protocolKeeper.GetCurrentVersion(ctx),10)))

	return tags
}
