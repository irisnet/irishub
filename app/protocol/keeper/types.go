package keeper

import "github.com/irisnet/irishub/types/common"

const AppVersionTag = "app_version"

type UpgradeConfig struct {
	ProposalID uint64
	Definition common.ProtocolDefinition
}
