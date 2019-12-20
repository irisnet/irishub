package cli

import (
	flag "github.com/spf13/pflag"

	"github.com/irisnet/irishub/modules/asset/internal/types"
)

const (
	FlagFamily          = "family"
	FlagSource          = "source"
	FlagSymbol          = "symbol"
	FlagCanonicalSymbol = "canonical-symbol"
	FlagName            = "name"
	FlagDecimal         = "decimal"
	FlagMinUnitAlias    = "min-unit-alias"
	FlagInitialSupply   = "initial-supply"
	FlagMaxSupply       = "max-supply"
	FlagMintable        = "mintable"
	FlagOwner           = "owner"
	FlagTo              = "to"
	FlagAmount          = "amount"
)

var (
	FsEditToken          = flag.NewFlagSet("", flag.ContinueOnError)
	FsTokenIssue         = flag.NewFlagSet("", flag.ContinueOnError)
	FsTokensQuery        = flag.NewFlagSet("", flag.ContinueOnError)
	FsTransferTokenOwner = flag.NewFlagSet("", flag.ContinueOnError)
	FsMintToken          = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsTokenIssue.String(FlagFamily, "", "the asset family, valid values can be fungible and non-fungible")
	FsTokenIssue.String(FlagSource, "", "the asset source, valid values can be native, external and gateway")
	FsTokenIssue.String(FlagSymbol, "", "the token symbol. Once created, it cannot be modified")
	FsTokenIssue.String(FlagCanonicalSymbol, "", "the source symbol of a gateway or external token")
	FsTokenIssue.String(FlagName, "", "the token name, e.g. IRIS Network")
	FsTokenIssue.Uint8(FlagDecimal, 0, "the token decimal. The maximum value is 18")
	FsTokenIssue.String(FlagMinUnitAlias, "", "the token symbol minimum alias")
	FsTokenIssue.Uint64(FlagInitialSupply, 0, "the initial supply token of token")
	FsTokenIssue.Uint64(FlagMaxSupply, types.MaximumAssetMaxSupply, "the max supply of the token")
	FsTokenIssue.Bool(FlagMintable, false, "whether the token can be minted, default false")

	FsTokensQuery.String(FlagSource, "", "the asset source, valid values can be native, external and gateway")
	FsTokensQuery.String(FlagOwner, "", "the owner address to be queried")

	FsEditToken.String(FlagName, "[do-not-modify]", "the token name, e.g. IRIS Network")
	FsEditToken.String(FlagCanonicalSymbol, "[do-not-modify]", "the source symbol of a gateway or external token")
	FsEditToken.String(FlagMinUnitAlias, "[do-not-modify]", "the token symbol minimum alias")
	FsEditToken.Uint64(FlagMaxSupply, 0, "the max supply of token")
	FsEditToken.String(FlagMintable, "", "whether the token can be minted, default false")

	FsTransferTokenOwner.String(FlagTo, "", "the new owner")

	FsMintToken.String(FlagTo, "", "address of mint token to")
	FsMintToken.Uint64(FlagAmount, 0, "amount of mint token")
}
