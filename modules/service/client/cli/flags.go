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
	FlagOwner             = "owner"
	FlagProvider          = "provider"
	FlagDeposit           = "deposit"
	FlagPricing           = "pricing"
	FlagQoS               = "qos"
	FlagOptions           = "options"
	FlagProviders         = "providers"
	FlagServiceFeeCap     = "service-fee-cap"
	FlagTimeout           = "timeout"
	FlagData              = "data"
	FlagRepeated          = "repeated"
	FlagFrequency         = "frequency"
	FlagTotal             = "total"
	FlagRequestID         = "request-id"
	FlagResult            = "result"
)

// common flagsets to add to various functions
var (
	FsDefineService        = flag.NewFlagSet("", flag.ContinueOnError)
	FsBindService          = flag.NewFlagSet("", flag.ContinueOnError)
	FsUpdateServiceBinding = flag.NewFlagSet("", flag.ContinueOnError)
	FsEnableServiceBinding = flag.NewFlagSet("", flag.ContinueOnError)
	FsCallService          = flag.NewFlagSet("", flag.ContinueOnError)
	FsRespondService       = flag.NewFlagSet("", flag.ContinueOnError)
	FsUpdateRequestContext = flag.NewFlagSet("", flag.ContinueOnError)
	FsQueryServiceBindings = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsDefineService.String(FlagName, "", "Service name")
	FsDefineService.String(FlagDescription, "", "Service description")
	FsDefineService.StringSlice(FlagTags, []string{}, "Service tags")
	FsDefineService.String(FlagAuthorDescription, "", "Service author description")
	FsDefineService.String(FlagSchemas, "", "Interface schemas content or file path")

	FsBindService.String(FlagServiceName, "", "Service name")
	FsBindService.String(FlagProvider, "", "Provider address, default to the owner")
	FsBindService.String(FlagDeposit, "", "Deposit of the binding")
	FsBindService.String(FlagPricing, "", "Pricing content or file path, which is an instance of the Service Pricing schema")
	FsBindService.Uint64(FlagQoS, 0, "Quality of service, in terms of minimum response time")
	FsBindService.String(FlagOptions, "", "Non-functional requirements options")

	FsUpdateServiceBinding.String(FlagDeposit, "", "Added deposit for the binding")
	FsUpdateServiceBinding.String(FlagPricing, "", "Pricing content or file path, which is an instance of the Service Pricing schema")
	FsUpdateServiceBinding.Uint64(FlagQoS, 0, "Quality of service, in terms of minimum response time, not updated if set to 0")
	FsUpdateServiceBinding.String(FlagOptions, "", "Non-functional requirements options")

	FsEnableServiceBinding.String(FlagDeposit, "", "Added deposit for enabling the binding")

	FsCallService.String(FlagServiceName, "", "Service name")
	FsCallService.StringSlice(FlagProviders, []string{}, "Provider list to request")
	FsCallService.String(FlagServiceFeeCap, "", "Maximum service fee to pay for a single request")
	FsCallService.String(FlagData, "", "Content or file path of the request input, which is an Input JSON schema instance")
	FsCallService.Int64(FlagTimeout, 0, "Request timeout")
	FsCallService.Bool(FlagRepeated, false, "Indicate if the request is repetitive")
	FsCallService.Uint64(FlagFrequency, 0, "Request frequency when repeated, default to timeout")
	FsCallService.Int64(FlagTotal, 0, "Request count when repeated, -1 means unlimited")

	FsRespondService.String(FlagRequestID, "", "ID of the request to respond to")
	FsRespondService.String(FlagResult, "", "Content or file path of the response result, which is an Result JSON schema instance")
	FsRespondService.String(FlagData, "", "Content or file path of the response output, which is an Output JSON schema instance")

	FsUpdateRequestContext.StringSlice(FlagProviders, []string{}, "Provider list to request, not updated if empty")
	FsUpdateRequestContext.String(FlagServiceFeeCap, "", "Maximum service fee to pay for a single request, not updated if empty")
	FsUpdateRequestContext.Int64(FlagTimeout, 0, "Request timeout, not updated if set to 0")
	FsUpdateRequestContext.Uint64(FlagFrequency, 0, "Request frequency, not updated if set to 0")
	FsUpdateRequestContext.Int64(FlagTotal, 0, "Request count, not updated if set to 0")

	FsQueryServiceBindings.String(FlagOwner, "", "The owner of bindings, which is optional")
}
