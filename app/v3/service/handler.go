package service

import (
	"fmt"

	"github.com/irisnet/irishub/app/v3/service/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

// NewHandler returns a handler for all the "service" type messages
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
			errMsg := fmt.Sprintf("unrecognized service message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleMsgSvcDef handles MsgSvcDef
func handleMsgSvcDef(ctx sdk.Context, k Keeper, msg MsgSvcDef) sdk.Result {
	if err := k.AddServiceDefinition(
		ctx, msg.Name, msg.ChainID, msg.Description, msg.Tags,
		msg.Author, msg.AuthorDescription, msg.IDLContent,
	); err != nil {
		return err.Result()
	}

	ctx.Logger().Info("Create service definition", "name", msg.Name, "author", msg.Author.String())

	return sdk.Result{}
}

// handleMsgSvcBind handles MsgSvcBind
func handleMsgSvcBind(ctx sdk.Context, k Keeper, msg MsgSvcBind) sdk.Result {
	if err := k.AddServiceBinding(
		ctx, msg.DefChainID, msg.DefName, msg.BindChainID,
		msg.Provider, msg.BindingType, msg.Deposit, msg.Prices, msg.Level,
	); err != nil {
		return err.Result()
	}

	ctx.Logger().Info("Add service binding", "def_name", msg.DefName, "def_chain_id", msg.DefChainID,
		"provider", msg.Provider.String(), "binding_type", msg.BindingType.String())

	return sdk.Result{}
}

// handleMsgSvcBindUpdate handles MsgSvcBindingUpdate
func handleMsgSvcBindUpdate(ctx sdk.Context, k Keeper, msg MsgSvcBindingUpdate) sdk.Result {
	svcBinding, err := k.UpdateServiceBinding(ctx, msg.DefChainID, msg.DefName, msg.BindChainID,
		msg.Provider, msg.BindingType, msg.Deposit, msg.Prices, msg.Level)
	if err != nil {
		return err.Result()
	}

	ctx.Logger().Info("Update service binding", "def_name", msg.DefName, "def_chain_id", msg.DefChainID,
		"provider", msg.Provider.String(), "binding_type", svcBinding.BindingType.String())

	return sdk.Result{}
}

// handleMsgSvcDisable handles MsgSvcDisable
func handleMsgSvcDisable(ctx sdk.Context, k Keeper, msg MsgSvcDisable) sdk.Result {
	if err := k.Disable(
		ctx, msg.DefChainID, msg.DefName, msg.BindChainID, msg.Provider,
	); err != nil {
		return err.Result()
	}

	ctx.Logger().Info("Disable service binding", "def_name", msg.DefName, "def_chain_id", msg.DefChainID,
		"provider", msg.Provider.String())

	return sdk.Result{}
}

// handleMsgSvcEnable handles MsgSvcEnable
func handleMsgSvcEnable(ctx sdk.Context, k Keeper, msg MsgSvcEnable) sdk.Result {
	if err := k.Enable(
		ctx, msg.DefChainID, msg.DefName, msg.BindChainID, msg.Provider, msg.Deposit,
	); err != nil {
		return err.Result()
	}

	ctx.Logger().Info("Enable service binding", "def_name", msg.DefName, "def_chain_id", msg.DefChainID,
		"provider", msg.Provider.String())

	return sdk.Result{}
}

// handleMsgSvcRefundDeposit handles MsgSvcRefundDeposit
func handleMsgSvcRefundDeposit(ctx sdk.Context, k Keeper, msg MsgSvcRefundDeposit) sdk.Result {
	if err := k.RefundDeposit(
		ctx, msg.DefChainID, msg.DefName, msg.BindChainID, msg.Provider,
	); err != nil {
		return err.Result()
	}

	ctx.Logger().Info("Refund deposit", "def_name", msg.DefName, "def_chain_id", msg.DefChainID,
		"provider", msg.Provider.String())

	return sdk.Result{}
}

// handleMsgSvcRequest handles MsgSvcRequest
func handleMsgSvcRequest(ctx sdk.Context, k Keeper, msg MsgSvcRequest) sdk.Result {
	request, err := k.AddRequest(ctx, msg.DefChainID, msg.DefName, msg.BindChainID, msg.ReqChainID,
		msg.Consumer, msg.Provider, msg.MethodID, msg.Input, msg.ServiceFee, msg.Profiling)
	if err != nil {
		return err.Result()
	}

	ctx.Logger().Debug("Service request", "def_name", request.DefName, "def_chain_id", request.DefChainID,
		"provider", request.Provider.String(), "consumer", request.Consumer.String(), "method_id", request.MethodID,
		"service_fee", request.ServiceFee, "request_id", request.RequestID())

	resTags := sdk.NewTags(
		types.TagRequestID, []byte(request.RequestID()),
		types.TagProvider, []byte(request.Provider.String()),
		types.TagConsumer, []byte(request.Consumer.String()),
		types.TagServiceFee, []byte(request.ServiceFee.String()),
	)

	return sdk.Result{
		Tags: resTags,
	}
}

// handleMsgSvcResponse handles MsgSvcResponse
func handleMsgSvcResponse(ctx sdk.Context, k Keeper, msg MsgSvcResponse) sdk.Result {
	response, err := k.AddResponse(ctx, msg.ReqChainID, msg.RequestID, msg.Provider, msg.Output, msg.ErrorMsg)
	if err != nil {
		return err.Result()
	}

	ctx.Logger().Debug("Service response", "request_id", msg.RequestID,
		"provider", response.Provider.String(), "consumer", response.Consumer.String())

	resTags := sdk.NewTags(
		types.TagRequestID, []byte(msg.RequestID),
		types.TagConsumer, []byte(response.Consumer.String()),
		types.TagProvider, []byte(response.Provider.String()),
	)

	return sdk.Result{
		Tags: resTags,
	}
}

// handleMsgSvcRefundFees handles MsgSvcRefundFees
func handleMsgSvcRefundFees(ctx sdk.Context, k Keeper, msg MsgSvcRefundFees) sdk.Result {
	if err := k.RefundFee(ctx, msg.Consumer); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

// handleMsgSvcWithdrawFees handles MsgSvcWithdrawFees
func handleMsgSvcWithdrawFees(ctx sdk.Context, k Keeper, msg MsgSvcWithdrawFees) sdk.Result {
	if err := k.WithdrawFee(ctx, msg.Provider); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

// handleMsgSvcWithdrawTax handles MsgSvcWithdrawTax
func handleMsgSvcWithdrawTax(ctx sdk.Context, k Keeper, msg MsgSvcWithdrawTax) sdk.Result {
	if err := k.WithdrawTax(ctx, msg.Trustee, msg.DestAddress, msg.Amount); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}
