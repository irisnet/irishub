package keys

import (
	"github.com/spf13/cobra"
	"github.com/irisnet/irishub/client/keys/utils"
)

// CMD

// listKeysCmd represents the list command
var listKeysCmd = &cobra.Command{
	Use:   "list",
	Short: "List all keys",
	Long: `Return a list of all public keys stored by this key manager
along with their associated name and address.`,
	RunE: runListCmd,
}

func runListCmd(cmd *cobra.Command, args []string) error {
	kb, err := utils.GetKeyBase()
	if err != nil {
		return err
	}

	infos, err := kb.List()
	if err == nil {
		utils.PrintInfos(infos)
	}
	return err
}
