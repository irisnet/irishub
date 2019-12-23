// nolint
package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagAddress     = "address"
	FlagDescription = "description"
)

// common flagsets to add to various functions
var (
	FsAddGuardian    = flag.NewFlagSet("", flag.ContinueOnError)
	FsDeleteGuardian = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsAddGuardian.String(FlagAddress, "", "bech32 encoded account address")
	FsAddGuardian.String(FlagDescription, "", "description of account")
	FsDeleteGuardian.String(FlagAddress, "", "bech32 encoded account address")
}
