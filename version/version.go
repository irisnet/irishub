package version

import (
	"fmt"

	"encoding/json"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/irisnet/irishub/modules/upgrade"
	"github.com/spf13/cobra"
	"net/http"
	"github.com/irisnet/irishub/client/context"
)

// Version - Iris Version
const Version = "0.4.0"

func GetCmdVersion(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version info",
		RunE: func(cmd *cobra.Command, args []string) error {

			fmt.Printf("v%s\n", Version)

			ctx := context.NewCLIContext()

			var res_versionID []byte
			var err error
			res_versionID, err = ctx.QueryStore(upgrade.GetCurrentVersionKey(), storeName)
			if err!=nil{
				return  nil
			}
			var versionID int64
			cdc.MustUnmarshalBinary(res_versionID, &versionID)

			var res_version []byte
			res_version, err  = ctx.QueryStore(upgrade.GetVersionIDKey(versionID), storeName)
			if err!=nil{
				return  nil
			}
			var version upgrade.Version
			cdc.MustUnmarshalBinary(res_version, &version)

			fmt.Println(version.Id)
			fmt.Println("Current version: Start Height    = ", version.Start)
			fmt.Println("Current version: Proposal Id     = ", version.ProposalID)
			return nil
		},
	}
	return cmd
}

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
			IrisVersion:    Version,
			UpgradeVersion: version.Id,
			StartHeight:    version.Start,
			ProposalId:     version.ProposalID,
		}

		output, err := json.MarshalIndent(versionInfo, "", "  ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}
