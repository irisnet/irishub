package cli

import (
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/modules/record"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetCmdSubmitFile implements submitting upload file transaction command.
func GetCmdSubmitFile(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit",
		Short: "Submit a proposal with a file",
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := viper.GetString(flagFilename)
			description := viper.GetString(flagDescription)
			strProposalType := viper.GetString(flagProposalType)
			strAmount := viper.GetString(flagAmount)
			//todo upload to ipfs
			strFilepath := viper.GetString(flagPath)
			if _, err := os.Stat(strFilepath); os.IsNotExist(err) {
				// file does not exist
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			txCtx := context.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			fromAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			amount, err := cliCtx.ParseCoins(strAmount)
			if err != nil {
				return err
			}

			proposalType := strProposalType

			msg := record.NewMsgSubmitFile(filename, strFilepath, description, proposalType, fromAddr, amount)

			if cliCtx.GenerateOnly {
				return utils.PrintUnsignedStdTx(txCtx, cliCtx, []sdk.Msg{msg})
			}
			// Build and sign the transaction, then broadcast to Tendermint
			// proposalID must be returned, and it is a part of response.
			cliCtx.PrintResponse = true

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagFilename, "", "name of file")
	cmd.Flags().String(flagDescription, "", "description of file")
	cmd.Flags().String(flagPath, "", "full path of file (include filename)")

	return cmd
}
