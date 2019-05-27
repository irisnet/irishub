package gov

var _ Proposal = (*PlainTextProposal)(nil)

type PlainTextProposal struct {
	BasicProposal
}
