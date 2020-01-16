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
	FlagTokenID       = "token-id"
	FlagTo            = "to"
	FlagAmount        = "amount"
)

var (
	FsEditToken          = flag.NewFlagSet("", flag.ContinueOnError)
	FsTokenIssue         = flag.NewFlagSet("", flag.ContinueOnError)
	FsTokensQuery        = flag.NewFlagSet("", flag.ContinueOnError)
	FsFeeQuery           = flag.NewFlagSet("", flag.ContinueOnError)
	FsTransferTokenOwner = flag.NewFlagSet("", flag.ContinueOnError)
	FsMintToken          = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsTokenIssue.String(FlagSymbol, "", "the token symbol. Once created, it cannot be modified")
	FsTokenIssue.String(FlagName, "", "the token name, e.g. IRIS Network")
	FsTokenIssue.String(FlagMinUnit, "", "the minimum unit name of token, e.g. wei")
	FsTokenIssue.Uint8(FlagScale, 0, "the token decimal. The maximum value is 18")
	FsTokenIssue.Uint64(FlagInitialSupply, 0, "the initial supply token of token")
	FsTokenIssue.Uint64(FlagMaxSupply, asset.MaximumAssetMaxSupply, "the max supply of the token")
	FsTokenIssue.Bool(FlagMintable, false, "whether the token can be minted, default false")

	FsTokensQuery.String(FlagOwner, "", "the owner address to be queried")
	FsTokensQuery.String(FlagTokenID, "", "The unique identifier of the token")

	FsFeeQuery.String(FlagSymbol, "", "the token symbol. Once created, it cannot be modified")

	FsEditToken.String(FlagName, "[do-not-modify]", "the token name, e.g. IRIS Network")
	FsEditToken.Uint64(FlagMaxSupply, 0, "the max supply of token")
	FsEditToken.String(FlagMintable, "", "whether the token can be minted, default false")

	FsTransferTokenOwner.String(FlagTo, "", "the new owner")

	FsMintToken.String(FlagTo, "", "address of mint token to")
	FsMintToken.Uint64(FlagAmount, 0, "amount of mint token")
}
