// nolint
package tags

import (
	sdk "github.com/irisnet/irishub/types"
)

var (
	ActionSubmitProposal   = []byte("submit-proposal")
	ActionDeposit          = []byte("deposit")
	ActionVote             = []byte("vote")
	ActionProposalDropped  = []byte("proposal-dropped")
	ActionProposalPassed   = []byte("proposal-passed")
	ActionProposalRejected = []byte("proposal-rejected")

	Action            = sdk.TagAction
	Proposer          = "proposer"
	ProposalID        = "proposal-id"
	VotingPeriodStart = "voting-period-start"
	Depositor         = "depositor"
	Voter             = "voter"
	Param             = "param"
	Usage             = "usage"
	Percent           = "percent"
	DestAddress       = "dest-address"
)
