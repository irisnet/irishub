package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ValidateServiceFeeCap verifies whether the service fee cap is legal
func ValidateServiceFeeCap(serviceFeeCap sdk.Coins) error {
	if !serviceFeeCap.IsValid() {
		return errorsmod.Wrapf(ErrInvalidServiceFeeCap, serviceFeeCap.String())
	}
	return nil
}
