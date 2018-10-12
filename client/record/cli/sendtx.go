package cli

import (
	"os"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	ipfs "github.com/ipfs/go-ipfs-api"
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
		Short: "Submit a transaction with a file hash",
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := viper.GetString(flagFilename)
			description := viper.GetString(flagDescription)

			strFilepath := viper.GetString(flagPath)
			file, err := os.Stat(strFilepath)

			if err != nil {
				// file does not exist
				return err
			}

			//upload to ipfs
			sh := ipfs.NewShell("localhost:5001")
			f, err := os.Open(strFilepath)
			if err != nil {
				return err
			}

			dataHash, err := sh.Add(f)

			if err != nil {
				return err
			}

			//file size
			dataSize := file.Size()
			//pinedNode

			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			txCtx := context.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			fromAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			submitTime := time.Now().Unix()

			msg := record.NewMsgSubmitFile(filename,
				strFilepath,
				description,
				submitTime,
				fromAddr,
				dataHash,
				dataSize)

			if cliCtx.GenerateOnly {
				return utils.PrintUnsignedStdTx(txCtx, cliCtx, []sdk.Msg{msg})
			}

			// Build and sign the transaction, then broadcast to Tendermint
			cliCtx.PrintResponse = true

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagFilename, "", "name of file")
	cmd.Flags().String(flagDescription, "", "description of file")
	cmd.Flags().String(flagPath, "", "full path of file (include filename)")

	return cmd
}
