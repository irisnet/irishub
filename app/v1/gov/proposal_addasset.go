package gov

import "github.com/irisnet/irishub/app/v1/asset"

var _ Proposal = (*AddAssetProposal)(nil)

type AddAssetProposal struct {
	BasicProposal
	Assert asset.BaseAsset
}
