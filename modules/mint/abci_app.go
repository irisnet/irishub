package mint

import (
	"fmt"

	"github.com/irisnet/irishub/modules/mint/tags"
	sdk "github.com/irisnet/irishub/types"
)

// Called every block, process inflation on the first block of every hour
func BeginBlocker(ctx sdk.Context, k Keeper) sdk.Tags {
	logger := ctx.Logger().With("handler", "beginBlock").With("module", fmt.Sprintf("iris/%s", "mint"))
	// Get block BFT time and block height
	blockTime := ctx.BlockHeader().Time
	blockHeight := ctx.BlockHeader().Height
	minter := k.GetMinter(ctx)
	if blockHeight <= 1 { // don't inflate token in the first block
		minter.LastUpdate = blockTime
		k.SetMinter(ctx, minter)
		return nil
	}

	// Calculate block mint amount
	params := k.GetParamSet(ctx)
	annualProvisions := minter.NextAnnualProvisions(params)
	mintedCoin := minter.BlockProvision(annualProvisions)

	// Increase loosen token and add minted coin to feeCollector
	k.bk.IncreaseLoosenToken(ctx, sdk.Coins{mintedCoin})
	k.fk.AddCollectedFees(ctx, sdk.Coins{mintedCoin})

	// Update last block BFT time
	lastInflationTime := minter.LastUpdate
	minter.LastUpdate = blockTime
	k.SetMinter(ctx, minter)

	logger.Info(fmt.Sprintf("Minted coin %s at height %d", mintedCoin, blockHeight), "time", blockTime)
	// Add tags
	return sdk.NewTags(
		tags.LastInflationTime, []byte(lastInflationTime.String()),
		tags.InflationTime, []byte(blockTime.String()),
		tags.MintCoin, []byte(mintedCoin.String()),
	)
}
