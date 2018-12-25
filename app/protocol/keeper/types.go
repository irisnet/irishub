package keeper

import (
	sdk "github.com/irisnet/irishub/types"
)

const AppVersionTag = "app_version"

type UpgradeConfig struct {
	ProposalID uint64
	Definition sdk.ProtocolDefinition
}
