package slashing

import (
	sdk "github.com/irisnet/irishub/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"strconv"
)

// slashing begin block functionality
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, sk Keeper) (tags sdk.Tags) {
	ctx = ctx.WithCoinFlowTrigger(sdk.SlashBeginBlocker)
	ctx = ctx.WithLogger(ctx.Logger().With("handler", "beginBlock").With("module", "iris/slashing"))

	// Tag the height
	tags = sdk.NewTags("height", []byte(strconv.FormatInt(req.Header.Height, 10)))

	// Iterate over all the validators  which *should* have signed this block
	// store whether or not they have actually signed it and slash/unbond any
	// which have missed too many blocks in a row (downtime slashing)
	for _, voteInfo := range req.LastCommitInfo.GetVotes() {
		absenceSlashTags := sk.handleValidatorSignature(ctx, voteInfo.Validator.Address, voteInfo.Validator.Power, voteInfo.SignedLastBlock)
		tags = tags.AppendTags(absenceSlashTags)
	}

	// Iterate through any newly discovered evidence of infraction
	// Slash any validators (and since-unbonded stake within the unbonding period)
	// who contributed to valid infractions
	for _, evidence := range req.ByzantineValidators {
		switch evidence.Type {
		case tmtypes.ABCIEvidenceTypeDuplicateVote:
			doubleSignSlashTag := sk.handleDoubleSign(ctx, evidence.Validator.Address, evidence.Height, evidence.Validator.Power)
			tags = tags.AppendTags(doubleSignSlashTag)
		default:
			ctx.Logger().Error("ignored unknown evidence type", "type", evidence.Type)
		}
	}

	return
}

// slashing end block functionality
func EndBlocker(ctx sdk.Context, req abci.RequestEndBlock, sk Keeper) (tags sdk.Tags) {
	ctx = ctx.WithCoinFlowTrigger(sdk.SlashEndBlocker)
	ctx = ctx.WithLogger(ctx.Logger().With("handler", "endBlock").With("module", "iris/slashing"))
	// Tag the height
	tags = sdk.NewTags("height", []byte(strconv.FormatInt(req.Height, 10)))

	if ctx.ValidTxCounter().Count() < ctx.BlockHeader().NumTxs {
		proposalCensorshipTag := sk.handleProposerCensorship(ctx,
			ctx.BlockHeader().ProposerAddress,
			ctx.BlockHeight())
		tags = tags.AppendTags(proposalCensorshipTag)
	}
	return
}
