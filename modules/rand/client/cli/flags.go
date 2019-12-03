package cli

import (
	"github.com/irisnet/irishub/modules/rand"
	flag "github.com/spf13/pflag"
)

const (
	FlagReqID         = "request-id"
	FlagConsumer      = "consumer"
	FlagBlockInterval = "block-interval"
	FlagGenHeight   = "gen-height"
)

var (
	FsRequestRand = flag.NewFlagSet("", flag.ContinueOnError)
	FsQueryRand   = flag.NewFlagSet("", flag.ContinueOnError)
	FsQueryQueue  = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsRequestRand.Uint64(FlagBlockInterval, rand.DefaultBlockInterval, "the block interval")
	FsQueryRand.String(FlagReqID, "", "the request id")
	FsQueryQueue.Int64(FlagGenHeight, 0, "optional height at which rands are generated")
}