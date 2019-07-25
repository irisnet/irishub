package keeper

import (
	"fmt"
	"strconv"

	"github.com/irisnet/irishub/modules/stake/types"
	sdk "github.com/irisnet/irishub/types"
)

const (
	// slash-unbondind-[delegatorAddr]-[validatorAddr]
	SlashUnbondindDelegation = "slash-unbondind-%s-%s"
	// slash-redelegation-[delegatorAddr]-[validatorAddrSrc]-[validatorAddrDst]
	SlashRedelegation = "slash-redelegation-%s-%s-%s"
	// slash-validator-[validatorAddr]
	SlashValidator = "slash-validator-%s"
	// slash-validator-redelegation-[validatorAddrDst]-[validatorAddrSrc]-[delegatorAddr]
	SlashValidatorRedelegation = "slash-validator-redelegation-%s-%s-%s"
)

// Slash a validator for an infraction committed at a known height
// Find the contributing stake at that height and burn the specified slashFactor
// of it, updating unbonding delegation & redelegations appropriately
//
// CONTRACT:
//    slashFactor is non-negative
// CONTRACT:
//    Infraction committed equal to or less than an unbonding period in the past,
//    so all unbonding delegations and redelegations from that height are stored
// CONTRACT:
//    Slash will not slash unbonded validators (for the above reason)
// CONTRACT:
//    Infraction committed at the current height or at a past height,
//    not at a height in the future
func (k Keeper) Slash(ctx sdk.Context, consAddr sdk.ConsAddress, infractionHeight int64, power int64, slashFactor sdk.Dec) (tags sdk.Tags) {
	logger := ctx.Logger()

	if slashFactor.LT(sdk.ZeroDec()) {
		panic(fmt.Errorf("attempted to slash with a negative slash factor: %v", slashFactor))
	}

	// Amount of slashing = slash slashFactor * power at time of infraction
	slashAmount := sdk.NewDec(power).Mul(slashFactor)
	// ref https://github.com/irisnet/irishub/issues/1348
	// ref https://github.com/irisnet/irishub/issues/1471
	//Multiply 1*10^18 to calculate equivalent stake denom amount
	slashAmount = slashAmount.MulInt(sdk.AttoScaleFactor).TruncateDec()

	validator, found := k.GetValidatorByConsAddr(ctx, consAddr)
	if !found {
		// If not found, the validator must have been overslashed and removed - so we don't need to do anything
		// NOTE:  Correctness dependent on invariant that unbonding delegations / redelegations must also have been completely
		//        slashed in this case - which we don't explicitly check, but should be true.
		// Log the slash attempt for future reference (maybe we should tag it too)
		logger.Error(fmt.Sprintf(
			"WARNING: Ignored attempt to slash a nonexistent validator with address %s, we recommend you investigate immediately",
			consAddr))
		return
	}

	// should not be slashing unbonded
	if validator.Status == sdk.Unbonded {
		panic(fmt.Sprintf("should not be slashing unbonded validator: %s", validator.GetOperator()))
	}

	operatorAddress := validator.GetOperator()
	k.OnValidatorModified(ctx, operatorAddress)

	// Track remaining slash amount for the validator
	// This will decrease when we slash unbondings and
	// redelegations, as that stake has since unbonded
	remainingSlashAmount := slashAmount

	switch {
	case infractionHeight > ctx.BlockHeight():

		// Can't slash infractions in the future
		panic(fmt.Sprintf(
			"impossible attempt to slash future infraction at height %d but we are at height %d",
			infractionHeight, ctx.BlockHeight()))

	case infractionHeight == ctx.BlockHeight():

		// Special-case slash at current height for efficiency - we don't need to look through unbonding delegations or redelegations
		logger.Info("Slashing at current height, not scanning unbonding delegations & redelegations", "infraction_height", infractionHeight)

	case infractionHeight < ctx.BlockHeight():

		// Iterate through unbonding delegations from slashed validator
		unbondingDelegations := k.GetUnbondingDelegationsFromValidator(ctx, operatorAddress)
		for _, unbondingDelegation := range unbondingDelegations {
			amountSlashed, slashUnbondingTags := k.slashUnbondingDelegation(ctx, unbondingDelegation, infractionHeight, slashFactor)
			tags = tags.AppendTags(slashUnbondingTags)
			if amountSlashed.IsZero() {
				continue
			}
			remainingSlashAmount = remainingSlashAmount.Sub(amountSlashed)
		}

		// Iterate through redelegations from slashed validator
		redelegations := k.GetRedelegationsFromValidator(ctx, operatorAddress)
		for _, redelegation := range redelegations {
			amountSlashed, slashRedelegationTags := k.slashRedelegation(ctx, validator, redelegation, infractionHeight, slashFactor)
			tags = tags.AppendTags(slashRedelegationTags)
			if amountSlashed.IsZero() {
				continue
			}
			remainingSlashAmount = remainingSlashAmount.Sub(amountSlashed)
		}
	}

	// cannot decrease balance below zero
	tokensToBurn := sdk.MinDec(remainingSlashAmount, validator.Tokens)
	tokensToBurn = sdk.MaxDec(tokensToBurn, sdk.ZeroDec()) // defensive.
	tags = tags.AppendTag(fmt.Sprintf(SlashValidator, validator.OperatorAddr), []byte(tokensToBurn.String()))
	// Deduct from validator's bonded tokens and update the validator.
	// The deducted tokens are returned to pool.LooseTokens.
	validator = k.RemoveValidatorTokens(ctx, validator, tokensToBurn)
	if !tokensToBurn.Sub(tokensToBurn.TruncateDec()).IsZero() {
		panic("slash decimal token in redelegation")
	}
	k.bankKeeper.DecreaseLoosenToken(ctx, sdk.Coins{sdk.NewCoin(types.StakeDenom, tokensToBurn.TruncateInt())})
	slashToken, err := strconv.ParseFloat(tokensToBurn.QuoInt(sdk.AttoScaleFactor).String(), 64)
	if err == nil {
		k.metrics.SlashedToken.With("validator_address", validator.GetConsAddr().String()).Add(slashToken)
	}
	// Log that a slash occurred!
	logger.Info("Validator slashed", "consensus_address", validator.GetConsAddr().String(),
		"operator_address", validator.GetOperator().String(), "slash_factor", slashFactor.String(), "slash_tokens", tokensToBurn)
	// TODO Return event(s), blocked on https://github.com/tendermint/tendermint/pull/1803
	return
}

