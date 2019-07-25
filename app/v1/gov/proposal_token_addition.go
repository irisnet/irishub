package gov

import (
	"fmt"
	"github.com/irisnet/irishub/app/v1/asset/exported"
	sdk "github.com/irisnet/irishub/types"
)

var _ Proposal = (*TokenAdditionProposal)(nil)

type TokenAdditionProposal struct {
	BasicProposal
	FToken exported.FungibleToken `json:"f_token"`
}

func (atp TokenAdditionProposal) HumanString(converter sdk.CoinsConverter) string {
	bps := atp.BasicProposal.HumanString(converter)
	return fmt.Sprintf(`%s
  %s`,
		bps, atp.FToken.String())
}

func (atp *TokenAdditionProposal) Validate(ctx sdk.Context, k Keeper, verify bool) sdk.Error {
	if err := atp.BasicProposal.Validate(ctx, k, verify); err != nil {
		return err
	}

	tokenId := atp.FToken.GetUniqueID()
	if k.ak.HasToken(ctx, tokenId) {
		return exported.ErrAssetAlreadyExists(k.codespace, fmt.Sprintf("token already exists: %s", tokenId))
	}
	return nil
}

func (atp *TokenAdditionProposal) Execute(ctx sdk.Context, gk Keeper) sdk.Error {
	logger := ctx.Logger()
	_, err := gk.ak.IssueToken(ctx, atp.FToken)
	if err != nil {
		logger.Error("Execute TokenAdditionProposal failed", "height", ctx.BlockHeight(), "proposalId", atp.ProposalID, "token_id", atp.FToken.Id, "err", err.Error())
		return err
	}
	logger.Info("Execute TokenAdditionProposal success", "height", ctx.BlockHeight(), "proposalId", atp.ProposalID, "token_id", atp.FToken.Id)
	return nil
}
