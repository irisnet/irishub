package mint

import (
	"time"

	sdk "github.com/irisnet/irishub/types"
)

// Called every block, process inflation on the first block of every hour
func BeginBlocker(ctx sdk.Context, k Keeper) {

	blockTime := ctx.BlockHeader().Time
	minter := k.GetMinter(ctx)
	if blockTime.Sub(minter.InflationLastTime) < time.Hour { // only mint on the hour!
		return
	}

	params := k.GetParams(ctx)
	minter.InflationLastTime = blockTime
	minter, mintedCoin := minter.ProcessProvisions(params)
	k.bk.IncreaseLoosenToken(ctx, sdk.Coins{mintedCoin})
	k.fck.AddCollectedFees(ctx, sdk.Coins{mintedCoin})
	k.SetMinter(ctx, minter)
}
