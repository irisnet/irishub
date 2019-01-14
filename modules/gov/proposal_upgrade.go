package gov

var _ Proposal = (*SoftwareUpgradeProposal)(nil)

type Upgrade struct {
	Version      uint64		`json:"version"`
	Software     string		`json:"software"`
	SwitchHeight uint64		`json:"switch_height"`
}

type SoftwareUpgradeProposal struct {
	TextProposal
	Upgrade	Upgrade  `json:"upgrade"`
}

func (sp SoftwareUpgradeProposal) GetUpgrade() Upgrade { return sp.Upgrade }
func (sp *SoftwareUpgradeProposal) SetUpgrade(upgrade Upgrade) {
	sp.Upgrade = upgrade
}
