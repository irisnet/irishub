package upgrade

import (
	sdk "github.com/irisnet/irishub/types"
	"strconv"
	protocolKeeper "github.com/irisnet/irishub/app/protocol/keeper"
)

// do switch
func EndBlocker(ctx sdk.Context, keeper Keeper) (tags sdk.Tags) {
	tags = sdk.NewTags()
	upgradeConfig,ok := keeper.pk.GetUpgradeConfig(ctx)
	if ok {
		if ctx.BlockHeader().Version.App == upgradeConfig.Definition.Version {
			keeper.SetSignal(ctx, upgradeConfig.Definition.Version, (sdk.ConsAddress)(ctx.BlockHeader().ProposerAddress).String())
		} else {
			keeper.DeleteSignal(ctx, upgradeConfig.Definition.Version, (sdk.ConsAddress)(ctx.BlockHeader().ProposerAddress).String())
		}

		if uint64(ctx.BlockHeight())+1 == upgradeConfig.Definition.Height {
			success := tally(ctx, upgradeConfig.Definition.Version, keeper)

			if success {
				keeper.pk.SetCurrentProtocolVersion(ctx, upgradeConfig.Definition.Version)
			} else {
				keeper.pk.SetLastFailureVersion(ctx, upgradeConfig.Definition.Version)
			}

			appVersion := NewVersion(upgradeConfig, success)
			keeper.AddNewVersion(ctx, appVersion)
			keeper.pk.ClearUpgradeConfig(ctx)
		}
	}

	// TODO: const CurrentVersionTagKey CurrentSoftwareTagKey
	tags = tags.AppendTag(protocolKeeper.AppVersionTag, []byte(strconv.FormatUint(keeper.pk.GetCurrentProtocolVersion(ctx),10)))

	return tags
}
