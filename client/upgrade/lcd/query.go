package lcd

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/modules/upgrade"
	"github.com/irisnet/irishub/modules/upgrade/params"
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

		res_height, _ := cliCtx.QueryStore([]byte("gov/"+upgradeparams.ProposalAcceptHeightParameter.GetStoreKey()), "params")
		res_proposalID, _ := cliCtx.QueryStore([]byte("gov/"+upgradeparams.CurrentUpgradeProposalIdParameter.GetStoreKey()), "params")
		var height int64
		var proposalID int64
		cdc.MustUnmarshalBinaryLengthPrefixed(res_height, &height)
		cdc.MustUnmarshalBinaryLengthPrefixed(res_proposalID, &proposalID)

		res_versionID, _ := cliCtx.QueryStore(upgrade.GetCurrentVersionKey(), storeName)
		var versionID int64
		cdc.MustUnmarshalBinaryLengthPrefixed(res_versionID, &versionID)

		res_version, _ := cliCtx.QueryStore(upgrade.GetVersionIDKey(versionID), storeName)
		var version upgrade.Version
		cdc.MustUnmarshalBinaryLengthPrefixed(res_version, &version)
		output, err := cdc.MarshalJSONIndent(version, "", "  ")
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Write(output)
	}
}
