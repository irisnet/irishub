package types

import (
	fmt "fmt"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ValidateInput verifies whether the given input is legal
func ValidateInput(input Input) error {
	if !(input.Coin.IsValid() && input.Coin.IsPositive()) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid input (%s)", input.Coin.String())
	}

	if strings.HasPrefix(input.Coin.Denom, LptTokenPrefix) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid input denom, should not begin with (%s)", LptTokenPrefix)
	}

	if _, err := sdk.AccAddressFromBech32(input.Address); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid input address (%s)", err)
	}
	return nil
}

// ValidateOutput verifies whether the given output is legal
func ValidateOutput(output Output) error {
	if !(output.Coin.IsValid() && output.Coin.IsPositive()) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid output (%s)", output.Coin.String())
	}

	if strings.HasPrefix(output.Coin.Denom, LptTokenPrefix) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid output denom, should not begin with (%s)", LptTokenPrefix)
	}

	if _, err := sdk.AccAddressFromBech32(output.Address); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid output address (%s)", err)
	}
	return nil
}

// ValidateDeadline verifies whether the given deadline is legal
func ValidateDeadline(deadline int64) error {
	if deadline <= 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("deadline %d must be greater than 0", deadline))
	}
	return nil
}

// ValidateMaxToken verifies whether the maximum token is legal
func ValidateMaxToken(maxToken sdk.Coin) error {
	if !(maxToken.IsValid() && maxToken.IsPositive()) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid maxToken (%s)", maxToken.String())
	}

	if strings.HasPrefix(maxToken.Denom, LptTokenPrefix) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "max token must be non-liquidity token")
	}
	return nil
}

// ValidateToken verifies whether the exact token is legal
func ValidateToken(exactToken sdk.Coin) error {
	if !(exactToken.IsValid() && exactToken.IsPositive()) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid exactToken (%s)", exactToken.String())
	}

	if strings.HasPrefix(exactToken.Denom, LptTokenPrefix) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "token must be non-liquidity token")
	}
	return nil
}

// ValidateCounterpartyDenom verifies whether the counterparty denom is legal
func ValidateCounterpartyDenom(counterpartydenom string) error {
	if counterpartydenom == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "counterparty denom should not be empty")
	}
	return nil
}

// ValidateExactStandardAmt verifies whether the standard token amount is legal
func ValidateExactStandardAmt(standardAmt sdkmath.Int) error {
	if !standardAmt.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "standard token amount must be positive")
	}
	return nil
}

// ValidateLiquidity verifies whether the minimum liquidity is legal
func ValidateLiquidity(minLiquidity sdkmath.Int) error {
	if minLiquidity.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "liquidity can not be negative")
	}
	return nil
}

// ValidateMinToken verifies whether the minimum token amount is legal
func ValidateMinToken(minToken sdkmath.Int) error {
	if minToken.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "minimum token amount can not be negative")
	}
	return nil
}

// ValidateWithdrawLiquidity verifies whether the given liquidity is legal
func ValidateWithdrawLiquidity(liquidity sdk.Coin) error {
	if !liquidity.IsValid() || !liquidity.IsPositive() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid withdrawLiquidity (%s)", liquidity.String())
	}

	if err := ValidateLptDenom(liquidity.Denom); err != nil {
		return err
	}
	return nil
}

// ValidateMinStandardAmt verifies whether the minimum standard amount is legal
func ValidateMinStandardAmt(minStandardAmt sdkmath.Int) error {
	if minStandardAmt.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("minimum standard token amount %s can not be negative", minStandardAmt.String()))
	}
	return nil
}

// ValidateLptDenom returns nil if the Liquidity pool token denom is valid
func ValidateLptDenom(lptDenom string) error {
	if _, err := ParseLptDenom(lptDenom); err != nil {
		return sdkerrors.Wrap(ErrInvalidDenom, lptDenom)
	}
	return nil
}

// ValidatePoolSequenceId returns nil if the pool id is valid
func ValidatePoolSequenceId(poolId uint64) error {
	if poolId == 0 {
		return sdkerrors.Wrap(ErrReservePoolNotExists, "pool-id is not valid")
	}
	return nil
}
