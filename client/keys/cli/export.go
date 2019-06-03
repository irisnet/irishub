package keys

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/keys"
	"github.com/irisnet/irishub/crypto/keystore"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagOutfile = "output-file"
)

func exportKeyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export <name>",
		Short: "Export an existing key",
		Long: `Add a public/private key pair to the key store.
If you select --seed/-s you can recover a key from the seed
phrase, otherwise, a new key will be generated.`,
		Example: "iriscli keys export <key name>",
		RunE:    runExportCmd,
		Args:    cobra.ExactArgs(1),
	}
	cmd.Flags().String(flagOutfile, "", "The key store will be written to the given file instead of STDOUT")
	cmd.Flags().Bool(client.FlagIndentResponse, false, "Add indent to JSON response")
	return cmd
}

func runExportCmd(_ *cobra.Command, args []string) error {
	var err error
	var name string

	if len(args) != 1 || len(args[0]) == 0 {
		return errors.New("you must provide a name for the key")
	}
	name = args[0]
	kb, err := keys.GetKeyBase()
	if err != nil {
		return err
	}

	_, err = kb.Get(name)
	if err != nil {
		return err
	}

	pass, err := keys.GetPassphrase(name)
	if err != nil {
		return err
	}

	privKey, err := kb.ExportPrivateKeyObject(name, pass)
	if err != nil {
		return err
	}

	passphrase, err := keys.ReadKeystorePassphraseFromStdin()
	if err != nil {
		return err
	}

	km := keystore.NewKeyManager(privKey)
	encryptedKeyJSON, err := km.ExportAsKeyStore(passphrase)
	if err != nil {
		return err
	}

	var jsonString []byte
	if viper.GetBool(client.FlagIndentResponse) {
		jsonString, err = json.MarshalIndent(encryptedKeyJSON, "", "  ")
	} else {
		jsonString, err = json.Marshal(encryptedKeyJSON)
	}

	if viper.GetString(flagOutfile) == "" {
		fmt.Printf("%s\n", jsonString)
		return nil
	}

	fp, err := os.OpenFile(
		viper.GetString(flagOutfile), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644,
	)
	defer fp.Close()

	if err != nil {
		return err
	}

	fmt.Fprintf(fp, "%s\n", jsonString)

	return nil
}
