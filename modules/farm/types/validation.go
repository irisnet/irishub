package types

import (
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// MaxNameLength length of the service name
	MaxPoolNameLength = 70
	// MaxDescriptionLength length of the service and author description
	MaxDescriptionLength = 280
)

var (
	// the pool name only accepts alphanumeric characters, _ and -, beginning with alpha character
	regexpPoolName = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)
)

// ValidatePoolName validates the pool name
func ValidatePoolName(poolName string) error {
	if !regexpPoolName.MatchString(poolName) || len(poolName) > MaxPoolNameLength {
		return sdkerrors.Wrap(ErrInvalidPoolName, poolName)
	}
	return nil
}

// ValidateDescription validates the pool name
func ValidateDescription(description string) error {
	if len(description) > MaxDescriptionLength {
		return sdkerrors.Wrap(ErrInvalidPoolName, description)
	}
	return nil
}

// ValidateLpTokenDenom validates the lp token denom
func ValidateLpTokenDenom(denom string) error {
	return sdk.ValidateDenom(denom)
}

// ValidateCoins validates the coin
func ValidateCoins(coins ...sdk.Coin) error {
	return sdk.NewCoins(coins...).Validate()
}

// ValidateAddress validates the address
func ValidateAddress(sender string) error {
	_, err := sdk.AccAddressFromBech32(sender)
	return err
}

// ValidateReward validates the coin
func ValidateReward(rewardPerBlock, totalReward sdk.Coins) error {
	if len(rewardPerBlock) != len(totalReward) {
		return sdkerrors.Wrapf(ErrNotMatch, "The length of rewardPerBlock and totalReward must be the same")
	}

	for _, r := range totalReward {
		if !rewardPerBlock.AmountOf(r.Denom).IsPositive() {
			return sdkerrors.Wrapf(ErrNotMatch, "rewardPerBlock and totalReward token types are inconsistent")
		}
	}
	return nil
}
