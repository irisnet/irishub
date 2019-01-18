package slashing

import (
	"encoding/binary"
	"fmt"

	sdk "github.com/irisnet/irishub/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/crypto"
)

// slashing begin block functionality
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, sk Keeper) (tags sdk.Tags) {

	// Tag the height
	heightBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(heightBytes, uint64(req.Header.Height))
	tags = sdk.NewTags("height", heightBytes)

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
			doubleSignSlashTag := sk.handleDoubleSign(ctx, evidence.Validator.Address, evidence.Height, evidence.Time, evidence.Validator.Power)
			tags = tags.AppendTags(doubleSignSlashTag)
		default:
			ctx.Logger().With("module", "iris/slashing").Error(fmt.Sprintf("ignored unknown evidence type: %s", evidence.Type))
		}
	}

	return
}

// slashing end block functionality
func EndBlocker(ctx sdk.Context, req abci.RequestEndBlock, sk Keeper) (tags sdk.Tags) {

	// Tag the height
	heightBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(heightBytes, uint64(req.Height))
	tags = sdk.NewTags("height", heightBytes)

	if int64(ctx.CheckValidNum()) < ctx.BlockHeader().NumTxs {
		ctx.Logger().With("module", "iris/slashing").
			Info("the malefactor proposer proposed a invalid block",
				"proposer address", crypto.Address(ctx.BlockHeader().ProposerAddress).String(),
				"block height", ctx.BlockHeight())

		proposalCensorshipTag := sk.handleProposerCensorship(ctx,
			ctx.BlockHeader().ProposerAddress,
			ctx.BlockHeight())
		tags = tags.AppendTags(proposalCensorshipTag)
	}
	return
}
