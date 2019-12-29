package keys

import (
	"bufio"
	"errors"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys/mintkey"
	"github.com/irisnet/irishub/crypto/keystore"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const flagInteractive = "interactive"

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

	interactive := viper.GetBool(flagInteractive)

	decryptPassword, err := input.GetPassword("Enter passphrase to decrypt your key:", buf)
	if err != nil {
		return err
	}

	encryptPassword, err := input.GetPassword("Enter passphrase to encrypt the exported key:", buf)
	if err != nil {
		return err
	}

	if interactive {
		encryptPassword, err = input.GetString(
			"Enter passphrase to encrypt the imported key. "+
				"Most users should just hit enter to use the default, \"\"", buf)
		if err != nil {
			return err
		}

		// if they use one, make them re-enter it
		if len(encryptPassword) != 0 {
			p2, err := input.GetString("Repeat the passphrase:", buf)
			if err != nil {
				return err
			}

			if encryptPassword != p2 {
				return errors.New("passphrases don't match")
			}
		}
	}

	km, err := keystore.NewKeyStoreKeyManager(args[1], decryptPassword)
	if err != nil {
		return err
	}

	privArmor := mintkey.EncryptArmorPrivKey(km.GetPrivKey(), encryptPassword)

	if err := kb.Import(args[0], privArmor); err != nil {
		return err
	}
	return nil
}
