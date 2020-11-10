package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irismod/modules/random/types"
)

type msgServer struct {
	Keeper
}


var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (m msgServer) RequestRandom(goCtx context.Context, msg *types.MsgRequestRandom) (*types.MsgRequestRandomResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	request, err := m.Keeper.RequestRandom(ctx, msg.Consumer, msg.BlockInterval, msg.Oracle, msg.ServiceFeeCap)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Consumer.String()),
			),
			sdk.NewEvent(
				types.EventTypeRequestRandom,
				sdk.NewAttribute(types.AttributeKeyRequestID, hex.EncodeToString(types.GenerateRequestID(request))),
				sdk.NewAttribute(types.AttributeKeyGenHeight, fmt.Sprintf("%d", request.Height+int64(msg.BlockInterval))),
			),
		},
	)
	return &types.MsgRequestRandomResponse{},nil
}
