package version

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/spf13/cobra"
)

// Version - Iris Version
const Version = "0.4.0-GOG"

func GetCmdVersion(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version info",
		RunE: func(cmd *cobra.Command, args []string) error {

			fmt.Printf("v%s\n", Version)
			return nil
		},
	}
	return cmd
}

