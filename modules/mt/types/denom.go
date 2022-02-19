package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewDenom return a new denom
func NewDenom(
	id, name string, data []byte, creator sdk.AccAddress,
) Denom {
	return Denom{
		Id:      id,
		Name:    name,
		Creator: creator.String(),
		Data:    data,
	}
}
