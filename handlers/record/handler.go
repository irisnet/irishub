package record

import (
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/handlers/record/tags"
	types "github.com/irisnet/irishub/types/record"
	record "github.com/irisnet/irishub/keepers/record"
)

// Handle all "record" type messages.
func NewHandler(keeper record.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgSubmitRecord:
			return handleMsgSubmitFile(ctx, keeper, msg)
		default:
			errMsg := "Unrecognized record msg type"
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgSubmitFile(ctx sdk.Context, keeper record.Keeper, msg types.MsgSubmitRecord) sdk.Result {

	keeper.AddRecord(ctx, msg)

	recordIDBytes := []byte(msg.RecordID)

	resTags := sdk.NewTags(
		tags.Action, tags.ActionSubmitRecord,
		tags.OwnerAddress, []byte(msg.OwnerAddress.String()),
		tags.RecordID, recordIDBytes,
	)

	return sdk.Result{
		Data: recordIDBytes,
		Tags: resTags,
	}
}
