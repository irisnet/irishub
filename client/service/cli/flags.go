// nolint
package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagName              = "name"
	FlagDescription       = "description"
	FlagTags              = "tags"
	FlagAuthorDescription = "author-description"
	FlagSchemas           = "schemas"
	FlagServiceName       = "service-name"
	FlagDeposit           = "deposit"
	FlagPricing           = "pricing"
	FlagWithdrawAddr      = "withdraw-addr"
	FlagMethodID          = "method-id"
	FlagServiceFee        = "service-fee"
	FlagReqData           = "request-data"
	FlagProfiling         = "profiling"
	FlagProvider          = "provider"
	FlagDefChainID        = "def-chain-id"
	FlagBindChainID       = "bind-chain-id"
	FlagReqChainID        = "request-chain-id"
	FlagReqID             = "request-id"
	FlagRespData          = "response-data"
	FlagErrMsg            = "error-msg"
	FlagDestAddress       = "dest-address"
	FlagWithdrawAmount    = "withdraw-amount"
)

// common flagsets to add to various functions
var (
	FsServiceDefine          = flag.NewFlagSet("", flag.ContinueOnError)
	FsServiceBind            = flag.NewFlagSet("", flag.ContinueOnError)
	FsServiceUpdateBinding   = flag.NewFlagSet("", flag.ContinueOnError)
	FsServiceSetWithdrawAddr = flag.NewFlagSet("", flag.ContinueOnError)
	FsServiceEnable          = flag.NewFlagSet("", flag.ContinueOnError)
	FsServiceRequest         = flag.NewFlagSet("", flag.ContinueOnError)
	FsServiceResponse        = flag.NewFlagSet("", flag.ContinueOnError)
	FsServiceWithdrawTax     = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsServiceDefine.String(FlagName, "", "service name")
	FsServiceDefine.String(FlagDescription, "", "service description")
	FsServiceDefine.StringSlice(FlagTags, []string{}, "service tags")
	FsServiceDefine.String(FlagAuthorDescription, "", "service author description")
	FsServiceDefine.String(FlagSchemas, "", "interface schemas content or path")

	FsServiceBind.String(FlagServiceName, "", "service name")
	FsServiceBind.String(FlagDeposit, "", "deposit of the binding")
	FsServiceBind.String(FlagPricing, "", "pricing content or path")
	FsServiceBind.String(FlagWithdrawAddr, "", "withdrawal address of the binding")

	FsServiceUpdateBinding.String(FlagDeposit, "", "added deposit for the binding")
	FsServiceUpdateBinding.String(FlagPricing, "", "new pricing of the binding, which is a Pricing JSON Schema instance")

	FsServiceSetWithdrawAddr.String(FlagWithdrawAddr, "", "withdrawal address of the binding")

	FsServiceEnable.String(FlagDeposit, "", "added deposit for enabling the binding")

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
