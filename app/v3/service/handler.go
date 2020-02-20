package service

import (
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

// NewHandler returns a handler for all the "service" type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgDefineService:
			return handleMsgDefineService(ctx, k, msg)
		case MsgBindService:
			return handleMsgBindService(ctx, k, msg)
		case MsgUpdateServiceBinding:
			return handleMsgUpdateServiceBinding(ctx, k, msg)
		case MsgSetWithdrawAddress:
			return handleMsgSetWithdrawAddress(ctx, k, msg)
		case MsgDisableService:
			return handleMsgDisableService(ctx, k, msg)
		case MsgEnableService:
			return handleMsgEnableService(ctx, k, msg)
		case MsgRefundServiceDeposit:
			return handleMsgRefundServiceDeposit(ctx, k, msg)
		// case MsgSvcRequest:
		// 	return handleMsgSvcRequest(ctx, k, msg)
		// case MsgSvcResponse:
		// 	return handleMsgSvcResponse(ctx, k, msg)
		// case MsgSvcRefundFees:
		// 	return handleMsgSvcRefundFees(ctx, k, msg)
		// case MsgSvcWithdrawFees:
		// 	return handleMsgSvcWithdrawFees(ctx, k, msg)
		// case MsgSvcWithdrawTax:
		// 	return handleMsgSvcWithdrawTax(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized service message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleMsgDefineService handles MsgDefineService
func handleMsgDefineService(ctx sdk.Context, k Keeper, msg MsgDefineService) sdk.Result {
	if err := k.AddServiceDefinition(
		ctx, msg.Name, msg.Description, msg.Tags,
		msg.Author, msg.AuthorDescription, msg.Schemas,
	); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

// handleMsgBindService handles MsgBindService
func handleMsgBindService(ctx sdk.Context, k Keeper, msg MsgBindService) sdk.Result {
	if err := k.AddServiceBinding(
		ctx, msg.ServiceName, msg.Provider, msg.Deposit,
		msg.Pricing, msg.WithdrawAddress,
	); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

// handleMsgUpdateServiceBinding handles MsgUpdateServiceBinding
func handleMsgUpdateServiceBinding(ctx sdk.Context, k Keeper, msg MsgUpdateServiceBinding) sdk.Result {
	if err := k.UpdateServiceBinding(
		ctx, msg.ServiceName, msg.Provider,
		msg.Deposit, msg.Pricing,
	); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

// handleMsgSetWithdrawAddress handles MsgSetWithdrawAddress
func handleMsgSetWithdrawAddress(ctx sdk.Context, k Keeper, msg MsgSetWithdrawAddress) sdk.Result {
	if err := k.SetWithdrawAddress(
		ctx, msg.ServiceName, msg.Provider, msg.WithdrawAddress,
	); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

// handleMsgDisableService handles MsgDisableService
func handleMsgDisableService(ctx sdk.Context, k Keeper, msg MsgDisableService) sdk.Result {
	if err := k.DisableService(ctx, msg.ServiceName, msg.Provider); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

// handleMsgEnableService handles MsgEnableService
func handleMsgEnableService(ctx sdk.Context, k Keeper, msg MsgEnableService) sdk.Result {
	if err := k.EnableService(ctx, msg.ServiceName, msg.Provider, msg.Deposit); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

// handleMsgRefundServiceDeposit handles MsgRefundServiceDeposit
func handleMsgRefundServiceDeposit(ctx sdk.Context, k Keeper, msg MsgRefundServiceDeposit) sdk.Result {
	if err := k.RefundDeposit(ctx, msg.ServiceName, msg.Provider); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

// TODO

// // handleMsgSvcRequest handles MsgSvcRequest
// func handleMsgSvcRequest(ctx sdk.Context, k Keeper, msg MsgSvcRequest) sdk.Result {
// 	request, err := k.AddRequest(ctx, msg.DefChainID, msg.DefName, msg.BindChainID, msg.ReqChainID,
// 		msg.Consumer, msg.Provider, msg.MethodID, msg.Input, msg.ServiceFee, msg.Profiling)
// 	if err != nil {
// 		return err.Result()
// 	}

// 	ctx.Logger().Debug("Service request", "def_name", request.DefName, "def_chain_id", request.DefChainID,
// 		"provider", request.Provider.String(), "consumer", request.Consumer.String(), "method_id", request.MethodID,
// 		"service_fee", request.ServiceFee, "request_id", request.RequestID())

// 	resTags := sdk.NewTags(
// 		types.TagRequestID, []byte(request.RequestID()),
// 		types.TagProvider, []byte(request.Provider.String()),
// 		types.TagConsumer, []byte(request.Consumer.String()),
// 		types.TagServiceFee, []byte(request.ServiceFee.String()),
// 	)

// 	return sdk.Result{
// 		Tags: resTags,
// 	}
// }

// // handleMsgSvcResponse handles MsgSvcResponse
// func handleMsgSvcResponse(ctx sdk.Context, k Keeper, msg MsgSvcResponse) sdk.Result {
// 	response, err := k.AddResponse(ctx, msg.ReqChainID, msg.RequestID, msg.Provider, msg.Output, msg.ErrorMsg)
// 	if err != nil {
// 		return err.Result()
// 	}

// 	ctx.Logger().Debug("Service response", "request_id", msg.RequestID,
// 		"provider", response.Provider.String(), "consumer", response.Consumer.String())

// 	resTags := sdk.NewTags(
// 		types.TagRequestID, []byte(msg.RequestID),
// 		types.TagConsumer, []byte(response.Consumer.String()),
// 		types.TagProvider, []byte(response.Provider.String()),
// 	)

// 	return sdk.Result{
// 		Tags: resTags,
// 	}
// }

// // handleMsgSvcRefundFees handles MsgSvcRefundFees
// func handleMsgSvcRefundFees(ctx sdk.Context, k Keeper, msg MsgSvcRefundFees) sdk.Result {
// 	if err := k.RefundFee(ctx, msg.Consumer); err != nil {
// 		return err.Result()
// 	}

// 	return sdk.Result{}
// }

// // handleMsgSvcWithdrawFees handles MsgSvcWithdrawFees
// func handleMsgSvcWithdrawFees(ctx sdk.Context, k Keeper, msg MsgSvcWithdrawFees) sdk.Result {
// 	if err := k.WithdrawFee(ctx, msg.Provider); err != nil {
// 		return err.Result()
// 	}

// 	return sdk.Result{}
// }

// // handleMsgSvcWithdrawTax handles MsgSvcWithdrawTax
// func handleMsgSvcWithdrawTax(ctx sdk.Context, k Keeper, msg MsgSvcWithdrawTax) sdk.Result {
// 	if err := k.WithdrawTax(ctx, msg.Trustee, msg.DestAddress, msg.Amount); err != nil {
// 		return err.Result()
// 	}

// 	return sdk.Result{}
// }
