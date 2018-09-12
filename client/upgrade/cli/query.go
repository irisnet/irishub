package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/irisnet/irishub/modules/upgrade"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func GetCmdInfo(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "query the information of upgrade module",
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			res_height, _ := cliCtx.QueryStore([]byte("gov/"+upgrade.GetCurrentProposalAcceptHeightKey()), "params")
			res_proposalID, _ := cliCtx.QueryStore([]byte("gov/"+upgrade.GetCurrentProposalIdKey()), "params")
			var height int64
			var proposalID int64
			cdc.MustUnmarshalBinary(res_height, &height)
			cdc.MustUnmarshalBinary(res_proposalID, &proposalID)

			res_versionID, _ := cliCtx.QueryStore(upgrade.GetCurrentVersionKey(), storeName)
			var versionID int64
			cdc.MustUnmarshalBinary(res_versionID, &versionID)

			res_version, _ := cliCtx.QueryStore(upgrade.GetVersionIDKey(versionID), storeName)
			var version upgrade.Version
			cdc.MustUnmarshalBinary(res_version, &version)
			output, err := wire.MarshalJSONIndent(cdc, version)
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			fmt.Println("CurrentProposalId           = ", proposalID)
			fmt.Println("CurrentProposalAcceptHeight = ", height)
			return nil
		},
	}
	return cmd
}

// Command to Get a Switch Information
func GetCmdQuerySwitch(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-switch",
		Short: "query switch details",
		RunE: func(cmd *cobra.Command, args []string) error {
			proposalID := viper.GetInt64(flagProposalID)
			voterStr := viper.GetString(flagVoter)

			voter, err := sdk.AccAddressFromBech32(voterStr)
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			res, err := cliCtx.QueryStore(upgrade.GetSwitchKey(proposalID, voter), storeName)
			if len(res) == 0 || err != nil {
				return fmt.Errorf("proposalID [%d] is not existed", proposalID)
			}

			var switchMsg upgrade.MsgSwitch
			cdc.MustUnmarshalBinary(res, &switchMsg)
			output, err := wire.MarshalJSONIndent(cdc, switchMsg)
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}

	cmd.Flags().String(flagProposalID, "", "proposalID of upgrade swtich being queried")
	cmd.Flags().String(flagVoter, "", "Address sign the switch msg")

	return cmd
}
