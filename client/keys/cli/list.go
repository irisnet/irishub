package keys

import (
	"github.com/irisnet/irishub/client/keys"
	"github.com/spf13/cobra"
	"github.com/irisnet/irishub/client"
)

// CMD

// listKeysCmd represents the list command
func listKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all keys",
		Long: `Return a list of all public keys stored by this key manager
along with their associated name and address.`,
		RunE: runListCmd,
	}
	cmd.Flags().Bool(client.FlagIndentResponse, false, "Add indent to JSON response")
	return cmd
}

func runListCmd(cmd *cobra.Command, args []string) error {
	kb, err := keys.GetKeyBase()
	if err != nil {
		return err
	}

	infos, err := kb.List()
	if err == nil {
		keys.PrintInfos(cdc, infos)
	}
	return err
}
