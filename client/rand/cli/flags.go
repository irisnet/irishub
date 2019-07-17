package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagReqID    = "request-id"
	FlagConsumer = "consumer"
	FlagHeight   = "height"
)

var (
	FsQueryRand     = flag.NewFlagSet("", flag.ContinueOnError)
	FsQueryRands    = flag.NewFlagSet("", flag.ContinueOnError)
	FsQueryRequest  = flag.NewFlagSet("", flag.ContinueOnError)
	FsQueryRequests = flag.NewFlagSet("", flag.ContinueOnError)
	FsQueryQueue    = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsQueryRand.String(FlagReqID, "", "the request id")
	FsQueryRands.String(FlagConsumer, "", "optional consumer address")
	FsQueryRequest.String(FlagReqID, "", "the request id")
	FsQueryRequests.String(FlagConsumer, "", "optional consumer address")
	FsQueryQueue.Int64(FlagHeight, 0, "optional height")
}
