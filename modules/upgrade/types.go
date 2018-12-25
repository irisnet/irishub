package upgrade

import (
	protocol "github.com/irisnet/irishub/app/protocol/keeper"
	sdk "github.com/irisnet/irishub/types"
)

type AppVersion struct {
	ProposalID uint64
	Success    bool
	Protocol   sdk.ProtocolDefinition
}

func NewVersion(upgradeConfig protocol.UpgradeConfig, success bool) AppVersion {
	return AppVersion{
		ProposalID: upgradeConfig.ProposalID,
		Success:    success,
		Protocol:   upgradeConfig.Definition,
	}
}
