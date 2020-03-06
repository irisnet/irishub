package cli

import (
	flag "github.com/spf13/pflag"

	"github.com/irisnet/irishub/app/v3/rand"
)

const (
	FlagReqID         = "request-id"
	FlagConsumer      = "consumer"
	FlagBlockInterval = "block-interval"
	FlagQueueHeight   = "queue-height"
	FlagOracle        = "oracle"
)

var (
	FsRequestRand = flag.NewFlagSet("", flag.ContinueOnError)
	FsQueryRand   = flag.NewFlagSet("", flag.ContinueOnError)
	FsQueryQueue  = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsRequestRand.Uint64(FlagBlockInterval, rand.DefaultBlockInterval, "the block interval")
	FsQueryRand.String(FlagReqID, "", "the request id")
	FsQueryQueue.Int64(FlagQueueHeight, 0, "optional height")
}
