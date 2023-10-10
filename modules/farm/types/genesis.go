package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	defaultGenesisState = GenesisState{
		Params: Params{
			PoolCreationFee:     DefaultPoolCreationFee,
			TaxRate:             DefaulttTaxRate,
			MaxRewardCategories: DefaultMaxRewardCategories,
		},
	}
)

// NewGenesisState constructs a new GenesisState instance
func NewGenesisState(params Params, pools []FarmPool, farmInfos []FarmInfo, sequence uint64, escrow []EscrowInfo) *GenesisState {
	return &GenesisState{
		params, pools, farmInfos, sequence, escrow,
	}
}

// DefaultGenesisState gets the default genesis state for testing
func DefaultGenesisState() *GenesisState {
	return &defaultGenesisState
}

func SetDefaultGenesisState(state GenesisState) {
	defaultGenesisState = state
}

// ValidateGenesis validates the provided farm genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	var maxSeq uint64
	for _, pool := range data.Pools {
		seq, err := ValidatepPoolId(pool.Id)
		if err != nil {
			return err
		}

		if seq > maxSeq {
			maxSeq = seq
		}

		if err := ValidateDescription(pool.Description); err != nil {
			return err
		}

		if err := ValidateAddress(pool.Creator); err != nil {
			return err
		}

		if err := ValidateCoins("TotalLptLocked", pool.TotalLptLocked); err != nil {
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

			// If the unexpired pool rule has been updated, rewardPerShare will not be zero.
			if !r.RewardPerShare.IsPositive() {
				// No reward has ever been distributed.
				if r.RemainingReward.Equal(r.TotalReward) {
					continue
				}
				// The pool is expired and the reward is refund to the creator
				if !r.RemainingReward.Equal(r.TotalReward) && pool.EndHeight == pool.LastHeightDistrRewards {
					continue
				}

				return fmt.Errorf("rewardPerShare must be positive, but got %s", r.RewardPerShare.String())
			}
		}
	}
	if data.Sequence < maxSeq {
		return fmt.Errorf("sequence must be equeal or greater than maxSeq, but got %d, %d", data.Sequence, maxSeq)
	}

	for _, info := range data.FarmInfos {
		if _, err := ValidatepPoolId(info.PoolId); err != nil {
			return err
		}

		if err := ValidateAddress(info.Address); err != nil {
			return err
		}

		if !info.Locked.IsPositive() {
			return fmt.Errorf("locked must be positive, but got %s", info.Locked.String())
		}

		if err := sdk.NewCoins(info.RewardDebt...).Validate(); err != nil {
			return err
		}
	}

	return ValidateCoins("PoolCreationFee", data.Params.PoolCreationFee)
}
