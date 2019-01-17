package version

import (
	"fmt"
	"github.com/irisnet/irishub/codec"
	"github.com/spf13/cobra"
	"strconv"
)

// Version - Iris Version
const ProtocolVersion = 1
const Version = "0.10.2"
// GitCommit set by build flags
var GitCommit = ""

// return version of CLI/node and commit hash
func GetVersion() string {
	v := Version
	if GitCommit != "" {
		v = v + "-" + GitCommit + "-" + strconv.Itoa(ProtocolVersion)
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
