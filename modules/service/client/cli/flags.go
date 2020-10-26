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
	FlagSuperMode         = "super-mode"
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
	FsDefineService.String(FlagName, "", "service name")
	FsDefineService.String(FlagDescription, "", "service description")
	FsDefineService.StringSlice(FlagTags, []string{}, "service tags")
	FsDefineService.String(FlagAuthorDescription, "", "service author description")
	FsDefineService.String(FlagSchemas, "", "interface schemas content or file path")

	FsBindService.String(FlagServiceName, "", "service name")
	FsBindService.String(FlagProvider, "", "provider address, default to the owner")
	FsBindService.String(FlagDeposit, "", "deposit of the binding")
	FsBindService.String(FlagPricing, "", "pricing content or file path, which is an instance of the Service Pricing schema")
	FsBindService.Uint64(FlagQoS, 0, "quality of service, in terms of minimum response time")
	FsBindService.String(FlagOptions, "", "non-functional requirements options")

	FsUpdateServiceBinding.String(FlagDeposit, "", "added deposit for the binding")
	FsUpdateServiceBinding.String(FlagPricing, "", "pricing content or file path, which is an instance of the Service Pricing schema")
	FsUpdateServiceBinding.Uint64(FlagQoS, 0, "quality of service, in terms of minimum response time, not updated if set to 0")
	FsUpdateServiceBinding.String(FlagOptions, "", "non-functional requirements options")

	FsEnableServiceBinding.String(FlagDeposit, "", "added deposit for enabling the binding")

	FsCallService.String(FlagServiceName, "", "service name")
	FsCallService.StringSlice(FlagProviders, []string{}, "provider list to request")
	FsCallService.String(FlagServiceFeeCap, "", "maximum service fee to pay for a single request")
	FsCallService.String(FlagData, "", "content or file path of the request input, which is an Input JSON schema instance")
	FsCallService.Int64(FlagTimeout, 0, "request timeout")
	FsCallService.Bool(FlagSuperMode, false, "indicate if the signer is a super user")
	FsCallService.Bool(FlagRepeated, false, "indicate if the request is repetitive")
	FsCallService.Uint64(FlagFrequency, 0, "request frequency when repeated, default to timeout")
	FsCallService.Int64(FlagTotal, 0, "request count when repeated, -1 means unlimited")

	FsRespondService.String(FlagRequestID, "", "ID of the request to respond to")
	FsRespondService.String(FlagResult, "", "content or file path of the response result, which is an Result JSON schema instance")
	FsRespondService.String(FlagData, "", "content or file path of the response output, which is an Output JSON schema instance")

	FsUpdateRequestContext.StringSlice(FlagProviders, []string{}, "provider list to request, not updated if empty")
	FsUpdateRequestContext.String(FlagServiceFeeCap, "", "maximum service fee to pay for a single request, not updated if empty")
	FsUpdateRequestContext.Uint64(FlagTimeout, 0, "request timeout, not updated if set to 0")
	FsUpdateRequestContext.Uint64(FlagFrequency, 0, "request frequency, not updated if set to 0")
	FsUpdateRequestContext.Int64(FlagTotal, 0, "request count, not updated if set to 0")

	FsQueryServiceBindings.String(FlagOwner, "", "the owner of bindings, which is optional")
}
