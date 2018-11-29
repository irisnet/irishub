package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// VersionCmd prints out the current sdk version
	VersionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the app version",
		Run:   printVersion,
	}
)

// CMD
func printVersion(cmd *cobra.Command, args []string) {
	fmt.Println(GetVersion())
}
