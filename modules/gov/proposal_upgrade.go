package gov

import sdk "github.com/irisnet/irishub/types"
var _ Proposal = (*SoftwareUpgradeProposal)(nil)

type Upgrade struct {
	Version      uint64		`json:"version"`
	Software     string		`json:"software"`
	SwitchHeight uint64		`json:"switch_height"`
	Threshold    sdk.Dec	`json:"threshold"`
}

type SoftwareUpgradeProposal struct {
	TextProposal
	Upgrade	Upgrade  `json:"upgrade"`
}

func (sp SoftwareUpgradeProposal) GetUpgrade() Upgrade { return sp.Upgrade }
func (sp *SoftwareUpgradeProposal) SetUpgrade(upgrade Upgrade) {
	sp.Upgrade = upgrade
}
