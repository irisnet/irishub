package version

import (
	"fmt"
	"strconv"

	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
)

// Version - Iris Version
const ProtocolVersion = 2
const Version = "0.16.0"

// GitCommit set by build flags
var GitCommit = ""

// return version of CLI/node and commit hash
func GetVersion() string {
	v := Version
	if GitCommit != "" {
		v = v + "-" + GitCommit + "-" + strconv.Itoa(ProtocolVersion) + "-" + types.NetworkType
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
