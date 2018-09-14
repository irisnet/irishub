package lcd

import (
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/modules/upgrade"
	irisVersion "github.com/irisnet/irishub/version"
	"net/http"
)

type VersionInfo struct {
	IrisVersion    string `json:"iris_version"`
	UpgradeVersion int64  `json:"upgrade_version"`
	StartHeight    int64  `json:"start_height"`
	ProposalId     int64  `json:"proposal_id"`
}

func VersionHandlerFn(ctx context.CLIContext, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res_versionID, _ := ctx.QueryStore(upgrade.GetCurrentVersionKey(), "upgrade")
		var versionID int64
		cdc.MustUnmarshalBinary(res_versionID, &versionID)

		res_version, _ := ctx.QueryStore(upgrade.GetVersionIDKey(versionID), "upgrade")
		var version upgrade.Version
		cdc.MustUnmarshalBinary(res_version, &version)

		versionInfo := VersionInfo{
			IrisVersion:    irisVersion.Version,
			UpgradeVersion: version.Id,
			StartHeight:    version.Start,
			ProposalId:     version.ProposalID,
		}

		output, err := cdc.MarshalJSONIndent(versionInfo, "", "  ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}
