package record

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Handle all "record" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSubmitFile:
			return handleMsgSubmitFile(ctx, keeper, msg)
		default:
			errMsg := "Unrecognized record msg type"
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgSubmitFile(ctx sdk.Context, keeper Keeper, msg MsgSubmitFile) sdk.Result {

	keeper.AddRecord(ctx, msg)

	return sdk.Result{
		Log: msg.RecordID,
	}
}
