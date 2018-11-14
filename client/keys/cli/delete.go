package keys

import (
	"fmt"

	"github.com/irisnet/irishub/client/keys"
	"github.com/spf13/cobra"
)

func deleteKeyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <name>",
		Short:   "Delete the given key",
		Example: "iriscli keys delete <key name>",
		RunE:    runDeleteCmd,
		Args:    cobra.ExactArgs(1),
	}
	return cmd
}

func runDeleteCmd(cmd *cobra.Command, args []string) error {
	name := args[0]

	kb, err := keys.GetKeyBaseWithWritePerm()
	if err != nil {
		return err
	}

	_, err = kb.Get(name)
	if err != nil {
		return err
	}

	buf := keys.BufferStdin()
	oldpass, err := keys.GetPassword(
		"DANGER - enter password to permanently delete key:", buf)
	if err != nil {
		return err
	}

	err = kb.Delete(name, oldpass)
	if err != nil {
		return err
	}
	fmt.Println("Password deleted forever (uh oh!)")
	return nil
}
