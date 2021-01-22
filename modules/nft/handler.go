package nft

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/nft/keeper"
	"github.com/irisnet/irismod/modules/nft/types"
)

// NewHandler routes the messages to the handlers
func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgIssueDenom:
			m := msg.Normalize()
			res, err := msgServer.IssueDenom(sdk.WrapSDKContext(ctx), &m)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgMintNFT:
			m := msg.Normalize()
			res, err := msgServer.MintNFT(sdk.WrapSDKContext(ctx), &m)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgTransferNFT:
			m := msg.Normalize()
			res, err := msgServer.TransferNFT(sdk.WrapSDKContext(ctx), &m)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgEditNFT:
			m := msg.Normalize()
			res, err := msgServer.EditNFT(sdk.WrapSDKContext(ctx), &m)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgBurnNFT:
			m := msg.Normalize()
			res, err := msgServer.BurnNFT(sdk.WrapSDKContext(ctx), &m)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized nft message type: %T", msg)
		}
	}
}
