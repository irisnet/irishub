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
	FlagFile               = "file"
	FlagProvider           = "provider"
	FlagBindChainID        = "bind-chain-id"
	FlagBindType           = "bind-type"
	FlagDeposit            = "deposit"
	FlagPrices             = "prices"
	FlagAvgRspTime         = "avg-rsp-time"
	FlagUsableTime         = "usable-time"
	FlagExpiration         = "expiration"
	FlagMethodID           = "method-id"
	FlagServiceFee         = "service-fee"
	FlagReqData            = "request-data"
	FlagRespData           = "response-data"
	FlagErrMsg             = "error-msg"
	FlagProfiling          = "profiling"
	FlagReqChainId         = "request-chain-id"
	FlagReqId              = "request-id"
)

var (
	FsDefChainID         = flag.NewFlagSet("", flag.ContinueOnError)
	FsServiceName        = flag.NewFlagSet("", flag.ContinueOnError)
	FsServiceDescription = flag.NewFlagSet("", flag.ContinueOnError)
	FsTags               = flag.NewFlagSet("", flag.ContinueOnError)
	FsAuthorDescription  = flag.NewFlagSet("", flag.ContinueOnError)
	FsIdlContent         = flag.NewFlagSet("", flag.ContinueOnError)
	FsFile               = flag.NewFlagSet("", flag.ContinueOnError)
	FsProvider           = flag.NewFlagSet("", flag.ContinueOnError)
	FsBindChainID        = flag.NewFlagSet("", flag.ContinueOnError)
	FsBindType           = flag.NewFlagSet("", flag.ContinueOnError)
	FsDeposit            = flag.NewFlagSet("", flag.ContinueOnError)
	FsPrices             = flag.NewFlagSet("", flag.ContinueOnError)
	FsAvgRspTime         = flag.NewFlagSet("", flag.ContinueOnError)
	FsUsableTime         = flag.NewFlagSet("", flag.ContinueOnError)
	FsExpiration         = flag.NewFlagSet("", flag.ContinueOnError)
	FsMethodID           = flag.NewFlagSet("", flag.ContinueOnError)
	FsServiceFee         = flag.NewFlagSet("", flag.ContinueOnError)
	FsReqData            = flag.NewFlagSet("", flag.ContinueOnError)
	FsRespData           = flag.NewFlagSet("", flag.ContinueOnError)
	FsErrMsg             = flag.NewFlagSet("", flag.ContinueOnError)
	FsProfiling          = flag.NewFlagSet("", flag.ContinueOnError)
	FsReqChainId         = flag.NewFlagSet("", flag.ContinueOnError)
	FsReqId              = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsDefChainID.String(FlagDefChainID, "", "the ID of the blockchain defined of the service")
	FsServiceName.String(FlagServiceName, "", "service name")
	FsServiceDescription.String(FlagServiceDescription, "", "service description")
	FsTags.StringSlice(FlagTags, []string{}, "service tags")
	FsAuthorDescription.String(FlagAuthorDescription, "", "service author description")
	FsIdlContent.String(FlagIdlContent, "", "content of service interface description language")
	FsFile.String(FlagFile, "", "path of file which contains service interface description language")

	FsProvider.String(FlagProvider, "", "bech32 encoded account created the service binding")
	FsBindChainID.String(FlagBindChainID, "", "the ID of the blockchain bond of the service")
	FsBindType.String(FlagBindType, "", "type of binding, valid values can be Local and Global")
	FsDeposit.String(FlagDeposit, "", "deposit of binding")
	FsPrices.StringSlice(FlagPrices, []string{}, "prices of binding, will contains all method")
	FsAvgRspTime.Int64(FlagAvgRspTime, 0, "the average service response time in milliseconds")
	FsUsableTime.Int64(FlagUsableTime, 0, "an integer represents the number of usable service invocations per 10,000")
	FsExpiration.String(FlagExpiration, "", "the blockchain height where this binding expires")

	FsMethodID.Int16(FlagMethodID, 0, "the method id called")
	FsServiceFee.String(FlagServiceFee, "", "fee to pay for a service invocation")
	FsReqData.BytesHex(FlagReqData, nil, "hex encoded request data of a service invocation")
	FsRespData.BytesHex(FlagRespData, nil, "hex encoded response data of a service invocation")
	FsErrMsg.BytesHex(FlagErrMsg, nil, "hex encoded response error msg of a service invocation")
	FsProfiling.Bool(FlagProfiling, false, "service invocation profiling model, default false")
	FsReqChainId.String(FlagReqChainId, "", "the ID of the blockchain that the service invocation initiated")
	FsReqId.String(FlagReqId, "", "the ID of the service invocation")
}
