package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/irisnet/irishub/modules/upgrade"
	irisVersion "github.com/irisnet/irishub/version"
	"github.com/spf13/cobra"
)

func GetCmdVersion(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show full node running version info",
		RunE: func(cmd *cobra.Command, args []string) error {

			fmt.Printf("v%s\n", irisVersion.Version)

			ctx := context.NewCLIContext()

			var res_versionID []byte
			var err error
			res_versionID, err = ctx.QueryStore(upgrade.GetCurrentVersionKey(), storeName)
			if err != nil {
				return nil
			}
			var versionID int64
			cdc.MustUnmarshalBinary(res_versionID, &versionID)

			var res_version []byte
			res_version, err = ctx.QueryStore(upgrade.GetVersionIDKey(versionID), storeName)
			if err != nil {
				return nil
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
