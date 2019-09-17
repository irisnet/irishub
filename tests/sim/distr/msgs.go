package distr

import (
	"fmt"
	"math/rand"

	"github.com/irisnet/irishub/mock/baseapp"
	"github.com/irisnet/irishub/mock/simulation"
	"github.com/irisnet/irishub/modules/auth"
	"github.com/irisnet/irishub/modules/distribution"
	sdk "github.com/irisnet/irishub/types"
)

// SimulateMsgWithdrawDelegatorRewardsAll
func SimulateMsgWithdrawDelegatorRewardsAll(m auth.AccountKeeper, k distribution.Keeper) simulation.Operation {
	handler := distribution.NewHandler(k)
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account, event func(string)) (
		action string, fOp []simulation.FutureOperation, err error) {

		account := simulation.RandomAcc(r, accs)
		msg := distribution.NewMsgWithdrawDelegatorRewardsAll(account.Address)

		if msg.ValidateBasic() != nil {
			return "", nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}

		ctx, write := ctx.CacheContext()
		result := handler(ctx, msg)
		if result.IsOK() {
			write()
		}

		event(fmt.Sprintf("distribution/MsgWithdrawDelegatorRewardsAll/%v", result.IsOK()))

		action = fmt.Sprintf("TestMsgWithdrawDelegatorRewardsAll: ok %v, msg %s", result.IsOK(), msg.GetSignBytes())
		return action, nil, nil
	}
}

// SimulateMsgWithdrawDelegatorReward
func SimulateMsgWithdrawDelegatorReward(m auth.AccountKeeper, k distribution.Keeper) simulation.Operation {
	handler := distribution.NewHandler(k)
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account, event func(string)) (
		action string, fOp []simulation.FutureOperation, err error) {

		delegatorAccount := simulation.RandomAcc(r, accs)
		validatorAccount := simulation.RandomAcc(r, accs)
		msg := distribution.NewMsgWithdrawDelegatorReward(delegatorAccount.Address, sdk.ValAddress(validatorAccount.Address))

		if msg.ValidateBasic() != nil {
			return "", nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}

		ctx, write := ctx.CacheContext()
		result := handler(ctx, msg)
		if result.IsOK() {
			write()
		}

		event(fmt.Sprintf("distribution/MsgWithdrawDelegatorReward/%v", result.IsOK()))

		action = fmt.Sprintf("TestMsgWithdrawDelegatorReward: ok %v, msg %s", result.IsOK(), msg.GetSignBytes())
		return action, nil, nil
	}
}

// SimulateMsgWithdrawValidatorRewardsAll
func SimulateMsgWithdrawValidatorRewardsAll(m auth.AccountKeeper, k distribution.Keeper) simulation.Operation {
	handler := distribution.NewHandler(k)
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account, event func(string)) (
		action string, fOp []simulation.FutureOperation, err error) {

		account := simulation.RandomAcc(r, accs)
		msg := distribution.NewMsgWithdrawValidatorRewardsAll(sdk.ValAddress(account.Address))

		if msg.ValidateBasic() != nil {
			return "", nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}

		ctx, write := ctx.CacheContext()
		result := handler(ctx, msg)
		if result.IsOK() {
			write()
		}

		event(fmt.Sprintf("distribution/MsgWithdrawValidatorRewardsAll/%v", result.IsOK()))

		action = fmt.Sprintf("TestMsgWithdrawValidatorRewardsAll: ok %v, msg %s", result.IsOK(), msg.GetSignBytes())
		return action, nil, nil
	}
}
