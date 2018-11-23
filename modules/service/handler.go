package service

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/service/tags"
	"fmt"
)

// handle all "service" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSvcDef:
			return handleMsgSvcDef(ctx, k, msg)
		case MsgSvcBind:
			return handleMsgSvcBind(ctx, k, msg)
		case MsgSvcBindingUpdate:
			return handleMsgSvcBindUpdate(ctx, k, msg)
		case MsgSvcDisable:
			return handleMsgSvcDisable(ctx, k, msg)
		case MsgSvcEnable:
			return handleMsgSvcEnable(ctx, k, msg)
		case MsgSvcRefundDeposit:
			return handleMsgSvcRefundDeposit(ctx, k, msg)
		case MsgSvcRequest:
			return handleMsgSvcRequest(ctx, k, msg)
		case MsgSvcResponse:
			return handleMsgSvcResponse(ctx, k, msg)
		case MsgSvcRefundFees:
			return handleMsgSvcRefundFees(ctx, k, msg)
		case MsgSvcWithdrawFees:
			return handleMsgSvcWithdrawFees(ctx, k, msg)
		default:
			return sdk.ErrTxDecode("invalid message parse in service module").Result()
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
	svcBinding := NewSvcBinding(msg.DefChainID, msg.DefName, msg.BindChainID, msg.Provider, msg.BindingType,
		msg.Deposit, msg.Prices, msg.Level, true, 0)
	err, _ := k.AddServiceBinding(ctx, svcBinding)
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
	svcBinding := NewSvcBinding(msg.DefChainID, msg.DefName, msg.BindChainID, msg.Provider, msg.BindingType,
		msg.Deposit, msg.Prices, msg.Level, false, 0)
	err, _ := k.UpdateServiceBinding(ctx, svcBinding)
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

func handleMsgSvcDisable(ctx sdk.Context, k Keeper, msg MsgSvcDisable) sdk.Result {
	err, _ := k.Disable(ctx, msg.DefChainID, msg.DefName, msg.BindChainID, msg.Provider)
	if err != nil {
		return err.Result()
	}
	resTags := sdk.NewTags(
		tags.Action, tags.ActionSvcDisable,
	)
	return sdk.Result{
		Tags: resTags,
	}
}

func handleMsgSvcEnable(ctx sdk.Context, k Keeper, msg MsgSvcEnable) sdk.Result {
	err, _ := k.Enable(ctx, msg.DefChainID, msg.DefName, msg.BindChainID, msg.Provider, msg.Deposit)
	if err != nil {
		return err.Result()
	}
	resTags := sdk.NewTags(
		tags.Action, tags.ActionSvcEnable,
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
		tags.Action, tags.ActionSvcRefundDeposit,
	)
	return sdk.Result{
		Tags: resTags,
	}
}

func handleMsgSvcRequest(ctx sdk.Context, k Keeper, msg MsgSvcRequest) sdk.Result {
	_, bindingFound := k.GetServiceBinding(ctx, msg.DefChainID, msg.DefName, msg.BindChainID, msg.Provider)
	if !bindingFound {
		return ErrSvcBindingNotExists(k.Codespace()).Result()
	}

	_, methodFound := k.GetMethod(ctx, msg.DefChainID, msg.DefName, msg.MethodID)
	if !methodFound {
		return ErrMethodNotExists(k.Codespace(), msg.MethodID).Result()
	}

	request := NewSvcRequest(msg.DefChainID, msg.DefName, msg.BindChainID, msg.ReqChainID, msg.Consumer, msg.Provider, msg.MethodID, msg.Input, msg.ServiceFee, msg.Profiling)

	k.AddRequest(ctx, request)
	k.AddActiveRequest(ctx, request)
	k.AddRequestExpiration(ctx, request)
	resTags := sdk.NewTags(
		tags.Action, tags.ActionSvcCall,
	)
	return sdk.Result{
		Tags: resTags,
	}
}

func handleMsgSvcResponse(ctx sdk.Context, k Keeper, msg MsgSvcResponse) sdk.Result {
	request, found := k.GetActiveRequest(ctx, msg.RequestHeight, msg.RequestIntraTxCounter)
	if !found {
		return ErrRequestNotActive(k.Codespace()).Result()
	}

	response := NewSvcResponse(msg.ReqChainID, msg.RequestHeight, msg.RequestIntraTxCounter, msg.Provider,
		request.Consumer, msg.Output, msg.ErrorMsg)

	k.AddResponse(ctx, response)

	// delete request from active request list and expiration list
	k.DeleteActiveRequest(ctx, request)
	k.DeleteRequestExpiration(ctx, request)

	k.AddReturnFee(ctx, response.Provider, request.ServiceFee)

	resTags := sdk.NewTags(
		tags.Action, tags.ActionSvcRespond,
	)
	return sdk.Result{
		Tags: resTags,
	}
}

func handleMsgSvcRefundFees(ctx sdk.Context, k Keeper, msg MsgSvcRefundFees) sdk.Result {
	k.RefundFee(ctx, msg.Consumer)
	resTags := sdk.NewTags(
		tags.Action, tags.ActionSvcRefundFees,
	)
	return sdk.Result{
		Tags: resTags,
	}
}

func handleMsgSvcWithdrawFees(ctx sdk.Context, k Keeper, msg MsgSvcWithdrawFees) sdk.Result {
	k.WithdrawFee(ctx, msg.Provider)
	resTags := sdk.NewTags(
		tags.Action, tags.ActionSvcWithdrawFees,
	)
	return sdk.Result{
		Tags: resTags,
	}
}

// Called every block, process inflation, update validator set
func EndBlocker(ctx sdk.Context, keeper Keeper) (resTags sdk.Tags) {

	logger := ctx.Logger().With("module", "service")
	resTags = sdk.NewTags()

	activeIterator := keeper.ActiveRequestQueueIterator(ctx, ctx.BlockHeight())
	for ; activeIterator.Valid(); activeIterator.Next() {
		var req SvcRequest
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(activeIterator.Value(), &req)
		keeper.AddReturnFee(ctx, req.Consumer, req.ServiceFee)
		keeper.DeleteActiveRequest(ctx, req)
		keeper.DeleteRequestExpiration(ctx, req)

		resTags = resTags.AppendTag(tags.Action, tags.ActionSvcCallTimeOut)
		logger.Info(fmt.Sprintf("request %s from %s timeout",
			req.RequestID(), req.Consumer))
	}
	activeIterator.Close()
	return resTags
}
