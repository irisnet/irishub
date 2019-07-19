package cli

import (
	"github.com/irisnet/irishub/app/v1/rand"
	flag "github.com/spf13/pflag"
)

const (
	FlagReqID         = "request-id"
	FlagConsumer      = "consumer"
	FlagBlockInterval = "block-interval"
	FlagQueueHeight   = "queue-height"
)

var (
	FsRequestRand = flag.NewFlagSet("", flag.ContinueOnError)
	FsQueryRand   = flag.NewFlagSet("", flag.ContinueOnError)
	FsQueryRands  = flag.NewFlagSet("", flag.ContinueOnError)
	FsQueryQueue  = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsRequestRand.Uint64(FlagBlockInterval, rand.DefaultBlockInterval, "the block interval")
	FsQueryRand.String(FlagReqID, "", "the request id")
	FsQueryRands.String(FlagConsumer, "", "optional consumer address")
	FsQueryQueue.Int64(FlagQueueHeight, 0, "optional height")
}
