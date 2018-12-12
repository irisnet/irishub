package version

import (
	"fmt"
	"github.com/irisnet/irishub/codec"
	"github.com/spf13/cobra"
)

// Version - Iris Version
const Version = "0.9.0-Alpha"
const ProtocolVersion = 0

// GitCommit set by build flags
var GitCommit = ""

// return version of CLI/node and commit hash
func GetVersion() string {
	v := Version
	if GitCommit != "" {
		v = v + "-" + GitCommit
	}
	return v
}

// ServeVersionCommand
func ServeVersionCommand(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show executable binary version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(GetVersion())
			return nil
		},
	}
	return cmd
}
