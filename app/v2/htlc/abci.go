package htlc

import (
	sdk "github.com/irisnet/irishub/types"
)

// EndBlocker handles block ending logic
func EndBlocker(ctx sdk.Context, keeper Keeper) (resTags sdk.Tags) {
	// TODO: check and reset state (Open => Expired)
	


	// TODO: alternative
	// check timeout =>
	// refund =>
	// delete HTLC from expire queue

	return nil
}
