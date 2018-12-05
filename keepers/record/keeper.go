package record

import (
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/codec"
	types "github.com/irisnet/irishub/types/record"
)

// nolint

// Record Keeper
type Keeper struct {

	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey

	// The codec codec for binary encoding/decoding.
	cdc *codec.Codec

	// Reserved codespace
	codespace sdk.CodespaceType
}

// NewKeeper returns a mapper that uses go-codec to (binary) encode and decode record types.
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, codespace sdk.CodespaceType) Keeper {
	return Keeper{
		storeKey:  key,
		cdc:       cdc,
		codespace: codespace,
	}
}

// Returns the go-codec codec.
func (keeper Keeper) WireCodec() *codec.Codec {
	return keeper.cdc
}

func (keeper Keeper) AddRecord(ctx sdk.Context, msg types.MsgSubmitRecord) {

	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(msg)
	store.Set(types.KeyRecord(msg.DataHash), bz)
}
