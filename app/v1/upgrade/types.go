package upgrade

import (
	sdk "github.com/irisnet/irishub/types"
)

type VersionInfo struct {
	UpgradeInfo sdk.UpgradeConfig
	Success     bool
}

func NewVersionInfo(upgradeConfig sdk.UpgradeConfig, success bool) VersionInfo {
	return VersionInfo{
		upgradeConfig,
		success,
	}
}
