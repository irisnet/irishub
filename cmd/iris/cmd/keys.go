package cmd

import (
	"bufio"
	"encoding/json"
	"io/ioutil"

	"github.com/spf13/cobra"

	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/client/keys"

	"github.com/irisnet/irishub/keystore"
)

// Commands registers a sub-tree of commands to interact with
// local private key storage.
func Commands(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keys",
		Short: "Manage your application's keys",
		Long: `Keyring management commands. These keys may be in any format supported by the
Tendermint crypto library and can be used by light-clients, full nodes, or any other application
that needs to sign with a private key.

The keyring supports the following backends:

    os          Uses the operating system's default credentials store.
    file        Uses encrypted file-based keystore within the app's configuration directory.
                This keyring will request a password each time it is accessed, which may occur
                multiple times in a single command resulting in repeated password prompts.
    kwallet     Uses KDE Wallet Manager as a credentials management application.
    pass        Uses the pass command line utility to store and retrieve keys.
    test        Stores keys insecurely to disk. It does not prompt for a password to be unlocked
                and it should be use only for testing purposes.

kwallet and pass backends depend on external tools. Refer to their respective documentation for more
information:
    KWallet     https://github.com/KDE/kwallet
    pass        https://www.passwordstore.org/

The pass backend requires GnuPG: https://gnupg.org/
`,
	}

	cmd.AddCommand(
		keys.MnemonicKeyCommand(),
		keys.AddKeyCommand(),
		keys.ExportKeyCommand(),
		importKeyCommand(),
		keys.ListKeysCmd(),
		keys.ShowKeysCmd(),
		flags.LineBreak,
		keys.DeleteKeyCommand(),
		keys.ParseKeyStringCommand(),
		keys.MigrateCommand(),
	)

	cmd.PersistentFlags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	cmd.PersistentFlags().String(flags.FlagKeyringDir, "", "The client Keyring directory; if omitted, the default 'home' directory will be used")
	cmd.PersistentFlags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|test)")
	cmd.PersistentFlags().String(cli.OutputFlag, "text", "Output format (text|json)")

	return cmd
}

func importKeyCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "import <name> <keyfile>",
		Short: "Import private keys into the local keybase",
		Long:  "Import a ASCII armored private key into the local keybase.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			buf := bufio.NewReader(cmd.InOrStdin())
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			bz, err := ioutil.ReadFile(args[1])
			if err != nil {
				return err
			}

			passphrase, err := input.GetPassword("Enter passphrase to decrypt your key:", buf)
			if err != nil {
				return err
			}

			armor, err := getArmor(bz, passphrase)
			if err != nil {
				return err
			}
			return clientCtx.Keyring.ImportPrivKey(args[0], armor, passphrase)
		},
	}
}

func getArmor(privBytes []byte, passphrase string) (string, error) {
	if !json.Valid(privBytes) {
		return string(privBytes), nil
	}
	return keystore.RecoveryAndExportPrivKeyArmor(privBytes, passphrase)
}
