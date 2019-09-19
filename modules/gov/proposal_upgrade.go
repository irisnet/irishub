package gov

import sdk "github.com/irisnet/irishub/types"

var _ Proposal = (*SoftwareUpgradeProposal)(nil)

type SoftwareUpgradeProposal struct {
	BasicProposal
	ProtocolDefinition sdk.ProtocolDefinition `json:"protocol_definition"`
}

func (sp SoftwareUpgradeProposal) GetProtocolDefinition() sdk.ProtocolDefinition {
	return sp.ProtocolDefinition
}
func (sp *SoftwareUpgradeProposal) SetProtocolDefinition(upgrade sdk.ProtocolDefinition) {
	sp.ProtocolDefinition = upgrade
}
