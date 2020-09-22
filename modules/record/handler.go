package record

import (
	"encoding/hex"

	"github.com/tendermint/tendermint/crypto/tmhash"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/record/keeper"
	"github.com/irisnet/irismod/modules/record/types"
)

// NewHandler returns a handler for all "record" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgCreateRecord:
			return handleMsgCreateRecord(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized record message type: %T", msg)
		}
	}
}

// handleMsgCreateRecord handles MsgCreateRecord
func handleMsgCreateRecord(ctx sdk.Context, k keeper.Keeper, msg *types.MsgCreateRecord) (*sdk.Result, error) {
	record := types.NewRecord(tmhash.Sum(ctx.TxBytes()), msg.Contents, msg.Creator)
	recordId := k.AddRecord(ctx, record)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator.String()),
		),
		sdk.NewEvent(
			types.EventTypeCreateRecord,
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator.String()),
			sdk.NewAttribute(types.AttributeKeyRecordID, hex.EncodeToString(recordId)),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
