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
	FsMintNFT       = flag.NewFlagSet("", flag.ContinueOnError)
	FsEditNFT       = flag.NewFlagSet("", flag.ContinueOnError)
	FsTransferNFT   = flag.NewFlagSet("", flag.ContinueOnError)
	FsQuerySupply   = flag.NewFlagSet("", flag.ContinueOnError)
	FsQueryOwner    = flag.NewFlagSet("", flag.ContinueOnError)
	FsTransferDenom = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsIssueDenom.String(FlagSchema, "", "Denom data structure definition")
	FsIssueDenom.String(FlagURI, "", "The uri for the class metadata stored off chain. It can define schema for Class and NFT `Data` attributes. Optional")
	FsIssueDenom.String(FlagURIHash, "", "The uri-hash is a hash of the document pointed by uri. Optional")
	FsIssueDenom.String(FlagDescription, "", "The description is a brief description of nft classification. Optional")
	FsIssueDenom.String(FlagDenomName, "", "The name of the denom")
	FsIssueDenom.String(FlagSymbol, "", "The symbol of the denom")
	FsIssueDenom.String(FlagData, "", "The data is the app specific metadata of the NFT class. Optional")
	FsIssueDenom.Bool(FlagMintRestricted, false, "mint restricted of nft under denom")
	FsIssueDenom.Bool(FlagUpdateRestricted, false, "update restricted of nft under denom")

	FsMintNFT.String(FlagURI, "", "The uri for supplemental off-chain tokenData (should return a JSON object)")
	FsMintNFT.String(FlagURIHash, "", "The uri_hash is a hash of the document pointed by uri. Optional")
	FsMintNFT.String(FlagRecipient, "", "The receiver of the nft, if not filled, the default is the sender of the transaction")
	FsMintNFT.String(FlagData, "", "The origin data of the nft")
	FsMintNFT.String(FlagTokenName, "", "The name of the nft")

	FsEditNFT.String(FlagURI, "[do-not-modify]", "URI for the supplemental off-chain token data (should return a JSON object)")
	FsEditNFT.String(FlagURIHash, "[do-not-modify]", "The uri_hash is a hash of the document pointed by uri. Optional")
	FsEditNFT.String(FlagData, "[do-not-modify]", "The token data of the nft")
	FsEditNFT.String(FlagTokenName, "[do-not-modify]", "The name of the nft")

	FsTransferNFT.String(FlagURI, "[do-not-modify]", "URI for the supplemental off-chain token data (should return a JSON object)")
	FsTransferNFT.String(FlagURIHash, "[do-not-modify]", "The uri_hash is a hash of the document pointed by uri. Optional")
	FsTransferNFT.String(FlagData, "[do-not-modify]", "The token data of the nft")
	FsTransferNFT.String(FlagTokenName, "[do-not-modify]", "The name of the nft")

	FsQuerySupply.String(FlagOwner, "", "The owner of the nft")

	FsQueryOwner.String(FlagDenomID, "", "The name of the collection")
}
