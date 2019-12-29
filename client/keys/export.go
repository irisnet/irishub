package keys

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/client/keys"

	"github.com/irisnet/irishub/crypto/keystore"
)

const (
	flagOutfile = "output-file"
)

func ExportKeyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export <name>",
		Short: "Export an existing key",
		Long: `Export an existing key to the file,
You can import the file by keys add --recover --keystore=<file>.`,
		Example: "iriscli keys export <key name>",
		RunE:    runExportCmd,
		Args:    cobra.ExactArgs(1),
	}
	cmd.Flags().String(flagOutfile, "", "The keystore info will be written to the given file instead of STDOUT")
	cmd.Flags().Bool(flags.FlagIndentResponse, false, "Add indent to JSON response")
	return cmd
}

func runExportCmd(cmd *cobra.Command, args []string) error {
	buf := bufio.NewReader(cmd.InOrStdin())
	kb, err := keys.NewKeyringFromHomeFlag(buf)
	if err != nil {
		return err
	}

	_, err = kb.Get(args[0])
	if err != nil {
		return err
	}

	decryptPassword, err := input.GetPassword("Enter passphrase to decrypt your key:", buf)
	if err != nil {
		return err
	}

	encryptPassword, err := input.GetPassword("Enter passphrase to encrypt the exported key:", buf)
	if err != nil {
		return err
	}

	privKey, err := kb.ExportPrivateKeyObject(args[0], decryptPassword)
	if err != nil {
		return err
	}

	km := keystore.NewKeyManager(privKey)
	encryptedKeyJSON, err := km.ExportAsKeyStore(encryptPassword)
	if err != nil {
		return err
	}

	var jsonString []byte
	if viper.GetBool(flags.FlagIndentResponse) {
		jsonString, err = json.MarshalIndent(encryptedKeyJSON, "", "  ")
	} else {
		jsonString, err = json.Marshal(encryptedKeyJSON)
	}

	if viper.GetString(flagOutfile) == "" {
		fmt.Println(string(jsonString))
		return nil
	}

	fp, err := os.OpenFile(
		viper.GetString(flagOutfile), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644,
	)
	defer fp.Close()

	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(fp, "%s\n", jsonString); err != nil {
		return err
	}

	return nil
}
