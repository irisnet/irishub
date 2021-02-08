package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagFeedName      = "feed-name"
	FlagAggregateFunc = "aggregate-func"
	FlagValueJsonPath = "value-json-path"
	FlagLatestHistory = "latest-history"
	FlagDescription   = "description"
	FlagServiceName   = "service-name"
	FlagProviders     = "providers"
	FlagInput         = "input"
	FlagTimeout       = "timeout"
	FlagServiceFeeCap = "service-fee-cap"
	FlagFrequency     = "frequency"
	FlagThreshold     = "threshold"
	FlagCreator       = "creator"
	FlagFeedState     = "state"
)

var (
	FsCreateFeed     = flag.NewFlagSet("", flag.ContinueOnError)
	FsStartFeed      = flag.NewFlagSet("", flag.ContinueOnError)
	FsPauseFeed      = flag.NewFlagSet("", flag.ContinueOnError)
	FsEditFeed       = flag.NewFlagSet("", flag.ContinueOnError)
	FsQueryFeed      = flag.NewFlagSet("", flag.ContinueOnError)
	FsQueryFeeds     = flag.NewFlagSet("", flag.ContinueOnError)
	FsQueryFeedValue = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsCreateFeed.String(FlagFeedName, "", "The unique identifier of the feed")
	FsCreateFeed.String(FlagAggregateFunc, "", "The name of predefined function for processing the service responses, e.g., avg, max, min, etc.")
	FsCreateFeed.String(FlagValueJsonPath, "", "The field name or path of service response result used to retrieve the value property of aggregate-func from response results")
	FsCreateFeed.Uint64(FlagLatestHistory, 0, "The maximum number of the latest history values to be saved for the feed, range [1, 100]")
	FsCreateFeed.String(FlagDescription, "", "The description of the feed.")
	FsCreateFeed.String(FlagServiceName, "", "The name of the service to be invoked by the feed")
	FsCreateFeed.StringSlice(FlagProviders, []string{}, "The list of service provider addresses")
	FsCreateFeed.String(FlagInput, "", "The input argument (JSON) used to invoke the service")
	FsCreateFeed.Int64(FlagTimeout, 0, "The maximum number of blocks to wait for a response since a request is sent, beyond which the request will be ignored")
	FsCreateFeed.String(FlagServiceFeeCap, "", "Only providers charging a fee lower than the cap will be invoked")
	FsCreateFeed.Uint64(FlagFrequency, 0, "The invocation frequency of sending repeated requests")
	FsCreateFeed.Uint32(FlagThreshold, 0, "The minimum number of responses needed for aggregation, range [1, Length(providers)]")
	FsCreateFeed.String(FlagCreator, "", "Address of the feed creator")

	FsStartFeed.String(FlagFeedName, "", "The unique identifier of the feed")
	FsStartFeed.String(FlagCreator, "", "Address of the feed creator")

	FsPauseFeed.String(FlagFeedName, "", "The unique identifier of the feed")
	FsPauseFeed.String(FlagCreator, "", "Address of the feed creator")

	FsEditFeed.Uint64(FlagLatestHistory, 0, "The maximum number of the latest history values to be saved for the feed, range [1, 100]")
	FsEditFeed.StringSlice(FlagProviders, []string{}, "The list of service provider addresses")
	FsEditFeed.String(FlagDescription, "", "The description of the feed.")
	FsEditFeed.Int64(FlagTimeout, 0, "The maximum number of blocks to wait for a response since a request is sent, beyond which the request will be ignored")
	FsEditFeed.String(FlagServiceFeeCap, "", "Only providers charging a fee lower than the cap will be invoked")
	FsEditFeed.Uint64(FlagFrequency, 0, "The invocation frequency of sending repeated requests")
	FsEditFeed.Uint32(FlagThreshold, 0, "The minimum number of responses needed for aggregation, range [1, Length(providers)]")
	FsEditFeed.String(FlagCreator, "", "Address of the feed creator")

	FsQueryFeeds.String(FlagFeedState, "", "The state of the feed, paused|running")
}
