package gov

var _ Proposal = (*SystemHaltProposal)(nil)

type SystemHaltProposal struct {
	TextProposal
}
