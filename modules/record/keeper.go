package record

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/codec"
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

func KeyRecord(dataHash string) []byte {
	return []byte(fmt.Sprintf("record:%s", dataHash))
}

func (keeper Keeper) AddRecord(ctx sdk.Context, msg MsgSubmitRecord) {

	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(msg)
	store.Set(KeyRecord(msg.DataHash), bz)
}
