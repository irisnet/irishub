package cli

import (
	"fmt"
	"io/ioutil"

	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/modules/auth"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
)

const (
	flagAppend    = "append"
	flagPrintSigs = "print-sigs"
	flagOffline   = "offline"
)

// GetSignCommand returns the sign command
func GetSignCommand(codec *amino.Codec, decoder auth.AccountDecoder) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sign <file>",
		Short: "Sign transactions generated offline",
		Long: `Sign transactions created with the --generate-only flag.
Read a transaction from <file>, sign it, and print its JSON encoding.

The --offline flag makes sure that the client will not reach out to the local cache.
Thus account number or sequence number lookups will not be performed and it is
recommended to set such parameters manually.`,
		Example: "iriscli bank sign <file> --name <key name> --chain-id=<chain-id>",
		RunE:    makeSignCmd(codec, decoder),
		Args:    cobra.ExactArgs(1),
	}
	cmd.Flags().String(client.FlagName, "", "Name of private key with which to sign")
	cmd.Flags().Bool(flagAppend, true, "Append the signature to the existing ones. If disabled, old signatures would be overwritten")
	cmd.Flags().Bool(flagPrintSigs, false, "Print the addresses that must sign the transaction and those who have already signed it, then exit")
	cmd.Flags().Bool(flagOffline, false, "Offline mode. Do not query local cache.")
	cmd.MarkFlagRequired(client.FlagChainID)
	return cmd
}

func makeSignCmd(cdc *amino.Codec, decoder auth.AccountDecoder) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) (err error) {
		stdTx, err := readAndUnmarshalStdTx(cdc, args[0])
		if err != nil {
			return
		}

		if viper.GetBool(flagPrintSigs) {
			printSignatures(stdTx)
			return nil
		}

		name := viper.GetString(client.FlagName)
		cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(decoder)
		txCtx := utils.NewTxContextFromCLI()

		newTx, err := utils.SignStdTx(txCtx, cliCtx, name, stdTx, viper.GetBool(flagAppend), viper.GetBool(flagOffline))
		if err != nil {
			return err
		}
		var json []byte
		if cliCtx.Indent {
			json, err = cdc.MarshalJSONIndent(newTx, "", "  ")
		} else {
			json, err = cdc.MarshalJSON(newTx)
		}
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", json)
		return
	}
}

func printSignatures(stdTx auth.StdTx) {
	fmt.Println("Signers:")
	for i, signer := range stdTx.GetSigners() {
		fmt.Printf(" %v: %v\n", i, signer.String())
	}
	fmt.Println("")
	fmt.Println("Signatures:")
	for i, sig := range stdTx.GetSignatures() {
		fmt.Printf(" %v: %v\n", i, sdk.AccAddress(sig.Address()).String())
	}
	return
}

func readAndUnmarshalStdTx(cdc *amino.Codec, filename string) (stdTx auth.StdTx, err error) {
	var bytes []byte
	if bytes, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	if err = cdc.UnmarshalJSON(bytes, &stdTx); err != nil {
		return
	}
	return
}
