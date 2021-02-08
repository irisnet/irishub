package cli

import (
	flag "github.com/spf13/pflag"

	randomtypes "github.com/irisnet/irismod/modules/random/types"
)

const (
	FlagReqID         = "request-id"
	FlagBlockInterval = "block-interval"
	FlagOracle        = "oracle"
	FlagServiceFeeCap = "service-fee-cap"
	FlagQueueHeight   = "queue-height"
)

var (
	FsRequestRand = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsRequestRand.Uint64(FlagBlockInterval, randomtypes.DefaultBlockInterval, "The block interval")
	FsRequestRand.Bool(FlagOracle, false, "Indicate if the request is initiated with oracle")
	FsRequestRand.String(FlagServiceFeeCap, "", "Maximum fee to pay for a service request")
}
