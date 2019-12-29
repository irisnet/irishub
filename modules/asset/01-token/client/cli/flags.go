// nolint
package cli

import (
	flag "github.com/spf13/pflag"

	"github.com/irisnet/irishub/modules/asset/01-token/internal/types"
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
	FlagRecipient     = "recipient"
	FlagAmount        = "amount"
)

// common flagsets to add to various functions
var (
	FsEditToken     = flag.NewFlagSet("", flag.ContinueOnError)
	FsTokenIssue    = flag.NewFlagSet("", flag.ContinueOnError)
	FsTokensQuery   = flag.NewFlagSet("", flag.ContinueOnError)
	FsTransferToken = flag.NewFlagSet("", flag.ContinueOnError)
	FsMintToken     = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsTokenIssue.String(FlagSymbol, "", "the token symbol. Once created, it cannot be modified")
	FsTokenIssue.String(FlagName, "", "the token name, e.g. IRIS Network")
	FsTokenIssue.Uint8(FlagScale, 0, "the token decimal. The maximum value is 18")
	FsTokenIssue.String(FlagMinUnit, "", "the token symbol minimum uint")
	FsTokenIssue.Uint64(FlagInitialSupply, 0, "the initial supply token of token")
	FsTokenIssue.Uint64(FlagMaxSupply, types.MaximumTokenMaxSupply, "the max supply of the token")
	FsTokenIssue.Bool(FlagMintable, false, "whether the token can be minted, default false")

	FsTokensQuery.String(FlagSymbol, "", "the token symbol. Once created, it cannot be modified")
	FsTokensQuery.String(FlagOwner, "", "the owner address to be queried")

	FsEditToken.String(FlagName, "[do-not-modify]", "the token name, e.g. IRIS Network")
	FsEditToken.Uint64(FlagMaxSupply, 0, "the max supply of token")
	FsEditToken.String(FlagMintable, "", "whether the token can be minted, default false")

	FsTransferToken.String(FlagRecipient, "", "the new owner")

	FsMintToken.String(FlagRecipient, "", "address of mint token to")
	FsMintToken.Uint64(FlagAmount, 0, "amount of mint token")
}
