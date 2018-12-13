package upgrade

import (
	protocol "github.com/irisnet/irishub/app/protocol/keeper"
	sdk "github.com/irisnet/irishub/types"
	"strconv"
	tmstate "github.com/tendermint/tendermint/state"
)

// do switch
func EndBlocker(ctx sdk.Context, keeper Keeper) (tags sdk.Tags) {
	tags = sdk.NewTags()
	upgradeConfig := keeper.pk.GetUpgradeConfig(ctx)

	emptyUpgradeConfig := protocol.UpgradeConfig{}
	if upgradeConfig != emptyUpgradeConfig {
		if ctx.BlockHeader().Version.App == upgradeConfig.Definition.Version {
			keeper.SetSignal(ctx, upgradeConfig.Definition.Version, (sdk.ConsAddress)(ctx.BlockHeader().ProposerAddress).String())
		}

		if uint64(ctx.BlockHeight())+1 >= upgradeConfig.Definition.Height {
			success := tally(ctx, upgradeConfig.Definition.Version, keeper)
			if success{
				tags = tags.AppendTag(tmstate.UpgradeTagKey, []byte("Please install the right protocol version from " + upgradeConfig.Definition.Software))
			}

			appVersion := NewVersion(upgradeConfig, success)
			keeper.AddNewVersion(ctx, appVersion)
			keeper.pk.ClearUpgradeConfig(ctx)
		}
	}

	tags = tags.AppendTag("CurrentVersion", []byte(strconv.FormatUint(keeper.pk.GetCurrentProtocolVersion(ctx),10)))

	return tags
}
