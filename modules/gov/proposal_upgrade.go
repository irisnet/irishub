package gov

import sdk "github.com/irisnet/irishub/types"
var _ Proposal = (*SoftwareUpgradeProposal)(nil)

type SoftwareUpgradeProposal struct {
	TextProposal
	Version      uint64
	Software     string
	SwitchHeight uint64
	Threshold    sdk.Dec
}
