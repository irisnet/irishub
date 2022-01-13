package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewDenom return a new denom
func NewDenom(
	id, name, schema, symbol, description, uri, uriHash, data string,
	creator sdk.AccAddress,
	mintRestricted, updateRestricted bool,
) Denom {
	return Denom{
		Id:               id,
		Name:             name,
		Schema:           schema,
		Creator:          creator.String(),
		Symbol:           symbol,
		MintRestricted:   mintRestricted,
		UpdateRestricted: updateRestricted,
		Description:      description,
		Uri:              uri,
		UriHash:          uriHash,
		Data:             data,
	}
}
