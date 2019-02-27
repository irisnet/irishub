package keys

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	client "github.com/irisnet/irishub/client/keys"
	"github.com/irisnet/irishub/crypto/keys"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagYes   = "yes"
	flagForce = "force"
)

func deleteKeyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <name>",
		Short:   "Delete the given key",
		Example: "iriscli keys delete <key name>",
		Long: `Delete a key from the store.
Note that removing offline or ledger keys will remove
only the public key references stored locally, i.e.
private keys stored in a ledger device cannot be deleted with
gaiacli.
`,
		RunE: runDeleteCmd,
		Args: cobra.ExactArgs(1),
	}

	cmd.Flags().BoolP(flagYes, "y", false,
		"Skip confirmation prompt when deleting offline or ledger key references")
	cmd.Flags().BoolP(flagForce, "f", false,
		"Remove the key unconditionally without asking for the passphrase")
	return cmd
}

func runDeleteCmd(cmd *cobra.Command, args []string) error {
	name := args[0]

	kb, err := client.GetKeyBaseWithWritePerm()
	if err != nil {
		return err
	}

	info, err := kb.Get(name)
	if err != nil {
		return err
	}

	buf := client.BufferStdin()
	if info.GetType() == keys.TypeLedger || info.GetType() == keys.TypeOffline {
		if !viper.GetBool(flagYes) {
			if err := confirmDeletion(buf); err != nil {
				return err
			}
		}
		if err := kb.Delete(name, "", true); err != nil {
			return err
		}
		fmt.Fprintln(os.Stderr, "Public key reference deleted")
		return nil
	}

	// skip passphrase check if run with --force
	skipPass := viper.GetBool(flagForce)
	var oldpass string
	if !skipPass {
		oldpass, err = client.GetPassword(
			"DANGER - enter password to permanently delete key:", buf)
		if err != nil {
			return err
		}
	}

	err = kb.Delete(name, oldpass, skipPass)
	if err != nil {
		return err
	}
	fmt.Println("Password deleted forever (uh oh!)")
	return nil
}

func confirmDeletion(buf *bufio.Reader) error {
	answer, err := client.GetConfirmation("Key reference will be deleted. Continue?", buf)
	if err != nil {
		return err
	}
	if !answer {
		return errors.New("aborted")
	}
	return nil
}