// jail a validator
func (k Keeper) Jail(ctx sdk.Context, consAddr sdk.ConsAddress) {
	validator := k.mustGetValidatorByConsAddr(ctx, consAddr)
	k.jailValidator(ctx, validator)
	// TODO Return event(s), blocked on https://github.com/tendermint/tendermint/pull/1803
	return
}

// unjail a validator
func (k Keeper) Unjail(ctx sdk.Context, consAddr sdk.ConsAddress) {
	validator := k.mustGetValidatorByConsAddr(ctx, consAddr)
	k.unjailValidator(ctx, validator)
	k.metrics.Jailed.With("validator_address", validator.GetConsAddr().String()).Set(0)
	// TODO Return event(s), blocked on https://github.com/tendermint/tendermint/pull/1803
	return
}

// slash an unbonding delegation and update the pool
// return the amount that would have been slashed assuming
// the unbonding delegation had enough stake to slash
// (the amount actually slashed may be less if there's
// insufficient stake remaining)
func (k Keeper) slashUnbondingDelegation(ctx sdk.Context, unbondingDelegation types.UnbondingDelegation,
	infractionHeight int64, slashFactor sdk.Dec) (slashAmount sdk.Dec, tags sdk.Tags) {

	now := ctx.BlockHeader().Time

	// If unbonding started before this height, stake didn't contribute to infraction
	if unbondingDelegation.CreationHeight < infractionHeight {
		return sdk.ZeroDec(), nil
	}

	if unbondingDelegation.MinTime.Before(now) {
		// Unbonding delegation no longer eligible for slashing, skip it
		// TODO Settle and delete it automatically?
		return sdk.ZeroDec(), nil
	}

	// Calculate slash amount proportional to stake contributing to infraction
	slashAmount = slashFactor.MulInt(unbondingDelegation.InitialBalance.Amount)

	// Don't slash more tokens than held
	// Possible since the unbonding delegation may already
	// have been slashed, and slash amounts are calculated
	// according to stake held at time of infraction
	unbondingSlashAmount := sdk.MinInt(slashAmount.RoundInt(), unbondingDelegation.Balance.Amount)

	// Update unbonding delegation if necessary
	if !unbondingSlashAmount.IsZero() {
		unbondingDelegation.Balance.Amount = unbondingDelegation.Balance.Amount.Sub(unbondingSlashAmount)
		tags = tags.AppendTag(fmt.Sprintf(SlashUnbondindDelegation, unbondingDelegation.DelegatorAddr, unbondingDelegation.ValidatorAddr), []byte(unbondingSlashAmount.String()))
		k.SetUnbondingDelegation(ctx, unbondingDelegation)
		k.bankKeeper.DecreaseLoosenToken(ctx, sdk.Coins{sdk.NewCoin(types.StakeDenom, unbondingSlashAmount)})
	}

	return
}

