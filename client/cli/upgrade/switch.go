package upgrade

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authctx "github.com/cosmos/cosmos-sdk/x/auth/client/context"
	"github.com/irisnet/irishub/modules/upgrade"
	"github.com/pkg/errors"
	"fmt"
	"os"
	"github.com/cosmos/cosmos-sdk/client/utils"
)

const (
	flagProposalID = "proposalID"
	flagTitle      = "title"
	flagVoter      = "voter"
)

// submit switch msg
func GetCmdSubmitSwitch(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-switch",
		Short: "Submit a switch msg for a upgrade propsal",
		RunE: func(cmd *cobra.Command, args []string) error {
			title := viper.GetString(flagTitle)
			proposalID := viper.GetInt64(flagProposalID)

			txCtx := authctx.NewTxContextFromCLI().WithCodec(cdc)
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			// get the from/to address
			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			msg := upgrade.NewMsgSwitch(title, proposalID, from)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			cliCtx.PrintResponse = true
			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagTitle, "", "title of switch")
	cmd.Flags().String(flagProposalID, "", "proposalID of upgrade proposal")

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
				return errors.Errorf("proposalID [%d] is not existed", proposalID)
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