package lcd

import (
	protocol "github.com/irisnet/irishub/app/protocol/keeper"
	"github.com/irisnet/irishub/client/context"
	upgcli "github.com/irisnet/irishub/client/upgrade"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/upgrade"
	"net/http"
)

type VersionInfo struct {
	IrisVersion    string `json:"iris_version"`
	UpgradeVersion int64  `json:"upgrade_version"`
	StartHeight    int64  `json:"start_height"`
	ProposalId     int64  `json:"proposal_id"`
}

func InfoHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res_protocolVersion, _ := cliCtx.QueryStore(protocol.CurrentProtocolVersionKey, "protocol")
		var protocolVersion uint64
		cdc.UnmarshalJSON(res_protocolVersion, &protocolVersion)

		res_upgradeConfig, _ := cliCtx.QueryStore(protocol.UpgradeConfigkey, "protocol")
		var upgradeConfig protocol.UpgradeConfig
		cdc.UnmarshalJSON(res_upgradeConfig, &upgradeConfig)

		res_appVersion, _ := cliCtx.QueryStore(upgrade.GetAppVersionKey(protocolVersion), storeName)
		var appVersion upgrade.AppVersion
		cdc.MustUnmarshalBinaryLengthPrefixed(res_appVersion, &appVersion)

		upgradeInfoOutput := upgcli.ConvertUpgradeInfoToUpgradeOutput(appVersion, upgradeConfig)

		output, err := cdc.MarshalJSONIndent(upgradeInfoOutput, "", "  ")
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Write(output)
	}
}
