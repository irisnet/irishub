package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenesisState constructs a new GenesisState instance
func NewGenesisState(params Params, pools []FarmPool, farmInfos []FarmInfo) *GenesisState {
	return &GenesisState{
		params, pools, farmInfos,
	}
}

// DefaultGenesisState gets the default genesis state for testing
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: Params{
			CreatePoolFee:      sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(5000)),
			MaxRewardCategoryN: 2,
		},
	}
}

// ValidateGenesis validates the provided farm genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	for _, pool := range data.Pools {
		if err := ValidatePoolName(pool.Name); err != nil {
			return err
		}

		if err := ValidateDescription(pool.Description); err != nil {
			return err
		}

		if err := ValidateAddress(pool.Creator); err != nil {
			return err
		}

		if err := ValidateCoins(pool.TotalLpTokenLocked); err != nil {
			return err
		}

		for _, r := range pool.Rules {
			if err := ValidateLpTokenDenom(r.Reward); err != nil {
				return err
			}

			if !r.TotalReward.IsPositive() {
				return fmt.Errorf("totalReward must be positive, but got %s", r.TotalReward.String())
			}

			if r.RemainingReward.IsNegative() {
				return fmt.Errorf("temainingReward must be greater than zero, but got %s", r.RemainingReward.String())
			}

			if !r.RewardPerBlock.IsPositive() {
				return fmt.Errorf("rewardPerBlock must be positive, but got %s", r.RewardPerBlock.String())
			}

			if !r.RewardPerShare.IsPositive() {
				return fmt.Errorf("rewardPerShare must be positive, but got %s", r.RewardPerShare.String())
			}
		}
	}

	for _, info := range data.FarmInfos {
		if err := ValidatePoolName(info.PoolName); err != nil {
			return err
		}

		if err := ValidateAddress(info.Address); err != nil {
			return err
		}

		if !info.Locked.IsPositive() {
			return fmt.Errorf("locked must be positive, but got %s", info.Locked.String())
		}

		if err := ValidateCoins(info.RewardDebt...); err != nil {
			return err
		}
	}

	return ValidateCoins(data.Params.CreatePoolFee)
}
