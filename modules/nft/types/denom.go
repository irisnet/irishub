package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewDenom return a new denom
func NewDenom(id, name, schema string, creator sdk.AccAddress) Denom {
	return Denom{
		Id:      id,
		Name:    name,
		Schema:  schema,
		Creator: creator.String(),
	}
}
