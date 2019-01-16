package gov

var _ Proposal = (*HaltProposal)(nil)

type HaltProposal struct {
	TextProposal
}
