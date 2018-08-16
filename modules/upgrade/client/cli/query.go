package cli

import (
	"github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/irisnet/irishub/modules/upgrade"
	"fmt"
)

func GetCmdInfo(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "query the information of upgrade module",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx := context.NewCoreContextFromViper()
			res_height, _ := ctx.QueryStore([]byte("gov/"+upgrade.GetCurrentProposalAcceptHeightKey()), "params")
            res_proposalID, _ := ctx.QueryStore([]byte("gov/"+upgrade.GetCurrentProposalIdKey()),"params")
			var height int64
			var proposalID int64
			cdc.MustUnmarshalBinary(res_height, &height)
			cdc.MustUnmarshalBinary(res_proposalID, &proposalID)

			res_versionID,_ := ctx.QueryStore(upgrade.GetCurrentVersionKey(),storeName)
			var versionID int64
			cdc.MustUnmarshalBinary(res_versionID, &versionID)

			res_version,_:= ctx.QueryStore(upgrade.GetVersionIDKey(versionID),storeName)
			var version upgrade.Version
			cdc.MustUnmarshalBinary(res_version, &version)
			output, err := wire.MarshalJSONIndent(cdc, version)
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			fmt.Println("CurrentProposalId           = ",proposalID)
			fmt.Println("CurrentProposalAcceptHeight = ",height)
			return nil
		},
	}
	return cmd
}