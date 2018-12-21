package gov

var _ Proposal = (*SoftwareUpgradeProposal)(nil)

type SoftwareUpgradeProposal struct {
	TextProposal
	Version      uint64
	Software     string
	SwitchHeight uint64
}
