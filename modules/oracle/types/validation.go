package types

import (
	"regexp"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"mods.irisnet.org/modules/service/exported"
)

const (
	// MaxLatestHistory defines the the maximum number of feed values saved
	MaxLatestHistory = 100
	// MaxAggregateFuncNameLen defines the the maximum length of the aggregation function name
	MaxAggregateFuncNameLen = 10
	// MaxDescriptionLen defines the the maximum length of the description
	MaxDescriptionLen = 280
)

// the feed name only accepts alphanumeric characters, _ and - /
var regexpFeedName = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9/_-]*$`)

// ValidateFeedName verifies if the feed name is legal
func ValidateFeedName(feedName string) error {
	if !regexpFeedName.MatchString(feedName) {
		return errorsmod.Wrap(ErrInvalidFeedName, feedName)
	}
	return nil
}

// ValidateDescription verifies if the description is legal
func ValidateDescription(desc string) error {
	if len(desc) > MaxDescriptionLen {
		return errorsmod.Wrap(ErrInvalidDescription, desc)
	}
	return nil
}

// ValidateAggregateFunc verifies if the aggregation function is legal
func ValidateAggregateFunc(aggregateFunc string) error {
	if len(aggregateFunc) == 0 || len(aggregateFunc) > MaxAggregateFuncNameLen {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "aggregate func must between [1, %d], got: %d", MaxAggregateFuncNameLen, len(aggregateFunc))
	}

	if _, err := GetAggregateFunc(aggregateFunc); err != nil {
		return err
	}
	return nil
}

// ValidateLatestHistory verifies if the latest history is legal
func ValidateLatestHistory(latestHistory uint64) error {
	if latestHistory < 1 || latestHistory > MaxLatestHistory {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "latest history is invalid, should be between 1 and %d", MaxLatestHistory)
	}
	return nil
}

// ValidateCreator verifies if the creator is legal
func ValidateCreator(creator string) error {
	if _, err := sdk.AccAddressFromBech32(creator); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator")
	}
	return nil
}

// ValidateServiceName verifies whether the service name is legal
func ValidateServiceName(serviceName string) error {
	return exported.ValidateServiceName(serviceName)
}

// ValidateResponseThreshold verifies whether the given threshold is legal
func ValidateResponseThreshold(responseThreshold uint32, maxCnt int) error {
	if (maxCnt != 0 && int(responseThreshold) > maxCnt) || responseThreshold < 1 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "response threshold should be between 1 and %d", maxCnt)
	}
	return nil
}

// ValidateTimeout verifies whether the given timeout and frequency are legal
func ValidateTimeout(timeout int64, frequency uint64) error {
	if frequency < uint64(timeout) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "timeout [%d] should be no more than frequency [%d]", timeout, frequency)
	}
	return nil
}

// ValidateServiceFeeCap verifies whether the given service fee cap is legal
func ValidateServiceFeeCap(serviceFeeCap sdk.Coins) error {
	if !serviceFeeCap.IsValid() {
		return errorsmod.Wrapf(ErrInvalidServiceFeeCap, serviceFeeCap.String())
	}
	return nil
}

// Modified returns true if the given target string is modified
// False otherwise
func Modified(target string) bool {
	return target != DoNotModify
}
