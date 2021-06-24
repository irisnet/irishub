package farm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/farm/keeper"
	"github.com/irisnet/irismod/modules/farm/types"
)

// EndBlocker handles block beginning logic for farm
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	logger := k.Logger(ctx).With("handler", "endBlocker")
	k.IteratorExpiredPool(ctx, ctx.BlockHeight(), func(pool types.FarmPool) {
		logger.Info(
			"The farm pool has expired, refund to creator",
			"poolName", pool.Name,
			"endHeight", pool.EndHeight,
			"lastHeightDistrRewards", pool.LastHeightDistrRewards,
			"totalLpTokenLocked", pool.TotalLpTokenLocked,
			"creator", pool.Creator,
		)
		if _, err := k.Refund(ctx, pool); err != nil {
			logger.Error("The farm pool refund failed",
				"poolName", pool.Name,
				"creator", pool.Creator,
				"errMsg", err.Error(),
			)
		}
	})
}
