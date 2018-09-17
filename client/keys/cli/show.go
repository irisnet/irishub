package keys

import (
	"github.com/spf13/cobra"
	"github.com/irisnet/irishub/client/keys/utils"
)

var showKeysCmd = &cobra.Command{
	Use:   "show <name>",
	Short: "Show key info for the given name",
	Long:  `Return public details of one local key.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		info, err := utils.GetKey(name)
		if err == nil {
			utils.PrintInfo(info)
		}
		return err
	},
}


