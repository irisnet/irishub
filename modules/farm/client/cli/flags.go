// nolint
package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagDescription      = "description"
	FlagStartHeight      = "start-height"
	FlagRewardPerBlock   = "reward-per-block"
	FlagLPTokenDenom     = "lp-token-denom"
	FlagTotalReward      = "total-reward"
	FlagEditable         = "editable"
	FlagFarmPool         = "farm-pool"
	FlagAdditionalReward = "additional-reward"
)

// common flag sets to add to various functions
var (
	FsCreateFarmPool = flag.NewFlagSet("", flag.ContinueOnError)
	FsAdjustFarmPool = flag.NewFlagSet("", flag.ContinueOnError)
	FsQueryFarmPool  = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsCreateFarmPool.String(FlagDescription, "", "The simple description of a farm pool")
	FsCreateFarmPool.Int64(FlagStartHeight, 0, "The start height the farm pool")
	FsCreateFarmPool.String(FlagRewardPerBlock, "", "The reward per block,ex: 1iris,1atom")
	FsCreateFarmPool.String(FlagLPTokenDenom, "", "The token accepted by farm pool")
	FsCreateFarmPool.String(FlagTotalReward, "", "The Total reward for the farm pool")
	FsCreateFarmPool.Bool(FlagEditable, false, "Is it possible to adjust the parameters of the farm pool")

	FsAdjustFarmPool.String(FlagAdditionalReward, "", "Bonuses added to the farm pool")
	FsAdjustFarmPool.String(FlagRewardPerBlock, "", "The reward per block,ex: 1iris,1atom")

	FsQueryFarmPool.String(FlagFarmPool, "", "The farm pool name")
}
