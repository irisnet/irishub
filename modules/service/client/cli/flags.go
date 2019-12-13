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
	FlagMethodID           = "method-id"
	FlagServiceFee         = "service-fee"
	FlagReqData            = "request-data"
	FlagRespData           = "response-data"
	FlagErrMsg             = "error-msg"
	FlagProfiling          = "profiling"
	FlagReqChainID         = "request-chain-id"
	FlagReqID              = "request-id"
	FlagDestAddress        = "dest-address"
	FlagWithdrawAmount     = "withdraw-amount"
)

var (
	FsServiceDefinitionCreate = flag.NewFlagSet("", flag.ContinueOnError)
	FsServiceBindingCreate    = flag.NewFlagSet("", flag.ContinueOnError)
	FsServiceBindingUpdate    = flag.NewFlagSet("", flag.ContinueOnError)
	FsServiceDefinition       = flag.NewFlagSet("", flag.ContinueOnError)
	FsServiceBinding          = flag.NewFlagSet("", flag.ContinueOnError)
	FsServiceRequest          = flag.NewFlagSet("", flag.ContinueOnError)
	FsServiceResponse         = flag.NewFlagSet("", flag.ContinueOnError)
	FsServiceWithdrawTax      = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsServiceDefinitionCreate.String(FlagServiceName, "", "service name")
	FsServiceDefinitionCreate.String(FlagServiceDescription, "", "service description")
	FsServiceDefinitionCreate.StringSlice(FlagTags, []string{}, "service tags")
	FsServiceDefinitionCreate.String(FlagAuthorDescription, "", "service author description")
	FsServiceDefinitionCreate.String(FlagIdlContent, "", "content of service interface description language")
	FsServiceDefinitionCreate.String(FlagFile, "", "path of file which contains service interface description language")

	FsServiceDefinition.String(FlagDefChainID, "", "the ID of the blockchain defined of the service")
	FsServiceDefinition.String(FlagServiceName, "", "service name")

	FsServiceBindingCreate.String(FlagBindType, "", "type of binding, valid values can be Local and Global")
	FsServiceBindingCreate.String(FlagDeposit, "", "deposit of binding")
	FsServiceBindingCreate.StringSlice(FlagPrices, []string{}, "prices of binding, will contains all method")
	FsServiceBindingCreate.Int64(FlagAvgRspTime, 0, "the average service response time in milliseconds")
	FsServiceBindingCreate.Int64(FlagUsableTime, 0, "an integer represents the number of usable service invocations per 10,000")

	FsServiceBindingUpdate.String(FlagBindType, "", "type of binding, valid values can be Local and Global")
	FsServiceBindingUpdate.String(FlagDeposit, "", "deposit of binding")
	FsServiceBindingUpdate.StringSlice(FlagPrices, []string{}, "prices of binding, will contains all method")
	FsServiceBindingUpdate.Int64(FlagAvgRspTime, 0, "the average service response time in milliseconds")
	FsServiceBindingUpdate.Int64(FlagUsableTime, 0, "an integer represents the number of usable service invocations per 10,000")

	FsServiceBinding.String(FlagBindChainID, "", "the ID of the blockchain bond of the service")
	FsServiceBinding.String(FlagProvider, "", "bech32 encoded account created the service binding")

	FsServiceRequest.Int16(FlagMethodID, 0, "the method id called")
	FsServiceRequest.String(FlagServiceFee, "", "fee to pay for a service invocation")
	FsServiceRequest.BytesHex(FlagReqData, nil, "hex encoded request data of a service invocation")
	FsServiceRequest.Bool(FlagProfiling, false, "service invocation profiling model, default false")

	FsServiceResponse.BytesHex(FlagRespData, nil, "hex encoded response data of a service invocation")
	FsServiceResponse.BytesHex(FlagErrMsg, nil, "hex encoded response error msg of a service invocation")
	FsServiceResponse.String(FlagReqChainID, "", "the ID of the blockchain that the service invocation initiated")
	FsServiceResponse.String(FlagReqID, "", "the ID of the service invocation")

	FsServiceWithdrawTax.String(FlagDestAddress, "", "bech32 encoded address of the destination account")
	FsServiceWithdrawTax.String(FlagWithdrawAmount, "", "withdraw amount")
}
