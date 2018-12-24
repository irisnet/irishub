package mint

import (
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/mint/tags"
)

// Called every block, process inflation on the first block of every hour
func BeginBlocker(ctx sdk.Context, k Keeper) sdk.Tags {

	blockTime := ctx.BlockHeader().Time
	blockHeight := ctx.BlockHeader().Height
	minter := k.GetMinter(ctx)
	if blockHeight <= 1 { // don't inflate token in the first block
		minter.LastUpdate = blockTime
		k.SetMinter(ctx, minter)
		return nil
	}

	params := k.GetParams(ctx)
	minter.AnnualProvisions = minter.NextAnnualProvisions(params)
	mintedCoin := minter.BlockProvision(params, blockTime)

	k.bk.IncreaseLoosenToken(ctx, sdk.Coins{mintedCoin})
	k.fck.AddCollectedFees(ctx, sdk.Coins{mintedCoin})

	lastInflationTime := minter.LastUpdate
	minter.LastUpdate = blockTime
	k.SetMinter(ctx, minter)

	return sdk.NewTags(
		tags.LastInflationTime, []byte(lastInflationTime.String()),
		tags.InflationTime, []byte(blockTime.String()),
		tags.MintCoin, []byte(mintedCoin.String()),
	)
}
