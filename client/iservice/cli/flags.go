package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagDefChainID         = "def-chain-id"
	FlagServiceName        = "name"
	FlagService1Name       = "service-name"
	FlagServiceDescription = "service-description"
	FlagTags               = "tags"
	FlagAuthorDescription  = "author-description"
	FlagIdlContent         = "idl-content"
	FlagBroadcast          = "broadcast"
	FlagMessaging          = "messaging"
	FlagFile               = "file"
	FlagProvider           = "provider"
	FlagBindChainID        = "bind-chain-id"
	FlagBindType           = "bind-type"
	FlagDeposit            = "deposit"
	FlagPrices             = "prices"
	FlagAvgRspTime         = "avg-rsp-time"
	FlagUsableTime         = "usable-time"
	FlagExpiration         = "expiration"
)

var (
	FsDefChainID         = flag.NewFlagSet("", flag.ContinueOnError)
	FsServiceName        = flag.NewFlagSet("", flag.ContinueOnError)
	FsService1Name       = flag.NewFlagSet("", flag.ContinueOnError)
	FsServiceDescription = flag.NewFlagSet("", flag.ContinueOnError)
	FsTags               = flag.NewFlagSet("", flag.ContinueOnError)
	FsAuthorDescription  = flag.NewFlagSet("", flag.ContinueOnError)
	FsIdlContent         = flag.NewFlagSet("", flag.ContinueOnError)
	FsBroadcast          = flag.NewFlagSet("", flag.ContinueOnError)
	FsMessaging          = flag.NewFlagSet("", flag.ContinueOnError)
	FsFile               = flag.NewFlagSet("", flag.ContinueOnError)
	FsProvider           = flag.NewFlagSet("", flag.ContinueOnError)
	FsBindChainID        = flag.NewFlagSet("", flag.ContinueOnError)
	FsBindType           = flag.NewFlagSet("", flag.ContinueOnError)
	FsDeposit            = flag.NewFlagSet("", flag.ContinueOnError)
	FsPrices             = flag.NewFlagSet("", flag.ContinueOnError)
	FsAvgRspTime         = flag.NewFlagSet("", flag.ContinueOnError)
	FsUsableTime         = flag.NewFlagSet("", flag.ContinueOnError)
	FsExpiration         = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsDefChainID.String(FlagDefChainID, "", "the ID of the blockchain defined of the iService")
	FsServiceName.String(FlagServiceName, "", "service name")
	FsService1Name.String(FlagService1Name, "", "service name")
	FsServiceDescription.String(FlagServiceDescription, "", "service description")
	FsTags.StringSlice(FlagTags, []string{}, "service tags")
	FsAuthorDescription.String(FlagAuthorDescription, "", "service author description")
	FsIdlContent.String(FlagIdlContent, "", "content of service interface description language")
	FsMessaging.String(FlagMessaging, "", "service messaging type, valid values can be Unicast and Multicast")
	FsBroadcast.String(FlagBroadcast, "", "service broadcast type, valid values can be Broadcast and Unicast")
	FsFile.String(FlagFile, "", "path of file which contains service interface description language")

	FsProvider.String(FlagProvider, "", "bech32 encoded account created the iService binding")
	FsBindChainID.String(FlagBindChainID, "", "the ID of the blockchain bond of the iService")
	FsBindType.String(FlagBindType, "", "type of binding, valid values can be Local and Global")
	FsDeposit.String(FlagDeposit, "", "deposit of binding")
	FsPrices.StringSlice(FlagPrices, []string{}, "prices of binding, will contains all method")
	FsAvgRspTime.String(FlagAvgRspTime, "", "the average service response time in milliseconds")
	FsUsableTime.String(FlagUsableTime, "", "an integer represents the number of usable service invocations per 10,000")
	FsExpiration.String(FlagExpiration, "", "the blockchain height where this binding expires")
}
