package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagTokenName = "name"
	FlagTokenURI  = "uri"
	FlagTokenData = "data"
	FlagRecipient = "recipient"
	FlagOwner     = "owner"

	FlagDenomName = "name"
	FlagDenom     = "denom"
	FlagSchema    = "schema"
)

var (
	FsIssueDenom  = flag.NewFlagSet("", flag.ContinueOnError)
	FsMintNFT     = flag.NewFlagSet("", flag.ContinueOnError)
	FsEditNFT     = flag.NewFlagSet("", flag.ContinueOnError)
	FsTransferNFT = flag.NewFlagSet("", flag.ContinueOnError)
	FsQuerySupply = flag.NewFlagSet("", flag.ContinueOnError)
	FsQueryOwner  = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsIssueDenom.String(FlagSchema, "", "Denom data structure definition")
	FsIssueDenom.String(FlagDenomName, "", "The name of the denom")

	FsMintNFT.String(FlagTokenURI, "", "URI for supplemental off-chain tokenData (should return a JSON object)")
	FsMintNFT.String(FlagRecipient, "", "Receiver of the nft, if not filled, the default is the sender of the transaction")
	FsMintNFT.String(FlagTokenData, "", "The origin data of nft")
	FsMintNFT.String(FlagTokenName, "", "The name of nft")

	FsEditNFT.String(FlagTokenURI, "[do-not-modify]", "URI for supplemental off-chain tokenData (should return a JSON object)")
	FsEditNFT.String(FlagTokenData, "[do-not-modify]", "The tokenData of nft")
	FsEditNFT.String(FlagTokenName, "[do-not-modify]", "The name of nft")

	FsTransferNFT.String(FlagTokenURI, "[do-not-modify]", "URI for supplemental off-chain tokenData (should return a JSON object)")
	FsTransferNFT.String(FlagTokenData, "[do-not-modify]", "The tokenData of nft")
	FsTransferNFT.String(FlagTokenName, "[do-not-modify]", "The name of nft")

	FsQuerySupply.String(FlagOwner, "", "The owner of a nft")

	FsQueryOwner.String(FlagDenom, "", "The name of a collection")
}
