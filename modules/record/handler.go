package record

import (
	//"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Handle all "record" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSubmitFile:
			return handleMsgSubmitFile(ctx, keeper, msg)
		default:
			errMsg := "Unrecognized record msg type"
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgSubmitFile(ctx sdk.Context, keeper Keeper, msg MsgSubmitFile) sdk.Result {

	err := msg.ValidateBasic()
	if err != nil {
		return err.Result()
	}

	//todo
	/*proposal := keeper.NewProposal(ctx, msg.Title, msg.Description, msg.ProposalType, msg.Param)

	err, votingStarted := keeper.AddDeposit(ctx, proposal.GetProposalID(), msg.Proposer, msg.Amount)
	if err != nil {
		return err.Result()
	}

	proposalIDBytes := keeper.cdc.MustMarshalBinaryBare(proposal.GetProposalID())*/

	return sdk.Result{
		//Data: proposalIDBytes,
	}
}
