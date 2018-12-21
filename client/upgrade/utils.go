package upgrade

import (
	protocol "github.com/irisnet/irishub/app/protocol/keeper"
	"github.com/irisnet/irishub/modules/upgrade"
)

type UpgradeInfoOutput struct {
	AppVerion          upgrade.AppVersion     `json:"version"`
	LastFailureVersion uint64                 `json:"last_failure_version"`
	UpgradeConfig      protocol.UpgradeConfig `json:"upgrade_config"`
}

func ConvertUpgradeInfoToUpgradeOutput(appVersion upgrade.AppVersion, upgradeConfig protocol.UpgradeConfig, lastFailureVersion uint64) UpgradeInfoOutput {

	return UpgradeInfoOutput{
		AppVerion:          appVersion,
		LastFailureVersion: lastFailureVersion,
		UpgradeConfig:      upgradeConfig,
	}
}
