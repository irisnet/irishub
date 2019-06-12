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
	FlagName           = "name"
	FlagDecimal        = "decimal"
	FlagSymbolMinAlias = "symbol-mi-alias"
	FlagInitialSupply  = "initial-supply"
	FlagMaxSupply      = "max-supply"
	FlagMintable       = "mintable"

	FlagOwner    = "owner"
	FlagMoniker  = "moniker"
	FlagIdentity = "identity"
	FlagDetails  = "details"
	FlagWebsite  = "website"
)

var (
	FsAssetIssue    = flag.NewFlagSet("", flag.ContinueOnError)
	FsGatewayCreate = flag.NewFlagSet("", flag.ContinueOnError)
	FsGatewayEdit   = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsAssetIssue.String(FlagFamily, "", "the asset family, valid values can be fungible and non-fungible")
	FsAssetIssue.String(FlagSource, "", "the asset source, valid values can be native, external and gateway")
	FsAssetIssue.String(FlagGateway, "", "the gateway name of gateway asset. For use in conjunction with --asset=gateway")
	FsAssetIssue.String(FlagSymbol, "", "the asset symbol. Once created, it cannot be modified")
	FsAssetIssue.String(FlagName, "", "the asset name")
	FsAssetIssue.Uint8(FlagDecimal, 0, "the asset decimal. The maximum value is 18")
	FsAssetIssue.String(FlagSymbolMinAlias, "", "the asset symbol minimum alias")
	FsAssetIssue.Uint64(FlagInitialSupply, 0, "the initial supply token of asset")
	FsAssetIssue.Uint64(FlagMaxSupply, asset.MaximumAssetMaxSupply, "the max supply token of asset")
	FsAssetIssue.Bool(FlagMintable, false, "whether the asset can be minted, default false")

	FsGatewayCreate.String(FlagMoniker, "", "the unique gateway name")
	FsGatewayCreate.String(FlagIdentity, "", "the gateway identity")
	FsGatewayCreate.String(FlagDetails, "", "the gateway description")
	FsGatewayCreate.String(FlagWebsite, "", "the external website")

	FsGatewayEdit.String(FlagMoniker, "", "the unique gateway name")
	FsGatewayEdit.String(FlagIdentity, "", "the gateway identity")
	FsGatewayEdit.String(FlagDetails, "", "the gateway description")
	FsGatewayEdit.String(FlagWebsite, "", "the external website")
}
