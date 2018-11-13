package upgrade

import (
	"github.com/irisnet/irishub/modules/upgrade"
)

type UpgradeInfoOutput struct {
	CurrentProposalId           uint64           `json:"current_proposal_id"` //  proposalID of the proposal
	CurrentProposalAcceptHeight int64           `json:"current_proposal_accept_height"`
	Verion                      upgrade.Version `json:"version"`
}

func ConvertUpgradeInfoToUpgradeOutput(version upgrade.Version, proposalId uint64, hight int64) UpgradeInfoOutput {

	return UpgradeInfoOutput{
		CurrentProposalId:           proposalId,
		CurrentProposalAcceptHeight: hight,
		Verion: version,
	}
}
