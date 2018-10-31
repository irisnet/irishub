package iservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// handle all "iservice" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSvcDef:
			return handleMsgSvcDef(ctx, k, msg)
		case MsgSvcBind:
			return handleMsgSvcBind(ctx, k, msg)
		default:
			return sdk.ErrTxDecode("invalid message parse in staking module").Result()
		}
	}
}
func handleMsgSvcDef(ctx sdk.Context, k Keeper, msg MsgSvcDef) sdk.Result {
	_, found := k.GetServiceDefinition(ctx, msg.ChainId, msg.Name)
	if found {
		return ErrSvcDefExists(k.Codespace(), msg.Name).Result()
	}
	k.AddServiceDefinition(ctx, msg.SvcDef)
	err := k.AddMethods(ctx, msg.SvcDef)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleMsgSvcBind(ctx sdk.Context, k Keeper, msg MsgSvcBind) sdk.Result {
	_, found := k.GetServiceBinding(ctx, msg.DefChainID, msg.DefName, msg.BindChainID, msg.Provider)
	if found {
		return ErrSvcBindingExists(k.Codespace(), msg.Provider).Result()
	}
	return sdk.Result{}
}
