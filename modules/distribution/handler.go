package distribution

import (
	"github.com/irisnet/irishub/modules/distribution/keeper"
	"github.com/irisnet/irishub/modules/distribution/tags"
	"github.com/irisnet/irishub/modules/distribution/types"
	sdk "github.com/irisnet/irishub/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		// NOTE msg already has validate basic run
		switch msg := msg.(type) {
		case types.MsgSetWithdrawAddress:
			return handleMsgModifyWithdrawAddress(ctx, msg, k)
		case types.MsgWithdrawDelegatorRewardsAll:
			return handleMsgWithdrawDelegatorRewardsAll(ctx, msg, k)
		case types.MsgWithdrawDelegatorReward:
			return handleMsgWithdrawDelegatorReward(ctx, msg, k)
		case types.MsgWithdrawValidatorRewardsAll:
			return handleMsgWithdrawValidatorRewardsAll(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("invalid message parse in distribution module").Result()
		}
	}
}

//_____________________________________________________________________

// These functions assume everything has been authenticated,
// now we just perform action and save

func handleMsgModifyWithdrawAddress(ctx sdk.Context, msg types.MsgSetWithdrawAddress, k keeper.Keeper) sdk.Result {

	k.SetDelegatorWithdrawAddr(ctx, msg.DelegatorAddr, msg.WithdrawAddr)

	tags := sdk.NewTags(
		tags.Delegator, []byte(msg.DelegatorAddr.String()),
	)
	return sdk.Result{
		Tags: tags,
	}
}

func handleMsgWithdrawDelegatorRewardsAll(ctx sdk.Context, msg types.MsgWithdrawDelegatorRewardsAll, k keeper.Keeper) sdk.Result {
	ctx.Logger().Debug("Withdraw delegator reward", "delegator", msg.DelegatorAddr.String())
	reward, withdrawTags := k.WithdrawDelegationRewardsAll(ctx, msg.DelegatorAddr)
	rewardTruncate, _ := reward.TruncateDecimal()
	resultTags := sdk.NewTags(
		tags.Delegator, []byte(msg.DelegatorAddr.String()),
		tags.Reward, []byte(rewardTruncate.String()),
		tags.WithdrawAddr, []byte(k.GetDelegatorWithdrawAddr(ctx, msg.DelegatorAddr).String()),
	)
	resultTags = resultTags.AppendTags(withdrawTags)
	ctx.Logger().Debug("Withdraw result", "reward", rewardTruncate.String())
	return sdk.Result{
		Tags: resultTags,
	}
}

func handleMsgWithdrawDelegatorReward(ctx sdk.Context, msg types.MsgWithdrawDelegatorReward, k keeper.Keeper) sdk.Result {
	ctx.Logger().Debug("Withdraw delegation reward", "delegator", msg.DelegatorAddr.String(), "validator", msg.ValidatorAddr.String())
	reward, err := k.WithdrawDelegationReward(ctx, msg.DelegatorAddr, msg.ValidatorAddr)
	if err != nil {
		return err.Result()
	}
	rewardTruncate, _ := reward.TruncateDecimal()
	tags := sdk.NewTags(
		tags.Delegator, []byte(msg.DelegatorAddr.String()),
		tags.Validator, []byte(msg.ValidatorAddr.String()),
		tags.Reward, []byte(rewardTruncate.String()),
		tags.WithdrawAddr, []byte(k.GetDelegatorWithdrawAddr(ctx, msg.DelegatorAddr).String()),
	)
	ctx.Logger().Debug("Withdraw result", "reward", rewardTruncate.String())
	return sdk.Result{
		Tags: tags,
	}
}

func handleMsgWithdrawValidatorRewardsAll(ctx sdk.Context, msg types.MsgWithdrawValidatorRewardsAll, k keeper.Keeper) sdk.Result {
	ctx.Logger().Debug("Withdraw validator reward", "validator", msg.ValidatorAddr.String())
	reward, withdrawTags, err := k.WithdrawValidatorRewardsAll(ctx, msg.ValidatorAddr)
	if err != nil {
		return err.Result()
	}
	rewardTruncate, _ := reward.TruncateDecimal()
	resultTags := sdk.NewTags(
		tags.Validator, []byte(msg.ValidatorAddr.String()),
		tags.Reward, []byte(rewardTruncate.String()),
		tags.WithdrawAddr, []byte(k.GetDelegatorWithdrawAddr(ctx, sdk.AccAddress(msg.ValidatorAddr)).String()),
	)
	resultTags = resultTags.AppendTags(withdrawTags)
	ctx.Logger().Debug("Withdraw result", "reward", rewardTruncate.String())
	return sdk.Result{
		Tags: resultTags,
	}
}
