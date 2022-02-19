package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagTokenName   = "name"
	FlagURI         = "uri"
	FlagURIHash     = "uri-hash"
	FlagDescription = "description"
	FlagRecipient   = "recipient"
	FlagOwner       = "owner"
	FlagData        = "data"

	FlagDenomName        = "name"
	FlagDenomID          = "denom-id"
	FlagSchema           = "schema"
	FlagSymbol           = "symbol"
	FlagMintRestricted   = "mint-restricted"
	FlagUpdateRestricted = "update-restricted"
)

var (
	FsIssueDenom    = flag.NewFlagSet("", flag.ContinueOnError)
	FsMintMT       = flag.NewFlagSet("", flag.ContinueOnError)
	FsEditMT       = flag.NewFlagSet("", flag.ContinueOnError)
	FsTransferMT   = flag.NewFlagSet("", flag.ContinueOnError)
	FsQuerySupply   = flag.NewFlagSet("", flag.ContinueOnError)
	FsQueryOwner    = flag.NewFlagSet("", flag.ContinueOnError)
	FsTransferDenom = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsIssueDenom.String(FlagSchema, "", "Denom data structure definition")
	FsIssueDenom.String(FlagURI, "", "The uri for the class metadata stored off chain. It can define schema for Class and MT `Data` attributes. Optional")
	FsIssueDenom.String(FlagURIHash, "", "The uri-hash is a hash of the document pointed by uri. Optional")
	FsIssueDenom.String(FlagDescription, "", "The description is a brief description of mt classification. Optional")
	FsIssueDenom.String(FlagDenomName, "", "The name of the denom")
	FsIssueDenom.String(FlagSymbol, "", "The symbol of the denom")
	FsIssueDenom.String(FlagData, "", "The data is the app specific metadata of the MT class. Optional")
	FsIssueDenom.Bool(FlagMintRestricted, false, "mint restricted of mt under denom")
	FsIssueDenom.Bool(FlagUpdateRestricted, false, "update restricted of mt under denom")

	FsMintMT.String(FlagURI, "", "The uri for supplemental off-chain tokenData (should return a JSON object)")
	FsMintMT.String(FlagURIHash, "", "The uri_hash is a hash of the document pointed by uri. Optional")
	FsMintMT.String(FlagRecipient, "", "The receiver of the mt, if not filled, the default is the sender of the transaction")
	FsMintMT.String(FlagData, "", "The origin data of the mt")
	FsMintMT.String(FlagTokenName, "", "The name of the mt")

	FsEditMT.String(FlagURI, "[do-not-modify]", "URI for the supplemental off-chain token data (should return a JSON object)")
	FsEditMT.String(FlagURIHash, "[do-not-modify]", "The uri_hash is a hash of the document pointed by uri. Optional")
	FsEditMT.String(FlagData, "[do-not-modify]", "The token data of the mt")
	FsEditMT.String(FlagTokenName, "[do-not-modify]", "The name of the mt")

	FsTransferMT.String(FlagURI, "[do-not-modify]", "URI for the supplemental off-chain token data (should return a JSON object)")
	FsTransferMT.String(FlagURIHash, "[do-not-modify]", "The uri_hash is a hash of the document pointed by uri. Optional")
	FsTransferMT.String(FlagData, "[do-not-modify]", "The token data of the mt")
	FsTransferMT.String(FlagTokenName, "[do-not-modify]", "The name of the mt")

	FsQuerySupply.String(FlagOwner, "", "The owner of the mt")

	FsQueryOwner.String(FlagDenomID, "", "The name of the collection")
}
