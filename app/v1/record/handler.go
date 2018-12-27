package record

import (
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/record/tags"
)

// Handle all "record" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSubmitRecord:
			return handleMsgSubmitFile(ctx, keeper, msg)
		default:
			errMsg := "Unrecognized record msg type"
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgSubmitFile(ctx sdk.Context, keeper Keeper, msg MsgSubmitRecord) sdk.Result {

	msg.Memo = "test for protocol1"
	keeper.AddRecord(ctx, msg)

	recordIDBytes := []byte(msg.RecordID)

	resTags := sdk.NewTags(
		tags.OwnerAddress, []byte(msg.OwnerAddress.String()),
		tags.RecordID, recordIDBytes,
	)

	return sdk.Result{
		Data: recordIDBytes,
		Tags: resTags,
	}
}
