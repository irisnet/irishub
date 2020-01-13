package upgrade

import (
	sdk "github.com/irisnet/irishub/types"
)

type VersionInfo struct {
	UpgradeInfo sdk.UpgradeConfig `json:"genesis_version"`
	Success     bool              `json:"success"`
}

func NewVersionInfo(upgradeConfig sdk.UpgradeConfig, success bool) VersionInfo {
	return VersionInfo{
		upgradeConfig,
		success,
	}
}
