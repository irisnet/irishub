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
	FlagTotal         = "total"
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
	FsCreateFeed.String(FlagFeedName, "", "unique identifier of the Feed")
	FsCreateFeed.String(FlagAggregateFunc, "", "name of predefined function for processing the Service responses, e.g., avg, max, min etc.")
	FsCreateFeed.String(FlagValueJsonPath, "", "json path used by aggregate function to retrieve the value property from responses")
	FsCreateFeed.Uint64(FlagLatestHistory, 0, "number of latest history values to be saved for the Feed, range [1, MaxHistory]")
	FsCreateFeed.String(FlagDescription, "", "description of the Feed")
	FsCreateFeed.String(FlagServiceName, "", "name of the Service to be invoked by the Feed")
	FsCreateFeed.StringSlice(FlagProviders, []string{}, "list of service provider addresses")
	FsCreateFeed.String(FlagInput, "", "input argument (JSON) used to invoke the Service")
	FsCreateFeed.Int64(FlagTimeout, 0, "number of blocks to wait since a request is sent, beyond which responses will be ignored")
	FsCreateFeed.String(FlagServiceFeeCap, "", "only providers charging a fee lower than the cap will be invoked")
	FsCreateFeed.Uint64(FlagFrequency, 0, "frequency of sending repeated requests")
	FsCreateFeed.Int64(FlagTotal, 0, "total number of calls for repetitive requests,  -1 means unlimited")
	FsCreateFeed.Uint16(FlagThreshold, 0, "minimum number of responses needed for aggregation, range [1, count(providers)]")
	FsCreateFeed.String(FlagCreator, "", "address of the Feed creator")

	FsStartFeed.String(FlagFeedName, "", "unique identifier of the Feed")
	FsStartFeed.String(FlagCreator, "", "address of the Feed creator")

	FsPauseFeed.String(FlagFeedName, "", "unique identifier of the Feed")
	FsPauseFeed.String(FlagCreator, "", "address of the Feed creator")

	FsEditFeed.String(FlagFeedName, "", "unique identifier of the Feed")
	FsEditFeed.Uint64(FlagLatestHistory, 0, "number of latest history values to be saved for the Feed, range [1, MaxHistory]")
	FsEditFeed.StringSlice(FlagProviders, []string{}, "list of service provider addresses")
	FsEditFeed.String(FlagDescription, "", "description of the Feed")
	FsEditFeed.Int64(FlagTimeout, 0, "number of blocks to wait since a request is sent, beyond which responses will be ignored")
	FsEditFeed.String(FlagServiceFeeCap, "", "only providers charging a fee lower than the cap will be invoked")
	FsEditFeed.Uint64(FlagFrequency, 0, "frequency of sending repeated requests")
	FsEditFeed.Int64(FlagTotal, 0, "total number of calls for repetitive requests,  -1 means unlimited")
	FsEditFeed.Uint16(FlagThreshold, 0, "minimum number of responses needed for aggregation, range [1, count(providers)]")
	FsEditFeed.String(FlagCreator, "", "address of the Feed creator")

	FsQueryFeeds.String(FlagFeedState, "", "the state of the Feed,paused|running")
}
