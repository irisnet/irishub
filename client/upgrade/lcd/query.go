package lcd

import (
	"github.com/irisnet/irishub/client/context"
	sdk "github.com/irisnet/irishub/types"
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

		res_currentVersion, _ := cliCtx.QueryStore(sdk.CurrentVersionKey, sdk.MainStore)
		var currentVersion uint64
		cdc.MustUnmarshalBinaryLengthPrefixed(res_currentVersion, &currentVersion)

		res_proposalID, _ := cliCtx.QueryStore(upgrade.GetSuccessVersionKey(currentVersion), storeName)
		var proposalID uint64
		cdc.MustUnmarshalBinaryLengthPrefixed(res_proposalID, &proposalID)

		res_currentVersionInfo, err := cliCtx.QueryStore(upgrade.GetProposalIDKey(proposalID), storeName)
		var currentVersionInfo upgrade.VersionInfo
		cdc.MustUnmarshalBinaryLengthPrefixed(res_currentVersionInfo, &currentVersionInfo)

		res_upgradeInProgress, _ := cliCtx.QueryStore(sdk.UpgradeConfigKey, sdk.MainStore)
		var upgradeInProgress sdk.UpgradeConfig
		if err == nil && len(res_upgradeInProgress) != 0 {
			cdc.MustUnmarshalBinaryLengthPrefixed(res_upgradeInProgress, &upgradeInProgress)
		}

		res_LastFailedVersion, err := cliCtx.QueryStore(sdk.LastFailedVersionKey, sdk.MainStore)
		var lastFailedVersion uint64
		if err == nil && len(res_LastFailedVersion) != 0 {
			cdc.MustUnmarshalBinaryLengthPrefixed(res_LastFailedVersion, &lastFailedVersion)
		} else {
			lastFailedVersion = 0
		}

		upgradeInfoOutput := upgcli.NewUpgradeInfoOutput(currentVersionInfo, lastFailedVersion, upgradeInProgress)

		output, err := cdc.MarshalJSONIndent(upgradeInfoOutput, "", "  ")
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Write(output)
	}
}
