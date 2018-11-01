package iservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/iservice/tags"
)

// handle all "iservice" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSvcDef:
			return handleMsgSvcDef(ctx, k, msg)
		case MsgSvcBind:
			return handleMsgSvcBind(ctx, k, msg)
		case MsgSvcBindingUpdate:
			return handleMsgSvcBindUpdate(ctx, k, msg)
		case MsgSvcRefundDeposit:
			return handleMsgSvcRefundDeposit(ctx, k, msg)
		default:
			return sdk.ErrTxDecode("invalid message parse in iservice module").Result()
		}
	}
}
func handleMsgSvcDef(ctx sdk.Context, k Keeper, msg MsgSvcDef) sdk.Result {
	_, found := k.GetServiceDefinition(ctx, msg.ChainId, msg.Name)
	if found {
		return ErrSvcDefExists(k.Codespace(), msg.ChainId, msg.Name).Result()
	}
	k.AddServiceDefinition(ctx, msg.SvcDef)
	err := k.AddMethods(ctx, msg.SvcDef)
	if err != nil {
		return err.Result()
	}
	resTags := sdk.NewTags(
		tags.Action, tags.ActionSvcDef,
	)
	return sdk.Result{
		Tags: resTags,
	}
}

func handleMsgSvcBind(ctx sdk.Context, k Keeper, msg MsgSvcBind) sdk.Result {
	err, _ := k.AddServiceBinding(ctx, msg.SvcBinding)
	if err != nil {
		return err.Result()
	}
	resTags := sdk.NewTags(
		tags.Action, tags.ActionSvcBind,
	)
	return sdk.Result{
		Tags: resTags,
	}
}

func handleMsgSvcBindUpdate(ctx sdk.Context, k Keeper, msg MsgSvcBindingUpdate) sdk.Result {
	err, _ := k.UpdateServiceBinding(ctx, msg.SvcBinding)
	if err != nil {
		return err.Result()
	}
	resTags := sdk.NewTags(
		tags.Action, tags.ActionSvcBindUpdate,
	)
	return sdk.Result{
		Tags: resTags,
	}
}

func handleMsgSvcRefundDeposit(ctx sdk.Context, k Keeper, msg MsgSvcRefundDeposit) sdk.Result {
	err, _ := k.RefundDeposit(ctx, msg.DefChainID, msg.DefName, msg.BindChainID, msg.Provider)
	if err != nil {
		return err.Result()
	}
	resTags := sdk.NewTags(
		tags.Action, tags.ActionSvcBindUpdate,
	)
	return sdk.Result{
		Tags: resTags,
	}
}
