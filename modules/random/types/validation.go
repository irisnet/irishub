package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ValidateServiceFeeCap verifies whether the  parameters are legal
func ValidateServiceFeeCap(serviceFeeCap sdk.Coins) error {
	if !serviceFeeCap.IsValid() {
		return sdkerrors.Wrapf(ErrInvalidServiceFeeCap, serviceFeeCap.String())
	}
	return nil
}
