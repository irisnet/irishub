package cli

import (
	"fmt"
	"os"
	"path/filepath"
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
		Short: "Submit the specified file",
		RunE: func(cmd *cobra.Command, args []string) error {
			description := viper.GetString(flagDescription)
			strFilepath := viper.GetString(flagPath)
			strPinedNode := viper.GetString(flagPinedNode)

			_, filename := filepath.Split(strFilepath)

			var fileInfo os.FileInfo
			var err error
			if fileInfo, err = os.Stat(strFilepath); os.IsNotExist(err) {
				fmt.Printf("File %v doesn't exists, please check correstponding path.\n", strFilepath)
				return err
			}
			dataSize := fileInfo.Size()

			//upload to ipfs
			sh := ipfs.NewShell(strPinedNode)
			f, err := os.Open(strFilepath)
			if err != nil {
				return err
			}
			dataHash, err := sh.Add(f)
			if err != nil {
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

			submitTime := time.Now().Unix()

			msg := record.NewMsgSubmitFile(filename,
				strFilepath,
				description,
				submitTime,
				fromAddr,
				dataHash,
				dataSize,
				strPinedNode,
			)

			if cliCtx.GenerateOnly {
				return utils.PrintUnsignedStdTx(txCtx, cliCtx, []sdk.Msg{msg})
			}

			// Build and sign the transaction, then broadcast to Tendermint
			cliCtx.PrintResponse = true

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagDescription, "record file", "description of file")
	cmd.Flags().String(flagPath, "", "full path of file (include filename)")
	cmd.Flags().String(flagPinedNode, "localhost:5001", "node to upload file,ip:port")

	return cmd
}
