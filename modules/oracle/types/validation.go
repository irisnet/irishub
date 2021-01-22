package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	MaxLatestHistory    = 100
	MaxNameLen          = 70
	MaxAggregateFuncLen = 10
	MaxValueJsonPath    = 70
	MaxDescriptionLen   = 280
)

// ValidateFeedName verify that the feedName is legal
func ValidateFeedName(feedName string) error {
	if len(feedName) == 0 || len(feedName) > MaxNameLen {
		return sdkerrors.Wrap(ErrInvalidFeedName, feedName)
	}

	if !regPlainText.MatchString(feedName) {
		return sdkerrors.Wrap(ErrInvalidFeedName, feedName)
	}
	return nil
}

// ValidateDescription verify that the desc is legal
func ValidateDescription(desc string) error {
	if len(desc) > MaxDescriptionLen {
		return sdkerrors.Wrap(ErrInvalidDescription, desc)
	}
	return nil
}

// ValidateAggregateFunc verify that the aggregateFunc is legal
func ValidateAggregateFunc(aggregateFunc string) error {
	if len(aggregateFunc) == 0 || len(aggregateFunc) > MaxAggregateFuncLen {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "aggregate func must between [1, %d], got: %d", MaxAggregateFuncLen, len(aggregateFunc))
	}

	if _, err := GetAggregateFunc(aggregateFunc); err != nil {
		return err
	}
	return nil
}

// ValidateValueJSONPath verify that the valueJsonPath is legal
func ValidateValueJSONPath(valueJSONPath string) error {
	valueJSONPath = strings.TrimSpace(valueJSONPath)
	if len(valueJSONPath) == 0 || len(valueJSONPath) > MaxValueJsonPath {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "the length of valueJson path func must less than %d, got: %d", MaxAggregateFuncLen, len(valueJSONPath))
	}
	return nil
}

// ValidateLatestHistory verify that the latestHistory is legal
func ValidateLatestHistory(latestHistory uint64) error {
	if latestHistory < 1 || latestHistory > MaxLatestHistory {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "latest history is invalid, should be between 1 and %d", MaxLatestHistory)
	}
	return nil
}

// ValidateCreator verify that the creator is legal
func ValidateCreator(creator string) error {
	if _, err := sdk.AccAddressFromBech32(creator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator")
	}
	return nil
}

// ValidateServiceName verifies whether the  parameters are legal
func ValidateServiceName(serviceName string) error {
	if len(serviceName) == 0 || len(serviceName) > MaxNameLen {
		return sdkerrors.Wrapf(ErrInvalidServiceName, serviceName)
	}
	if !regPlainText.MatchString(serviceName) {
		return sdkerrors.Wrapf(ErrInvalidServiceName, serviceName)
	}
	return nil
}

// ValidateResponseThreshold verifies whether the  parameters are legal
func ValidateResponseThreshold(responseThreshold uint32, maxCnt int) error {
	if (maxCnt != 0 && int(responseThreshold) > maxCnt) || responseThreshold < 1 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "response threshold should be between 1 and %d", maxCnt)
	}
	return nil
}

// ValidateTimeout verifies whether the  parameters are legal
func ValidateTimeout(timeout int64, frequency uint64) error {
	if frequency < uint64(timeout) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "timeout [%d] should be no more than frequency [%d]", timeout, frequency)
	}
	return nil
}

// ValidateServiceFeeCap verifies whether the  parameters are legal
func ValidateServiceFeeCap(serviceFeeCap sdk.Coins) error {
	if !serviceFeeCap.IsValid() {
		return sdkerrors.Wrapf(ErrInvalidServiceFeeCap, serviceFeeCap.String())
	}
	return nil
}

// Modified return whether the  parameters are modified
func Modified(target string) bool {
	return target != DoNotModify
}
