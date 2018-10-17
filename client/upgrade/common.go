package upgrade

import (
	"github.com/irisnet/irishub/modules/upgrade"
)

type UpgradeInfoOutput struct {
	CurrentProposalId           int64           `json:"current_proposal_id"` //  proposalID of the proposal
	CurrentProposalAcceptHeight int64           `json:"current_proposal_accept_height"`
	Verion                      upgrade.Version `json:"version"`
}

func ConvertDepositToDepositOutput(version upgrade.Version, proposalId, hight int64) UpgradeInfoOutput {

	return UpgradeInfoOutput{
		CurrentProposalId:           proposalId,
		CurrentProposalAcceptHeight: hight,
		Verion: version,
	}
}
