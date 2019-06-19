package gov

import (
	"fmt"
	"github.com/irisnet/irishub/app/v1/asset"
)

var _ Proposal = (*AddTokenProposal)(nil)

type AddTokenProposal struct {
	BasicProposal
	FToken asset.FungibleToken `json:"f_token"`
}

func (itp AddTokenProposal) String() string {
	return fmt.Sprintf(`Proposal %d:
  Title:              %s
  Type:               %s
  Status:             %s
  Submit Time:        %s
  Deposit End Time:   %s
  Total Deposit:      %s
  Voting Start Time:  %s
  Voting End Time:    %s
  Description:        %s
  %s`,
		itp.ProposalID, itp.Title, itp.ProposalType,
		itp.Status, itp.SubmitTime, itp.DepositEndTime,
		itp.TotalDeposit.MainUnitString(), itp.VotingStartTime, itp.VotingEndTime, itp.GetDescription(), itp.FToken.String())
}
