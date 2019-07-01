package cli

import (
	"github.com/irisnet/irishub/app/v1/asset"
	flag "github.com/spf13/pflag"
)

const (
	FlagFamily         = "family"
	FlagSource         = "source"
	FlagGateway        = "gateway"
	FlagSymbol         = "symbol"
	FlagSymbolAtSource = "symbol-at-source"
	FlagName           = "name"
	FlagDecimal        = "decimal"
	FlagSymbolMinAlias = "symbol-min-alias"
	FlagInitialSupply  = "initial-supply"
	FlagMaxSupply      = "max-supply"
	FlagMintable       = "mintable"

	FlagOwner    = "owner"
	FlagMoniker  = "moniker"
	FlagIdentity = "identity"
	FlagDetails  = "details"
	FlagWebsite  = "website"
	FlagTo       = "to"

	FlagSubject = "subject"
	FlagID      = "id"
	FlagAmount  = "amount"
)

var (
	FsEditToken            = flag.NewFlagSet("", flag.ContinueOnError)
	FsTokenIssue           = flag.NewFlagSet("", flag.ContinueOnError)
	FsTokensQuery          = flag.NewFlagSet("", flag.ContinueOnError)
	FsGatewayCreate        = flag.NewFlagSet("", flag.ContinueOnError)
	FsGatewayEdit          = flag.NewFlagSet("", flag.ContinueOnError)
	FsGatewayOwnerTransfer = flag.NewFlagSet("", flag.ContinueOnError)
	FsFeeQuery             = flag.NewFlagSet("", flag.ContinueOnError)
	FsTransferTokenOwner   = flag.NewFlagSet("", flag.ContinueOnError)
	FsMintToken            = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsTokenIssue.String(FlagFamily, "", "the asset family, valid values can be fungible and non-fungible")
	FsTokenIssue.String(FlagSource, "", "the asset source, valid values can be native, external and gateway")
	FsTokenIssue.String(FlagGateway, "", "the gateway name of gateway token. required if --source=gateway")
	FsTokenIssue.String(FlagSymbol, "", "the token symbol. Once created, it cannot be modified")
	FsTokenIssue.String(FlagSymbolAtSource, "", "the source symbol of a gateway or external token")
	FsTokenIssue.String(FlagName, "", "the token name, e.g. IRIS Network")
	FsTokenIssue.Uint8(FlagDecimal, 0, "the token decimal. The maximum value is 18")
	FsTokenIssue.String(FlagSymbolMinAlias, "", "the token symbol minimum alias")
	FsTokenIssue.Uint64(FlagInitialSupply, 0, "the initial supply token of token")
	FsTokenIssue.Uint64(FlagMaxSupply, asset.MaximumAssetMaxSupply, "the max supply of the token")
	FsTokenIssue.Bool(FlagMintable, false, "whether the token can be minted, default false")

	FsTokensQuery.String(FlagSource, "", "the asset source, valid values can be native, external and gateway")
	FsTokensQuery.String(FlagGateway, "", "the gateway name of gateway token. required if --source=gateway")
	FsTokensQuery.String(FlagOwner, "", "the owner address to be queried")

	FsGatewayCreate.String(FlagMoniker, "", "the unique gateway name")
	FsGatewayCreate.String(FlagIdentity, "", "the gateway identity")
	FsGatewayCreate.String(FlagDetails, "", "the gateway description")
	FsGatewayCreate.String(FlagWebsite, "", "the external website")

	FsGatewayEdit.String(FlagMoniker, asset.DoNotModify, "the unique gateway name")
	FsGatewayEdit.String(FlagIdentity, asset.DoNotModify, "the gateway identity")
	FsGatewayEdit.String(FlagDetails, asset.DoNotModify, "the gateway description")
	FsGatewayEdit.String(FlagWebsite, asset.DoNotModify, "the external website")

	FsGatewayOwnerTransfer.String(FlagMoniker, "", "the unique name of the gateway to be transferred")
	FsGatewayOwnerTransfer.String(FlagTo, "", "the new owner")

	FsFeeQuery.String(FlagSubject, "", "the fee type to be queried")
	FsFeeQuery.String(FlagMoniker, "", "the gateway name, required if the subject is gateway")
	FsFeeQuery.String(FlagID, "", "the token id, required if the subject is token")

	FsEditToken.String(FlagName, "[do-not-modify]", "the token name, e.g. IRIS Network")
	FsEditToken.String(FlagSymbolAtSource, "[do-not-modify]", "the source symbol of a gateway or external token")
	FsEditToken.String(FlagSymbolMinAlias, "[do-not-modify]", "the token symbol minimum alias")
	FsEditToken.Uint64(FlagMaxSupply, 0, "the max supply of token")
	FsEditToken.Bool(FlagMintable, false, "whether the token can be minted, default false")

	FsTransferTokenOwner.String(FlagTo, "", "the new owner")

	FsMintToken.String(FlagTo, "", "address of mint token to")
	FsMintToken.Uint64(FlagAmount, 0, "amount of mint token")
}
