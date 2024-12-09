package keepers

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

type DistrKeeperAdapter struct {
	keeper distrkeeper.Keeper
}

func NewDistrKeeperAdapter(keeper distrkeeper.Keeper) DistrKeeperAdapter {
	return DistrKeeperAdapter{keeper: keeper}
}

func (a DistrKeeperAdapter) GetFeePool(ctx context.Context) (distrtypes.FeePool, error) {
	return a.keeper.FeePool.Get(ctx)
}

func (a DistrKeeperAdapter) SetFeePool(ctx context.Context, feePool distrtypes.FeePool) error {
	return a.keeper.FeePool.Set(ctx, feePool)
}

type GovKeeperAdapter struct {
	keeper *govkeeper.Keeper
}

func NewGovKeeperAdapter(keeper *govkeeper.Keeper) GovKeeperAdapter {
	return GovKeeperAdapter{keeper: keeper}
}

func (a GovKeeperAdapter) SubmitProposal(ctx context.Context, messages []sdk.Msg, metadata, title, summary string, proposer sdk.AccAddress, expedited bool) (v1.Proposal, error) {
	return a.keeper.SubmitProposal(ctx, messages, metadata, title, summary, proposer, expedited)
}

func (a GovKeeperAdapter) AddDeposit(ctx context.Context, proposalID uint64, depositorAddr sdk.AccAddress, depositAmount sdk.Coins) (bool, error) {
	return a.keeper.AddDeposit(ctx, proposalID, depositorAddr, depositAmount)
}

func (a GovKeeperAdapter) GetProposal(ctx context.Context, proposalID uint64) (v1.Proposal, error) {
	return a.keeper.Proposals.Get(ctx, proposalID)
}

func (a GovKeeperAdapter) GetGovernanceAccount(ctx context.Context) sdk.ModuleAccountI {
	return a.keeper.GetGovernanceAccount(ctx)
}
