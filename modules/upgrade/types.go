package upgrade

import (
	protocol "github.com/irisnet/irishub/app/protocol/keeper"
	"github.com/irisnet/irishub/types/common"
)

type AppVersion struct {
	ProposalID uint64
	Success    bool
	Protocol   common.ProtocolDefinition
}

func NewVersion(upgradeConfig protocol.UpgradeConfig, success bool) AppVersion {
	return AppVersion{
		ProposalID: upgradeConfig.ProposalID,
		Success:    success,
		Protocol:   upgradeConfig.Definition,
	}
}
