package stake

import (
	"github.com/irisnet/irishub/app/v1/stake/keeper"
	"github.com/irisnet/irishub/app/v1/stake/tags"
	"github.com/irisnet/irishub/types"
	types2 "github.com/tendermint/tendermint/abci/types"
)

func EndBlocker(ctx types.Context, k keeper.Keeper) (validatorUpdates []types2.ValidatorUpdate) {
	ctx = ctx.WithCoinFlowTrigger(types.StakeEndBlocker)
	ctx = ctx.WithLogger(ctx.Logger().With("handler", "endBlock").With("module", "iris/stake"))
	endBlockerTags := types.EmptyTags()
	// Calculate validator set changes.
	//
	// NOTE: ApplyAndReturnValidatorSetUpdates has to come before
	// UnbondAllMatureValidatorQueue.
	// This fixes a bug when the unbonding period is instant (is the case in
	// some of the tests). The test expected the validator to be completely
	// unbonded after the Endblocker (go from Bonded -> Unbonding during
	// ApplyAndReturnValidatorSetUpdates and then Unbonding -> Unbonded during
	// UnbondAllMatureValidatorQueue).
	validatorUpdates = k.ApplyAndReturnValidatorSetUpdates(ctx)

	// Unbond all mature validators from the unbonding queue.
	k.UnbondAllMatureValidatorQueue(ctx)

	// Remove all mature unbonding delegations from the ubd queue.
	matureUnbonds := k.DequeueAllMatureUnbondingQueue(ctx, ctx.BlockHeader().Time)
	for _, dvPair := range matureUnbonds {
		err := k.CompleteUnbonding(ctx, dvPair.DelegatorAddr, dvPair.ValidatorAddr)
		if err != nil {
			continue
		}
		endBlockerTags.AppendTags(types.NewTags(
			tags.Action, ActionCompleteUnbonding,
			tags.Delegator, []byte(dvPair.DelegatorAddr.String()),
			tags.SrcValidator, []byte(dvPair.ValidatorAddr.String()),
		))
	}

	// Remove all mature redelegations from the red queue.
	matureRedelegations := k.DequeueAllMatureRedelegationQueue(ctx, ctx.BlockHeader().Time)
	for _, dvvTriplet := range matureRedelegations {
		err := k.CompleteRedelegation(ctx, dvvTriplet.DelegatorAddr, dvvTriplet.ValidatorSrcAddr, dvvTriplet.ValidatorDstAddr)
		if err != nil {
			continue
		}
		endBlockerTags.AppendTags(types.NewTags(
			tags.Action, tags.ActionCompleteRedelegation,
			tags.Delegator, []byte(dvvTriplet.DelegatorAddr.String()),
			tags.SrcValidator, []byte(dvvTriplet.ValidatorSrcAddr.String()),
			tags.DstValidator, []byte(dvvTriplet.ValidatorDstAddr.String()),
		))
	}
	k.UpdateMetrics(ctx)
	return
}
