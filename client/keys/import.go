package keys

import (
	"bufio"

	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys/mintkey"
	"github.com/irisnet/irishub/crypto/keystore"
	"github.com/spf13/cobra"
)

// ImportKeyCommand imports private keys from a keyfile.
func ImportKeyCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "import <name> <keyfile>",
		Short: "Import private keys into the local keybase",
		Long:  "Import a ASCII armored private key into the local keybase.",
		Args:  cobra.ExactArgs(2),
		RunE:  runImportCmd,
	}
}

func runImportCmd(cmd *cobra.Command, args []string) error {
	buf := bufio.NewReader(cmd.InOrStdin())
	kb, err := keys.NewKeyringFromHomeFlag(buf)
	if err != nil {
		return err
	}

	decryptPassword, err := input.GetPassword("Enter passphrase to decrypt your key:", buf)
	if err != nil {
		return err
	}

	km, err := keystore.NewKeyStoreKeyManager(args[1], decryptPassword)
	if err != nil {
		return err
	}

	privArmor := mintkey.EncryptArmorPrivKey(km.GetPrivKey(), "")

	if err := kb.ImportPrivKey(args[0], privArmor, ""); err != nil {
		return err
	}
	return nil
}
