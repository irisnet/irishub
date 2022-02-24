package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagDescription         = "description"
	FlagStartHeight         = "start-height"
	FlagRewardPerBlock      = "reward-per-block"
	FlagLPTokenDenom        = "lp-token-denom"
	FlagTotalReward         = "total-reward"
	FlagEditable            = "editable"
	FlagFarmPool            = "pool-id"
	FlagAdditionalReward    = "additional-reward"
	FlagProposalDescription = "proposal-description"
	FlagProposalTitle       = "proposal-title"
	FlagProposaldeposit     = "proposal-deposit"
	FlagFundApplied         = "proposal-fund-applied"
	FlagFundSelfBond        = "proposal-fund-self-bond"
)

// common flag sets to add to various functions
var (
	FsCreateFarmPool              = flag.NewFlagSet("", flag.ContinueOnError)
	FsCreatePoolWithCommunityPool = flag.NewFlagSet("", flag.ContinueOnError)
	FsAdjustFarmPool              = flag.NewFlagSet("", flag.ContinueOnError)
	FsQueryFarmPool               = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsCreateFarmPool.String(FlagDescription, "", "The simple description of a farm pool")
	FsCreateFarmPool.Int64(FlagStartHeight, 0, "The start height the farm pool")
	FsCreateFarmPool.String(FlagRewardPerBlock, "", "The reward per block,ex: 1iris,1atom")
	FsCreateFarmPool.String(FlagLPTokenDenom, "", "The token accepted by farm pool")
	FsCreateFarmPool.String(FlagTotalReward, "", "The Total reward for the farm pool")
	FsCreateFarmPool.Bool(FlagEditable, false, "Is it possible to adjust the parameters of the farm pool")

	FsCreatePoolWithCommunityPool.String(FlagDescription, "", "The simple description of a farm pool")
	FsCreatePoolWithCommunityPool.String(FlagProposalDescription, "", "The simple description of a proposal")
	FsCreatePoolWithCommunityPool.String(FlagProposalTitle, "", "The simple title of a proposal")
	FsCreatePoolWithCommunityPool.String(FlagProposaldeposit, "", "The initial deposit of a proposal")
	FsCreatePoolWithCommunityPool.String(FlagFundApplied, "", "The fund applied from the community pool")
	FsCreatePoolWithCommunityPool.String(FlagFundSelfBond, "", "The fund self bond")
	FsCreatePoolWithCommunityPool.String(FlagRewardPerBlock, "", "The reward per block,ex: 1iris,1atom")
	FsCreatePoolWithCommunityPool.String(FlagLPTokenDenom, "", "The token accepted by farm pool")

	FsAdjustFarmPool.String(FlagAdditionalReward, "", "Bonuses added to the farm pool")
	FsAdjustFarmPool.String(FlagRewardPerBlock, "", "The reward per block,ex: 1iris,1atom")

	FsQueryFarmPool.String(FlagFarmPool, "", "The farm pool id")
}
