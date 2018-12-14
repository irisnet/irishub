package service

import (
	sdk "github.com/irisnet/irishub/types"
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
		case MsgSvcWithdrawTax:
			return handleMsgSvcWithdrawTax(ctx, k, msg)
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
	svcBinding := NewSvcBinding(ctx, msg.DefChainID, msg.DefName, msg.BindChainID, msg.Provider, msg.BindingType,
		msg.Deposit, msg.Prices, msg.Level, true)
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
	svcBinding := NewSvcBinding(ctx, msg.DefChainID, msg.DefName, msg.BindChainID, msg.Provider, msg.BindingType,
		msg.Deposit, msg.Prices, msg.Level, false)
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
	bind, bindingFound := k.GetServiceBinding(ctx, msg.DefChainID, msg.DefName, msg.BindChainID, msg.Provider)
	if !bindingFound {
		return ErrSvcBindingNotExists(k.Codespace()).Result()
	}
	if !bind.Available {
		return ErrSvcBindingNotAvailable(k.Codespace()).Result()
	}

	_, methodFound := k.GetMethod(ctx, msg.DefChainID, msg.DefName, msg.MethodID)
	if !methodFound {
		return ErrMethodNotExists(k.Codespace(), msg.MethodID).Result()
	}

	if msg.Profiling {
		if _, found := k.gk.GetProfiler(ctx, msg.Consumer); !found {
			return ErrNotProfiler(k.Codespace(), msg.Consumer).Result()
		}
	}

	//Method id start at 1
	if len(bind.Prices) >= int(msg.MethodID) && !msg.ServiceFee.IsAllGTE(sdk.Coins{bind.Prices[msg.MethodID-1]}) {
		return ErrLtServiceFee(k.Codespace(), sdk.Coins{bind.Prices[msg.MethodID-1]}).Result()
	}

	request := NewSvcRequest(msg.DefChainID, msg.DefName, msg.BindChainID, msg.ReqChainID, msg.Consumer, msg.Provider, msg.MethodID, msg.Input, msg.ServiceFee, msg.Profiling)

	// request service fee is equal to service binding service fee if not profiling
	if len(bind.Prices) >= int(msg.MethodID) && !msg.Profiling {
		request.ServiceFee = sdk.Coins{bind.Prices[msg.MethodID-1]}
	} else {
		request.ServiceFee = nil
	}

	request, err := k.AddRequest(ctx, request)
	if err != nil {
		return err.Result()
	}
	resTags := sdk.NewTags(
		tags.Action, tags.ActionSvcCall,
		tags.RequestID, []byte(request.RequestID()),
		tags.Provider, []byte(request.Provider.String()),
		tags.Consumer, []byte(request.Consumer.String()),
	)
	return sdk.Result{
		Tags: resTags,
	}
}

func handleMsgSvcResponse(ctx sdk.Context, k Keeper, msg MsgSvcResponse) sdk.Result {
	eHeight, rHeight, counter, _ := ConvertRequestID(msg.RequestID)
	request, found := k.GetActiveRequest(ctx, eHeight, rHeight, counter)
	if !found {
		request.ExpirationHeight = eHeight
		request.RequestHeight = rHeight
		request.RequestIntraTxCounter = counter
		return ErrRequestNotActive(k.Codespace(), request.RequestID()).Result()
	}
	if !(msg.Provider.Equals(request.Provider)) {
		return ErrNotMatchingProvider(k.Codespace(), request.Provider).Result()
	}
	if request.ReqChainID != msg.ReqChainID {
		return ErrNotMatchingReqChainID(k.Codespace(), msg.ReqChainID).Result()
	}

	response := NewSvcResponse(msg.ReqChainID, eHeight, rHeight, counter, msg.Provider,
		request.Consumer, msg.Output, msg.ErrorMsg)

	k.AddResponse(ctx, response)

	// delete request from active request list and expiration list
	k.DeleteActiveRequest(ctx, request)
	k.DeleteRequestExpiration(ctx, request)

	err := k.AddIncomingFee(ctx, response.Provider, request.ServiceFee)
	if err != nil {
		return err.Result()
	}

	resTags := sdk.NewTags(
		tags.Action, tags.ActionSvcRespond,
		tags.RequestID, []byte(request.RequestID()),
		tags.Consumer, []byte(response.Consumer.String()),
		tags.Provider, []byte(response.Provider.String()),
	)
	return sdk.Result{
		Tags: resTags,
	}
}

func handleMsgSvcRefundFees(ctx sdk.Context, k Keeper, msg MsgSvcRefundFees) sdk.Result {
	err := k.RefundFee(ctx, msg.Consumer)
	if err != nil {
		return err.Result()
	}
	resTags := sdk.NewTags(
		tags.Action, tags.ActionSvcRefundFees,
	)
	return sdk.Result{
		Tags: resTags,
	}
}

func handleMsgSvcWithdrawFees(ctx sdk.Context, k Keeper, msg MsgSvcWithdrawFees) sdk.Result {
	err := k.WithdrawFee(ctx, msg.Provider)
	if err != nil {
		return err.Result()
	}
	resTags := sdk.NewTags(
		tags.Action, tags.ActionSvcWithdrawFees,
	)
	return sdk.Result{
		Tags: resTags,
	}
}

func handleMsgSvcWithdrawTax(ctx sdk.Context, k Keeper, msg MsgSvcWithdrawTax) sdk.Result {
	_, found := k.gk.GetTrustee(ctx, msg.Trustee)
	if !found {
		return ErrNotTrustee(k.Codespace(), msg.Trustee).Result()
	}
	oldTaxPool := k.GetServiceFeeTaxPool(ctx)
	newTaxPool, hasNeg := oldTaxPool.SafeMinus(msg.Amount)
	if hasNeg {
		errMsg := fmt.Sprintf("%s is less than %s", oldTaxPool, msg.Amount)
		return sdk.ErrInsufficientFunds(errMsg).Result()
	}
	if !newTaxPool.IsNotNegative() {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("%s is less than %s", oldTaxPool, msg.Amount)).Result()
	}
	_, _, err := k.ck.AddCoins(ctx, msg.DestAddress, msg.Amount)
	if err != nil {
		return err.Result()
	}

	k.SetServiceFeeTaxPool(ctx, newTaxPool)
	resTags := sdk.NewTags(
		tags.Action, tags.ActionSvcWithdrawTax,
	)
	return sdk.Result{
		Tags: resTags,
	}
}

// Called every block, update request status
func EndBlocker(ctx sdk.Context, keeper Keeper) (resTags sdk.Tags) {

	// Reset the intra-transaction counter.
	keeper.SetIntraTxCounter(ctx, 0)

	logger := ctx.Logger().With("module", "service")
	resTags = sdk.NewTags()

	activeIterator := keeper.ActiveRequestQueueIterator(ctx, ctx.BlockHeight())
	defer activeIterator.Close()
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

	return resTags
}
