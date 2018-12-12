package upgrade

import (
	"github.com/irisnet/irishub/modules/upgrade"
	protocol "github.com/irisnet/irishub/app/protocol/keeper"
)

type UpgradeInfoOutput struct {
	AppVerion  upgrade.AppVersion `json:"version"`
	UpgradeConfig protocol.UpgradeConfig `json:"upgrade_config"`
}

func ConvertUpgradeInfoToUpgradeOutput(appVersion upgrade.AppVersion, upgradeConfig protocol.UpgradeConfig) UpgradeInfoOutput {

	return UpgradeInfoOutput{
        AppVerion:appVersion,
        UpgradeConfig:upgradeConfig,
	}
}
