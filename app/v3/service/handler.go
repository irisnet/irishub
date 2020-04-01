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
		case MsgDisableServiceBinding:
			return handleMsgDisableServiceBinding(ctx, k, msg)
		case MsgEnableServiceBinding:
			return handleMsgEnableServiceBinding(ctx, k, msg)
		case MsgRefundServiceDeposit:
			return handleMsgRefundServiceDeposit(ctx, k, msg)
		case MsgCallService:
			return handleMsgCallService(ctx, k, msg)
		case MsgRespondService:
			return handleMsgRespondService(ctx, k, msg)
		case MsgPauseRequestContext:
			return handleMsgPauseRequestContext(ctx, k, msg)
		case MsgStartRequestContext:
			return handleMsgStartRequestContext(ctx, k, msg)
		case MsgKillRequestContext:
			return handleMsgKillRequestContext(ctx, k, msg)
		case MsgUpdateRequestContext:
			return handleMsgUpdateRequestContext(ctx, k, msg)
		case MsgWithdrawEarnedFees:
			return handleMsgWithdrawEarnedFees(ctx, k, msg)
		case MsgWithdrawTax:
			return handleMsgWithdrawTax(ctx, k, msg)
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

	tags := sdk.NewTags(
		TagAuthor, []byte(msg.Author.String()),
	)

	return sdk.Result{
		Tags: tags,
	}
}

// handleMsgBindService handles MsgBindService
func handleMsgBindService(ctx sdk.Context, k Keeper, msg MsgBindService) sdk.Result {
	if err := k.AddServiceBinding(
		ctx, msg.ServiceName, msg.Provider,
		msg.Deposit, msg.Pricing, msg.MinRespTime,
	); err != nil {
		return err.Result()
	}

	tags := sdk.NewTags(
		TagProvider, []byte(msg.Provider.String()),
	)

	return sdk.Result{
		Tags: tags,
	}
}

// handleMsgUpdateServiceBinding handles MsgUpdateServiceBinding
func handleMsgUpdateServiceBinding(ctx sdk.Context, k Keeper, msg MsgUpdateServiceBinding) sdk.Result {
	if err := k.UpdateServiceBinding(
		ctx, msg.ServiceName, msg.Provider,
		msg.Deposit, msg.Pricing, msg.MinRespTime,
	); err != nil {
		return err.Result()
	}

	tags := sdk.NewTags(
		TagProvider, []byte(msg.Provider.String()),
	)

	return sdk.Result{
		Tags: tags,
	}
}

// handleMsgSetWithdrawAddress handles MsgSetWithdrawAddress
func handleMsgSetWithdrawAddress(ctx sdk.Context, k Keeper, msg MsgSetWithdrawAddress) sdk.Result {
	k.SetWithdrawAddress(ctx, msg.Provider, msg.WithdrawAddress)

	tags := sdk.NewTags(
		TagProvider, []byte(msg.Provider.String()),
	)

	return sdk.Result{
		Tags: tags,
	}
}

// handleMsgDisableServiceBinding handles MsgDisableServiceBinding
func handleMsgDisableServiceBinding(ctx sdk.Context, k Keeper, msg MsgDisableServiceBinding) sdk.Result {
	if err := k.DisableServiceBinding(ctx, msg.ServiceName, msg.Provider); err != nil {
		return err.Result()
	}

	tags := sdk.NewTags(
		TagProvider, []byte(msg.Provider.String()),
	)

	return sdk.Result{
		Tags: tags,
	}
}

// handleMsgEnableServiceBinding handles MsgEnableServiceBinding
func handleMsgEnableServiceBinding(ctx sdk.Context, k Keeper, msg MsgEnableServiceBinding) sdk.Result {
	if err := k.EnableServiceBinding(ctx, msg.ServiceName, msg.Provider, msg.Deposit); err != nil {
		return err.Result()
	}

	tags := sdk.NewTags(
		TagProvider, []byte(msg.Provider.String()),
	)

	return sdk.Result{
		Tags: tags,
	}
}

// handleMsgRefundServiceDeposit handles MsgRefundServiceDeposit
func handleMsgRefundServiceDeposit(ctx sdk.Context, k Keeper, msg MsgRefundServiceDeposit) sdk.Result {
	if err := k.RefundDeposit(ctx, msg.ServiceName, msg.Provider); err != nil {
		return err.Result()
	}

	tags := sdk.NewTags(
		TagProvider, []byte(msg.Provider.String()),
	)

	return sdk.Result{
		Tags: tags,
	}
}

