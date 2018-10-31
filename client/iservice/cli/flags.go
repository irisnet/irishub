package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagDefChainID         = "def-chain-id"
	FlagServiceName        = "service-name"
	FlagServiceDescription = "service-description"
	FlagTags               = "tags"
	FlagAuthorDescription  = "author-description"
	FlagIdlContent         = "idl-content"
	FlagMessaging          = "messaging"
	FlagFile               = "file"
	FlagProvider           = "provider"
	FlagBindChainID        = "bind-chain-id"
	FlagBindType           = "bind-type"
	FlagDeposit            = "deposit"
	FlagPrices             = "price"
	FlagLevels             = "levels"
	FlagExpiration         = "expiration"
)

var (
	FsDefChainID         = flag.NewFlagSet("", flag.ContinueOnError)
	FsServiceName        = flag.NewFlagSet("", flag.ContinueOnError)
	FsServiceDescription = flag.NewFlagSet("", flag.ContinueOnError)
	FsTags               = flag.NewFlagSet("", flag.ContinueOnError)
	FsAuthorDescription  = flag.NewFlagSet("", flag.ContinueOnError)
	FsIdlContent         = flag.NewFlagSet("", flag.ContinueOnError)
	FsMessaging          = flag.NewFlagSet("", flag.ContinueOnError)
	FsFile               = flag.NewFlagSet("", flag.ContinueOnError)
	FsProvider           = flag.NewFlagSet("", flag.ContinueOnError)
	FsBindChainID        = flag.NewFlagSet("", flag.ContinueOnError)
	FsBindType           = flag.NewFlagSet("", flag.ContinueOnError)
	FsDeposit            = flag.NewFlagSet("", flag.ContinueOnError)
	FsPrices             = flag.NewFlagSet("", flag.ContinueOnError)
	FsLevels             = flag.NewFlagSet("", flag.ContinueOnError)
	FsExpiration         = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsDefChainID.String(FlagDefChainID, "", "the ID of the blockchain defined of the iService")
	FsServiceName.String(FlagServiceName, "", "service name")
	FsServiceDescription.String(FlagServiceDescription, "", "service description")
	FsTags.String(FlagTags, "", "service tags")
	FsAuthorDescription.String(FlagAuthorDescription, "", "service author description")
	FsIdlContent.String(FlagIdlContent, "", "content of service interface description language")
	FsMessaging.String(FlagMessaging, "", "service messaging type, valid values can be Unicast and Multicast")
	FsFile.String(FlagFile, "", "path of file which contains service interface description language")

	FsProvider.String(FlagProvider, "", "bech32 encoded account created the iService binding")
	FsBindChainID.String(FlagBindChainID, "", "the ID of the blockchain bond of the iService")
	FsBindType.String(FlagBindType, "", "	")
	FsDeposit.String(FlagDeposit, "", "path of file which contains service interface description language")
	FsPrices.String(FlagPrices, "", "path of file which contains service interface description language")
	FsLevels.String(FlagLevels, "", "path of file which contains service interface description language")
	FsExpiration.String(FlagExpiration, "", "path of file which contains service interface description language")
}
