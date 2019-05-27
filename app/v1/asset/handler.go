package asset

import (
	sdk "github.com/irisnet/irishub/types"
)

// handle all "asset" type messages.
// TODO
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		//switch msg := msg.(type) {
		//case ...
		//default:
		//	return sdk.ErrTxDecode("invalid message parse in asset module").Result()
		//}

		return sdk.ErrTxDecode("invalid message parse in asset module").Result()
	}
}

// Called every block, update request status
// TODO
func EndBlocker(ctx sdk.Context, keeper Keeper) (resTags sdk.Tags) {
	ctx = ctx.WithLogger(ctx.Logger().With("handler", "endBlock").With("module", "iris/asset"))

	resTags = sdk.NewTags()

	return resTags
}
