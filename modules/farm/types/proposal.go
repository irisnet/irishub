package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	ProposalTypeCreateFarmPool string = "CommunityPoolCreateFarm"
)

// Implements Proposal Interface
var _ govv1beta1.Content = &CommunityPoolCreateFarmProposal{}

func init() {
	govv1beta1.RegisterProposalType(ProposalTypeCreateFarmPool)
}

func (cfp *CommunityPoolCreateFarmProposal) GetTitle() string       { return cfp.Title }
func (cfp *CommunityPoolCreateFarmProposal) GetDescription() string { return cfp.Description }
func (cfp *CommunityPoolCreateFarmProposal) ProposalRoute() string  { return RouterKey }
func (cfp *CommunityPoolCreateFarmProposal) ProposalType() string   { return ProposalTypeCreateFarmPool }
func (cfp *CommunityPoolCreateFarmProposal) ValidateBasic() error {
	// Validate gov base proposal
	if err := govv1beta1.ValidateAbstract(cfp); err != nil {
		return err
	}
	if err := ValidateDescription(cfp.PoolDescription); err != nil {
		return err
	}

	if err := ValidateLpTokenDenom(cfp.LptDenom); err != nil {
		return err
	}

	if err := ValidateCoins("RewardsPerBlock", cfp.RewardPerBlock...); err != nil {
		return err
	}
	return ValidateFund(cfp.RewardPerBlock, cfp.FundApplied, cfp.FundSelfBond)
}

func (cfp CommunityPoolCreateFarmProposal) String() string {
	return fmt.Sprintf(`Community Pool Create Farm Proposal:
  Title:       %s
  Description: %s
  PoolDescription: %s
  LpTokenDenom: %s
  RewardPerBlock: %s
  FundApplied: %s
  FundSelfBond: %s
`, cfp.Title,
		cfp.Description,
		cfp.PoolDescription,
		cfp.LptDenom,
		sdk.Coins(cfp.RewardPerBlock),
		sdk.Coins(cfp.FundApplied),
		sdk.Coins(cfp.FundSelfBond),
	)
}
