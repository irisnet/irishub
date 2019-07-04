package gov

import (
	"fmt"
	"github.com/irisnet/irishub/app/v1/asset"
	sdk "github.com/irisnet/irishub/types"
)

var _ Proposal = (*AddTokenProposal)(nil)

type AddTokenProposal struct {
	BasicProposal
	FToken asset.FungibleToken `json:"f_token"`
}

func (atp AddTokenProposal) String() string {
	bps := atp.BasicProposal.String()
	return fmt.Sprintf(`%s
  %s`,
		bps, atp.FToken.String())
}

func (atp *AddTokenProposal) Validate(ctx sdk.Context, k Keeper) sdk.Error {
	if err := atp.BasicProposal.Validate(ctx, k); err != nil {
		return err
	}

	tokenId := atp.FToken.GetUniqueID()
	if k.ak.HasToken(ctx, tokenId) {
		return asset.ErrAssetAlreadyExists(k.codespace, tokenId)
	}
	return nil
}

func (atp *AddTokenProposal) Execute(ctx sdk.Context, gk Keeper) sdk.Error {
	logger := ctx.Logger()
	if err := atp.Validate(ctx, gk); err != nil {
		logger.Error("Execute AddTokenProposal failed", "height", ctx.BlockHeight(), "proposalId", atp.ProposalID, "token_id", atp.FToken.Id, "err", err.Error())
		return err
	}
	_, err := gk.ak.IssueToken(ctx, atp.FToken)
	if err != nil {
		logger.Error("Execute AddTokenProposal failed", "height", ctx.BlockHeight(), "proposalId", atp.ProposalID, "token_id", atp.FToken.Id, "err", err.Error())
		return err
	}
	logger.Info("Execute AddTokenProposal success", "height", ctx.BlockHeight(), "proposalId", atp.ProposalID, "token_id", atp.FToken.Id)
	return nil
}