// slash a redelegation and update the pool
// return the amount that would have been slashed assuming
// the unbonding delegation had enough stake to slash
// (the amount actually slashed may be less if there's
// insufficient stake remaining)
// nolint: unparam
func (k Keeper) slashRedelegation(ctx sdk.Context, validator types.Validator, redelegation types.Redelegation,
	infractionHeight int64, slashFactor sdk.Dec) (slashAmount sdk.Dec, tags sdk.Tags) {

	now := ctx.BlockHeader().Time

	// If redelegation started before this height, stake didn't contribute to infraction
	if redelegation.CreationHeight < infractionHeight {
		return sdk.ZeroDec(), nil
	}

	if redelegation.MinTime.Before(now) {
		// Redelegation no longer eligible for slashing, skip it
		// TODO Delete it automatically?
		return sdk.ZeroDec(), nil
	}

	// Calculate slash amount proportional to stake contributing to infraction
	slashAmount = slashFactor.MulInt(redelegation.InitialBalance.Amount)

	// Don't slash more tokens than held
	// Possible since the redelegation may already
	// have been slashed, and slash amounts are calculated
	// according to stake held at time of infraction
	redelegationSlashAmount := sdk.MinInt(slashAmount.RoundInt(), redelegation.Balance.Amount)

	// Update redelegation if necessary
	if !redelegationSlashAmount.IsZero() {
		redelegation.Balance.Amount = redelegation.Balance.Amount.Sub(redelegationSlashAmount)
		k.SetRedelegation(ctx, redelegation)
		tags = tags.AppendTag(fmt.Sprintf(SlashRedelegation, redelegation.DelegatorAddr, redelegation.ValidatorSrcAddr, redelegation.ValidatorDstAddr), []byte(redelegationSlashAmount.String()))
	}

	// Unbond from target validator
	sharesToUnbond := slashFactor.Mul(redelegation.SharesDst)
	if !sharesToUnbond.IsZero() {
		delegation, found := k.GetDelegation(ctx, redelegation.DelegatorAddr, redelegation.ValidatorDstAddr)
		if !found {
			// If deleted, delegation has zero shares, and we can't unbond any more
			return slashAmount, nil
		}
		if sharesToUnbond.GT(delegation.Shares) {
			sharesToUnbond = delegation.Shares
		}

		tokensToBurn, err := k.unbond(ctx, redelegation.DelegatorAddr, redelegation.ValidatorDstAddr, sharesToUnbond)
		if err != nil {
			panic(fmt.Errorf("error unbonding delegator: %v", err))
		}
		tags = tags.AppendTag(fmt.Sprintf(SlashValidatorRedelegation, redelegation.ValidatorDstAddr, redelegation.ValidatorSrcAddr, redelegation.DelegatorAddr), []byte(tokensToBurn.String()))
		k.bankKeeper.DecreaseLoosenToken(ctx, sdk.Coins{sdk.NewCoin(types.StakeDenom, tokensToBurn.TruncateInt())})
	}

	return
}
