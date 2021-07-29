package cli

import (
	"fmt"

	flag "github.com/spf13/pflag"

	"github.com/irisnet/irismod/modules/token/types"
)

const (
	FlagSymbol        = "symbol"
	FlagName          = "name"
	FlagScale         = "scale"
	FlagMinUnit       = "min-unit"
	FlagInitialSupply = "initial-supply"
	FlagMaxSupply     = "max-supply"
	FlagMintable      = "mintable"
	FlagTo            = "to"
	FlagAmount        = "amount"
)

var (
	FsIssueToken         = flag.NewFlagSet("", flag.ContinueOnError)
	FsEditToken          = flag.NewFlagSet("", flag.ContinueOnError)
	FsTransferTokenOwner = flag.NewFlagSet("", flag.ContinueOnError)
	FsMintToken          = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsIssueToken.String(FlagSymbol, "", "The token symbol. Once created, it cannot be modified")
	FsIssueToken.String(FlagName, "", "The token name, e.g. IRIS Network")
	FsIssueToken.String(FlagMinUnit, "", "The minimum unit name of the token, e.g. wei")
	FsIssueToken.Uint32(FlagScale, 0, fmt.Sprintf("The token decimals, the maximum value is %d", types.MaximumScale))
	FsIssueToken.Uint64(FlagInitialSupply, 0, "The initial supply of the token")
	FsIssueToken.Uint64(FlagMaxSupply, types.MaximumMaxSupply, "The maximum supply of the token")
	FsIssueToken.Bool(FlagMintable, false, "Whether the token can be minted, default to false")

	FsEditToken.String(FlagName, "[do-not-modify]", "The token name, e.g. IRIS Network")
	FsEditToken.Uint64(FlagMaxSupply, 0, "The maximum supply of the token")
	FsEditToken.String(FlagMintable, "", "Whether the token can be minted, default to false")

	FsTransferTokenOwner.String(FlagTo, "", "The new owner")

	FsMintToken.String(FlagTo, "", "Address to which the token is to be minted")
	FsMintToken.Uint64(FlagAmount, 0, "Amount of the token to be minted")
}