// handleMsgCallService handles MsgCallService
func handleMsgCallService(ctx sdk.Context, k Keeper, msg MsgCallService) sdk.Result {
	_, tags, err := k.CreateRequestContext(
		ctx, msg.ServiceName, msg.Providers, msg.Consumer, msg.Input, msg.ServiceFeeCap, msg.Timeout,
		msg.SuperMode, msg.Repeated, msg.RepeatedFrequency, msg.RepeatedTotal, RUNNING, 0, "")
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}

// handleMsgRespondService handles MsgRespondService
func handleMsgRespondService(ctx sdk.Context, k Keeper, msg MsgRespondService) sdk.Result {
	request, response, tags, err := k.AddResponse(ctx, msg.RequestID, msg.Provider, msg.Result, msg.Output)
	if err != nil {
		return err.Result()
	}

	tags = tags.AppendTags(
		sdk.NewTags(
			TagRequestID, []byte(msg.RequestID.String()),
			TagRequestContextID, []byte(request.RequestContextID.String()),
			TagConsumer, []byte(response.Consumer.String()),
			TagProvider, []byte(response.Provider.String()),
			TagServiceName, []byte(request.ServiceName),
		),
	)

	return sdk.Result{
		Tags: tags,
	}
}

// handleMsgPauseRequestContext handles MsgPauseRequestContext
func handleMsgPauseRequestContext(ctx sdk.Context, k Keeper, msg MsgPauseRequestContext) sdk.Result {
	if err := k.CheckAuthority(ctx, msg.Consumer, msg.RequestContextID, true); err != nil {
		return err.Result()
	}

	if err := k.PauseRequestContext(ctx, msg.RequestContextID, msg.Consumer); err != nil {
		return err.Result()
	}

	tags := sdk.NewTags(
		TagConsumer, []byte(msg.Consumer.String()),
	)

	return sdk.Result{
		Tags: tags,
	}
}

// handleMsgStartRequestContext handles MsgStartRequestContext
func handleMsgStartRequestContext(ctx sdk.Context, k Keeper, msg MsgStartRequestContext) sdk.Result {
	if err := k.CheckAuthority(ctx, msg.Consumer, msg.RequestContextID, true); err != nil {
		return err.Result()
	}

	if err := k.StartRequestContext(ctx, msg.RequestContextID, msg.Consumer); err != nil {
		return err.Result()
	}

	tags := sdk.NewTags(
		TagConsumer, []byte(msg.Consumer.String()),
	)

	return sdk.Result{
		Tags: tags,
	}
}

// handleMsgKillRequestContext handles MsgKillRequestContext
func handleMsgKillRequestContext(ctx sdk.Context, k Keeper, msg MsgKillRequestContext) sdk.Result {
	if err := k.CheckAuthority(ctx, msg.Consumer, msg.RequestContextID, true); err != nil {
		return err.Result()
	}

	if err := k.KillRequestContext(ctx, msg.RequestContextID, msg.Consumer); err != nil {
		return err.Result()
	}

	tags := sdk.NewTags(
		TagConsumer, []byte(msg.Consumer.String()),
	)

	return sdk.Result{
		Tags: tags,
	}
}

// handleMsgUpdateRequestContext handles MsgUpdateRequestContext
func handleMsgUpdateRequestContext(ctx sdk.Context, k Keeper, msg MsgUpdateRequestContext) sdk.Result {
	if err := k.CheckAuthority(ctx, msg.Consumer, msg.RequestContextID, true); err != nil {
		return err.Result()
	}

	if err := k.UpdateRequestContext(
		ctx, msg.RequestContextID, msg.Providers, msg.ServiceFeeCap,
		msg.Timeout, msg.RepeatedFrequency, msg.RepeatedTotal, msg.Consumer,
	); err != nil {
		return err.Result()
	}

	tags := sdk.NewTags(
		TagConsumer, []byte(msg.Consumer.String()),
	)

	return sdk.Result{
		Tags: tags,
	}
}

// handleMsgWithdrawEarnedFees handles MsgWithdrawEarnedFees
func handleMsgWithdrawEarnedFees(ctx sdk.Context, k Keeper, msg MsgWithdrawEarnedFees) sdk.Result {
	if err := k.WithdrawEarnedFees(ctx, msg.Provider); err != nil {
		return err.Result()
	}

	tags := sdk.NewTags(
		TagProvider, []byte(msg.Provider.String()),
	)

	return sdk.Result{
		Tags: tags,
	}
}

// handleMsgWithdrawTax handles MsgWithdrawTax
func handleMsgWithdrawTax(ctx sdk.Context, k Keeper, msg MsgWithdrawTax) sdk.Result {
	if err := k.WithdrawTax(ctx, msg.Trustee, msg.DestAddress, msg.Amount); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}
