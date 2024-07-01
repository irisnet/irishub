package types

import (
	fmt "fmt"
	"strings"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"gopkg.in/yaml.v2"
)

const (
	FormatHTLTAssetPrefix = "htlt"
)

// NewParams is the HTLC params constructor
func NewParams(assetParams []AssetParam) Params {
	return Params{
		AssetParams: assetParams,
	}
}

// DefaultParams returns the default coinswap module parameters
func DefaultParams() Params {
	return Params{[]AssetParam{}}
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func NewAssetParam(
	denom string, coinID int, limit SupplyLimit, active bool,
	deputyAddr string, fixedFee, minSwapAmount math.Int,
	maxSwapAmount math.Int, minBlockLock, maxBlockLock uint64,
) AssetParam {
	return AssetParam{
		Denom:         denom,
		SupplyLimit:   limit,
		Active:        active,
		DeputyAddress: deputyAddr,
		FixedFee:      fixedFee,
		MinSwapAmount: minSwapAmount,
		MaxSwapAmount: maxSwapAmount,
		MinBlockLock:  minBlockLock,
		MaxBlockLock:  maxBlockLock,
	}
}

// String returns a human readable string representation of the parameters.
func (p AssetParam) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// String returns a human readable string representation of the parameters.
func (p SupplyLimit) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Validate returns err if Params is invalid
func (p Params) Validate() error {
	return validateAssetParams(p.AssetParams)
}

func validateAssetParams(i interface{}) error {
	assetParams, ok := i.([]AssetParam)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	coinDenoms := make(map[string]bool)
	for _, asset := range assetParams {
		if err := sdk.ValidateDenom(asset.Denom); err != nil ||
			!strings.HasPrefix(asset.Denom, FormatHTLTAssetPrefix) ||
			strings.ToLower(asset.Denom) != asset.Denom ||
			len(asset.Denom) < MinDenomLength {
			return fmt.Errorf(fmt.Sprintf("invalid asset denom: %s", asset.Denom))
		}

		if asset.SupplyLimit.Limit.IsNegative() {
			return fmt.Errorf(
				fmt.Sprintf(
					"asset %s has invalid (negative) supply limit: %s",
					asset.Denom,
					asset.SupplyLimit.Limit,
				),
			)
		}

		if asset.SupplyLimit.TimeBasedLimit.IsNegative() {
			return fmt.Errorf(
				fmt.Sprintf(
					"asset %s has invalid (negative) supply time limit: %s",
					asset.Denom,
					asset.SupplyLimit.TimeBasedLimit,
				),
			)
		}

		if asset.SupplyLimit.TimeBasedLimit.GT(asset.SupplyLimit.Limit) {
			return fmt.Errorf(
				fmt.Sprintf(
					"asset %s cannot have supply time limit greater than supply limit: %s>%s",
					asset.Denom,
					asset.SupplyLimit.TimeBasedLimit,
					asset.SupplyLimit.Limit,
				),
			)
		}

		if _, found := coinDenoms[asset.Denom]; found {
			return fmt.Errorf(fmt.Sprintf("asset %s cannot have duplicate denom", asset.Denom))
		}

		coinDenoms[asset.Denom] = true

		if _, err := sdk.AccAddressFromBech32(asset.DeputyAddress); err != nil {
			return fmt.Errorf("invalid deputy address %s", asset.DeputyAddress)
		}

		if asset.FixedFee.IsNegative() {
			return fmt.Errorf(
				"asset %s cannot have a negative fixed fee %s",
				asset.Denom,
				asset.FixedFee,
			)
		}

		if asset.MinBlockLock < MinTimeLock {
			return fmt.Errorf(
				"asset %s has minimum time lock %d less than min htlc time lock %d",
				asset.Denom,
				asset.MinBlockLock,
				MinTimeLock,
			)
		}

		if asset.MaxBlockLock > MaxTimeLock {
			return fmt.Errorf(
				"asset %s has maximum time lock %d  greater than max htlc time lock %d",
				asset.Denom,
				asset.MaxBlockLock,
				MaxTimeLock,
			)
		}

		if asset.MinBlockLock > asset.MaxBlockLock {
			return fmt.Errorf(
				"asset %s has minimum time lock %d greater than maximum time lock %d",
				asset.Denom,
				asset.MinBlockLock,
				asset.MaxBlockLock,
			)
		}

		if !asset.MinSwapAmount.IsPositive() {
			return fmt.Errorf(
				fmt.Sprintf(
					"asset %s must have a positive minimum swap amount, got %s",
					asset.Denom,
					asset.MinSwapAmount,
				),
			)
		}

		if !asset.MaxSwapAmount.IsPositive() {
			return fmt.Errorf(
				fmt.Sprintf(
					"asset %s must have a positive maximum swap amount, got %s",
					asset.Denom,
					asset.MaxSwapAmount,
				),
			)
		}

		if asset.MinSwapAmount.GT(asset.MaxSwapAmount) {
			return fmt.Errorf(
				"asset %s has minimum swap amount %s greater than maximum swap amount %s",
				asset.Denom,
				asset.MinSwapAmount,
				asset.MaxSwapAmount,
			)
		}
	}

	return nil
}
