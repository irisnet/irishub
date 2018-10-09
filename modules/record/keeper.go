package record

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

// nolint

// Record Keeper
type Keeper struct {
	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey

	// The wire codec for binary encoding/decoding.
	cdc *wire.Codec

	// Reserved codespace
	codespace sdk.CodespaceType
}

// NewKeeper returns a mapper that uses go-wire to (binary) encode and decode record types.
func NewKeeper(cdc *wire.Codec, key sdk.StoreKey, codespace sdk.CodespaceType) Keeper {
	return Keeper{
		storeKey:  key,
		cdc:       cdc,
		codespace: codespace,
	}
}

// Returns the go-wire codec.
func (keeper Keeper) WireCodec() *wire.Codec {
	return keeper.cdc
}
