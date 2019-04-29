package upgrade

import (
	"github.com/irisnet/irishub/modules/upgrade"
	sdk "github.com/irisnet/irishub/types"
	"fmt"
)

type UpgradeInfoOutput struct {
	CurrentVersion    upgrade.VersionInfo `json:"current_version"`
	LastFailedVersion uint64              `json:"last_failed_version"`
	UpgradeInProgress sdk.UpgradeConfig   `json:"upgrade_in_progress"`
}

func NewUpgradeInfoOutput(currentVersion upgrade.VersionInfo, lastFailedVersion uint64, upgradeInProgress sdk.UpgradeConfig) UpgradeInfoOutput {
	return UpgradeInfoOutput{
		currentVersion,
		lastFailedVersion,
		upgradeInProgress,
	}
}

func (p UpgradeInfoOutput) String() string {
	success := "fail"
	if p.CurrentVersion.Success {
		success = "success"
	}
	return fmt.Sprintf(`Upgrade Info:
  Current Version[%v]:     %s     
  Last Failed Version:     %v
  Upgrade In Progress:     %s`,
		success, p.CurrentVersion.UpgradeInfo, p.LastFailedVersion, p.UpgradeInProgress)
}
