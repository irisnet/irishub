package keeper

import (
	"context"
	"encoding/hex"

	"github.com/cometbft/cometbft/crypto/tmhash"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"mods.irisnet.org/modules/record/types"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the record MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (m msgServer) CreateRecord(
	goCtx context.Context,
	msg *types.MsgCreateRecord,
) (*types.MsgCreateRecordResponse, error) {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	record := types.NewRecord(tmhash.Sum(ctx.TxBytes()), msg.Contents, creator)
	recordId := m.Keeper.AddRecord(ctx, record)

	hexID := hex.EncodeToString(recordId)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateRecord,
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyRecordID, hex.EncodeToString(recordId)),
		),
	})

	return &types.MsgCreateRecordResponse{
		Id: hexID,
	}, nil
}
