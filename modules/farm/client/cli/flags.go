// nolint
package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagDescription    = "description"
	FlagStartHeight    = "start-height"
	FlagRewardPerBlock = "reward-per-block"
	FlagLPTokenDenom   = "lp-token-denom"
	FlagTotalReward    = "total-reward"
	FlagDestructible   = "destructible"
	FlagFarmPool       = "farm-pool"
)

// common flag sets to add to various functions
var (
	FsCreateFarmPool = flag.NewFlagSet("", flag.ContinueOnError)
	FsQueryFarmPool  = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsCreateFarmPool.String(FlagDescription, "", "The simple description of a farm pool")
	FsCreateFarmPool.Int64(FlagStartHeight, 0, "The start height the farm pool ")
	FsCreateFarmPool.String(FlagRewardPerBlock, "", "The reward per block,ex: 1iris,1atom")
	FsCreateFarmPool.String(FlagLPTokenDenom, "", "The token accepted by farm pool")
	FsCreateFarmPool.String(FlagTotalReward, "", "The Total reward for the farm pool")
	FsCreateFarmPool.Bool(FlagDestructible, false, "Can farm activities end early")

	FsQueryFarmPool.String(FlagFarmPool, "", "The farm pool name")
}
