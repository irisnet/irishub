package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagIdentity      = "identity"
	FlagMoniker       = "moniker"
	FlagDetails       = "details"
	FlagWebsite       = "website"
	FlagRedeemAddress = "redeem-address"
	FlagOperators     = "operators"
)

var (
	FsGatewayCreate = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsGatewayCreate.String(FlagIdentity, "", "gateway identity")
	FsGatewayCreate.String(FlagMoniker, "", "gateway moniker")
	FsGatewayCreate.String(FlagDetails, "", "gateway description")
	FsGatewayCreate.String(FlagWebsite, "", "gateway website")
	FsGatewayCreate.String(FlagRedeemAddress, "", "the redeeming address of the gateway")
	FsGatewayCreate.StringSlice(FlagOperators, []string{}, "a set of addresses which are authorized to operate the gateway")
}
