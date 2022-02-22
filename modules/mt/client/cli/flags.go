package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagRecipient = "recipient"
	FlagData      = "data"
	FlagName      = "name"
	FlagMTID      = "mt-id"
	FlagAmount    = "amount"
)

var (
	FsIssueDenom    = flag.NewFlagSet("", flag.ContinueOnError)
	FsTransferDenom = flag.NewFlagSet("", flag.ContinueOnError)

	FsMintMT     = flag.NewFlagSet("", flag.ContinueOnError)
	FsEditMT     = flag.NewFlagSet("", flag.ContinueOnError)
	FsTransferMT = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsIssueDenom.String(FlagName, "", "The name of the denom")
	FsIssueDenom.String(FlagData, "", "The metadata of the denom")

	FsMintMT.String(FlagMTID, "", "The ID of the MT, leave empty when issuing, required when minting")
	FsMintMT.String(FlagAmount, "1", "The amount of the MT to mint")
	FsMintMT.String(FlagData, "", "The metadata of the MT")
	FsMintMT.String(FlagRecipient, "", "The recipient of the MT, default to the sender of the tx")

	FsEditMT.String(FlagData, "[do-not-modify]", "The metadata of the MT")
}
