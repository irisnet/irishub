package keys

import (
	"fmt"

	"bufio"
	"errors"
	"os"

	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/keys"
	ck "github.com/irisnet/irishub/crypto/keys"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagYes = "yes"
)

func deleteKeyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete the given key",
		Long: `Delete a key from the store.
Note that removing offline or ledger keys will remove
only the public key references stored locally, i.e.
private keys stored in a ledger device cannot be deleted with
iriscli.
`,
		Example: "iriscli keys delete <key name>",
		RunE:    runDeleteCmd,
		Args:    cobra.ExactArgs(1),
	}

	cmd.Flags().BoolP(flagYes, "y", false,
		"Skip confirmation prompt when deleting offline or ledger key references")
	return cmd
}

func runDeleteCmd(cmd *cobra.Command, args []string) error {
	name := args[0]

	kb, err := keys.GetKeyBaseWithWritePerm()
	if err != nil {
		return err
	}

	info, err := kb.Get(name)
	if err != nil {
		return err
	}

	buf := client.BufferStdin()
	if info.GetType() == ck.TypeLedger || info.GetType() == ck.TypeOffline {
		if !viper.GetBool(flagYes) {
			if err := confirmDeletion(buf); err != nil {
				return err
			}
		}
		if err := kb.Delete(name, ""); err != nil {
			return err
		}
		fmt.Fprintln(os.Stderr, "Public key reference deleted")
		return nil
	}

	oldpass, err := keys.GetPassword(
		"DANGER - enter password to permanently delete key:", buf)
	if err != nil {
		return err
	}

	err = kb.Delete(name, oldpass)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "Key deleted forever (uh oh!)")
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
