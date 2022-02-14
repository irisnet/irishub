package types

import (
	"fmt"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeCreateFarmPool string = "CommunityPoolCreateFarm"
)

// Implements Proposal Interface
var _ govtypes.Content = &CommunityPoolCreateFarmProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeCreateFarmPool)
	govtypes.RegisterProposalTypeCodec(&CommunityPoolCreateFarmProposal{}, "irismod/CommunityPoolCreateFarmProposal")
}

func (cfp *CommunityPoolCreateFarmProposal) GetTitle() string       { return cfp.Title }
func (cfp *CommunityPoolCreateFarmProposal) GetDescription() string { return cfp.Description }
func (cfp *CommunityPoolCreateFarmProposal) ProposalRoute() string  { return RouterKey }
func (cfp *CommunityPoolCreateFarmProposal) ProposalType() string   { return ProposalTypeCreateFarmPool }
func (cfp *CommunityPoolCreateFarmProposal) ValidateBasic() error {
	//TODO
	return govtypes.ValidateAbstract(cfp)
}

func (cfp CommunityPoolCreateFarmProposal) String() string {
	return fmt.Sprintf(`Community Pool Create Farm Proposal:
  Title:       %s
  Description: %s
`, cfp.Title, cfp.Description)
}
