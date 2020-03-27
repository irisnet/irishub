package cli

import (
	"github.com/irisnet/irishub/app/v3/asset"
	flag "github.com/spf13/pflag"
)

const (
	FlagSymbol        = "symbol"
	FlagName          = "name"
	FlagScale         = "scale"
	FlagMinUnit       = "min-unit"
	FlagInitialSupply = "initial-supply"
	FlagMaxSupply     = "max-supply"
	FlagMintable      = "mintable"
	FlagOwner         = "owner"
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
	FsIssueToken.String(FlagSymbol, "", "the token symbol. Once created, it cannot be modified")
	FsIssueToken.String(FlagName, "", "the token name, e.g. IRIS Network")
	FsIssueToken.String(FlagMinUnit, "", "the minimum unit name of the token, e.g. wei")
	FsIssueToken.Uint8(FlagScale, 0, "the token decimal. The maximum value is 18")
	FsIssueToken.Uint64(FlagInitialSupply, 0, "the initial supply of the token")
	FsIssueToken.Uint64(FlagMaxSupply, asset.MaximumAssetMaxSupply, "the max supply of the token")
	FsIssueToken.Bool(FlagMintable, false, "whether the token can be minted, default to false")

	FsEditToken.String(FlagName, "[do-not-modify]", "the token name, e.g. IRIS Network")
	FsEditToken.Uint64(FlagMaxSupply, 0, "the max supply of the token")
	FsEditToken.String(FlagMintable, "", "whether the token can be minted, default to false")

	FsTransferTokenOwner.String(FlagTo, "", "the new owner")

	FsMintToken.String(FlagTo, "", "address of minting token to")
	FsMintToken.Uint64(FlagAmount, 0, "amount of minting token")
}
