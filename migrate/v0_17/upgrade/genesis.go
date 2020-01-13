package upgrade

import sdk "github.com/cosmos/cosmos-sdk/types"

type GenesisState struct {
	GenesisVersion VersionInfo `json:"genesis_version"`
}

type VersionInfo struct {
	UpgradeInfo UpgradeConfig `json:"upgrade_info"`
	Success     bool          `json:"success"`
}

type UpgradeConfig struct {
	ProposalID uint64             `json:"proposal_id"`
	Protocol   ProtocolDefinition `json:"proposal"`
}

type ProtocolDefinition struct {
	Version   uint64  `json:"version"`
	Software  string  `json:"software"`
	Height    uint64  `json:"height"`
	Threshold sdk.Dec `json:"threshold"`
}
