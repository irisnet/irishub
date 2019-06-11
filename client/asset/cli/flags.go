package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagOwner    = "owner"
	FlagMoniker  = "moniker"
	FlagIdentity = "identity"
	FlagDetails  = "details"
	FlagWebsite  = "website"
)

var (
	FsGatewayCreate = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsGatewayCreate.String(FlagMoniker, "", "the unique gateway name")
	FsGatewayCreate.String(FlagIdentity, "", "the gateway identity")
	FsGatewayCreate.String(FlagDetails, "", "the gateway description")
	FsGatewayCreate.String(FlagWebsite, "", "the external website")
}
