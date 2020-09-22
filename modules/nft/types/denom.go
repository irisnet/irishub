package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewDenom return a new denom
func NewDenom(id, name, schema string, creator sdk.AccAddress) Denom {
	return Denom{
		Id:      id,
		Name:    name,
		Schema:  schema,
		Creator: creator,
	}
}

func ValidateDenomID(denomID string) error {
	denomID = strings.TrimSpace(denomID)
	if len(denomID) < MinDenomLen || len(denomID) > MaxDenomLen {
		return sdkerrors.Wrapf(ErrInvalidDenom, "invalid denom %s, only accepts value [%d, %d]", denomID, MinDenomLen, MaxDenomLen)
	}
	if !IsBeginWithAlpha(denomID) || !IsAlphaNumeric(denomID) {
		return sdkerrors.Wrapf(ErrInvalidDenom, "invalid denom %s, only accepts alphanumeric characters,and begin with an english letter", denomID)
	}
	return nil
}
