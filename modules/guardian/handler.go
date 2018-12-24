package guardian

import (
	sdk "github.com/irisnet/irishub/types"
)

// handle all "guardian" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgAddProfiler:
			return handleMsgAddProfiler(ctx, k, msg)
		case MsgAddTrustee:
			return handleMsgAddTrustee(ctx, k, msg)
		case MsgDeleteProfiler:
			return handleMsgDeleteProfiler(ctx, k, msg)
		case MsgDeleteTrustee:
			return handleMsgDeleteTrustee(ctx, k, msg)
		default:
			return sdk.ErrTxDecode("invalid message parse in guardian module").Result()
		}
	}
}

func handleMsgAddProfiler(ctx sdk.Context, k Keeper, msg MsgAddProfiler) sdk.Result {
	if profiler, found := k.GetProfiler(ctx, msg.AddedAddress);
		!found || profiler.AccountType != Genesis {
		return ErrInvalidOperator(DefaultCodespace, msg.AddedAddress).Result()
	}
	if _, found := k.GetProfiler(ctx, msg.AddedAddress); found {
		return ErrProfilerExists(DefaultCodespace, msg.Address).Result()
	}
	profiler := NewGuardian(msg.Description, Ordinary, msg.Address, msg.AddedAddress)
	err := k.AddProfiler(ctx, profiler)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
	}
}

func handleMsgAddTrustee(ctx sdk.Context, k Keeper, msg MsgAddTrustee) sdk.Result {
	if trustee, found := k.GetTrustee(ctx, msg.AddedAddress);
		!found || trustee.AccountType != Genesis {
		return ErrInvalidOperator(DefaultCodespace, msg.AddedAddress).Result()
	}
	if _, found := k.GetTrustee(ctx, msg.AddedAddress); found {
		return ErrTrusteeExists(DefaultCodespace, msg.Address).Result()
	}
	trustee := NewGuardian(msg.Description, Ordinary, msg.Address, msg.AddedAddress)
	err := k.AddTrustee(ctx, trustee)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
	}
}

func handleMsgDeleteProfiler(ctx sdk.Context, k Keeper, msg MsgDeleteProfiler) sdk.Result {
	if profiler, found := k.GetProfiler(ctx, msg.DeletedAddress);
		!found || profiler.AccountType != Genesis {
		return ErrInvalidOperator(DefaultCodespace, msg.DeletedAddress).Result()
	}
	if _, found := k.GetProfiler(ctx, msg.Address); found {
		return ErrProfilerNotExists(DefaultCodespace, msg.Address).Result()
	}

	err := k.DeleteProfiler(ctx, msg.Address)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
	}
}

func handleMsgDeleteTrustee(ctx sdk.Context, k Keeper, msg MsgDeleteTrustee) sdk.Result {
	if trustee, found := k.GetTrustee(ctx, msg.DeletedAddress);
		!found || trustee.AccountType != Genesis {
		return ErrInvalidOperator(DefaultCodespace, msg.DeletedAddress).Result()
	}
	if _, found := k.GetTrustee(ctx, msg.Address); found {
		return ErrTrusteeNotExists(DefaultCodespace, msg.Address).Result()
	}

	err := k.DeleteProfiler(ctx, msg.Address)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
	}
}
