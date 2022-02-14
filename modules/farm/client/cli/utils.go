package cli

import (
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/irisnet/irismod/modules/farm/types"
)

// ParseCommunityPoolCreateFarmProposalWithDeposit reads and parses a CommunityPoolCreateFarmProposalWithDeposit from a file.
func ParseCommunityPoolCreateFarmProposalWithDeposit(cdc codec.JSONCodec, proposalFile string) (types.CommunityPoolCreateFarmProposalWithDeposit, error) {
	proposal := types.CommunityPoolCreateFarmProposalWithDeposit{}

	contents, err := os.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err = cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
