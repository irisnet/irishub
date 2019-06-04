package gov

var _ Proposal = (*AddAssetProposal)(nil)

type AddAssetProposal struct {
	BasicProposal
}
