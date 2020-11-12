package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/random/types"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the random MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (m msgServer) RequestRandom(goCtx context.Context, msg *types.MsgRequestRandom) (*types.MsgRequestRandomResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	consumer, _ := sdk.AccAddressFromBech32(msg.Consumer)
	request, err := m.Keeper.RequestRandom(ctx, consumer, msg.BlockInterval, msg.Oracle, msg.ServiceFeeCap)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeRequestRandom,
				sdk.NewAttribute(types.AttributeKeyRequestID, hex.EncodeToString(types.GenerateRequestID(request))),
				sdk.NewAttribute(types.AttributeKeyConsumer, msg.Consumer),
				sdk.NewAttribute(types.AttributeKeyGenHeight, fmt.Sprintf("%d", request.Height+int64(msg.BlockInterval))),
				sdk.NewAttribute(types.AttributeKeyOracle, strconv.FormatBool(msg.Oracle)),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Consumer),
			),
		},
	)

	return &types.MsgRequestRandomResponse{}, nil
}
