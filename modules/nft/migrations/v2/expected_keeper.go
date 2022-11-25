package v2

import sdk "github.com/cosmos/cosmos-sdk/types"

// SaveDenom save the denom of class
type SaveDenom func(ctx sdk.Context, id, name, schema,
	symbol string,
	creator sdk.AccAddress,
	mintRestricted,
	updateRestricted bool,
	description,
	uri,
	uriHash,
	data string,
) error
