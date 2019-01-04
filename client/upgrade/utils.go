package upgrade

import (
	"github.com/irisnet/irishub/modules/upgrade"
	sdk "github.com/irisnet/irishub/types"
)

type UpgradeInfoOutput struct {
	CurrentVersion    upgrade.VersionInfo `json:"current_version"`
	LastFailedVersion int64               `json:"last_failed_version"`
	UpgradeInProgress sdk.UpgradeConfig   `json:"upgrade_in_progress"`
}

func NewUpgradeInfoOutput(currentVersion upgrade.VersionInfo, lastFailedVersion int64, upgradeInProgress sdk.UpgradeConfig) UpgradeInfoOutput {
	return UpgradeInfoOutput{
		currentVersion,
		lastFailedVersion,
		upgradeInProgress,
	}
}
