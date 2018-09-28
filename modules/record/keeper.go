package record

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

// Record Keeper
type Keeper struct {
	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey

	// The wire codec for binary encoding/decoding.
	cdc *wire.Codec
}

func NewKeeper(cdc *wire.Codec, key sdk.StoreKey) Keeper {
	return Keeper{
		storeKey: key,
		cdc:      cdc,
	}
}

// Returns the go-wire codec.
func (keeper Keeper) WireCodec() *wire.Codec {
	return keeper.cdc
}

func (keeper Keeper) NewRecord(ctx sdk.Context, addr string, time string, hash string, size string, node string) Record {
	return Record{}
}

func (k Keeper) Record(ctx sdk.Context, msg MsgRecord) {
	return
}
