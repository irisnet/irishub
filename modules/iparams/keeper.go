package iparams

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

// Keeper manages global parameter store
type Keeper struct {
	cdc *wire.Codec
	key sdk.StoreKey
}

// NewKeeper constructs a new Keeper
func NewKeeper(cdc *wire.Codec, key sdk.StoreKey) Keeper {
	return Keeper{
		cdc: cdc,
		key: key,
	}
}

// InitKeeper constructs a new Keeper with initial parameters
func InitKeeper(ctx sdk.Context, cdc *wire.Codec, key sdk.StoreKey, params ...interface{}) Keeper {
	if len(params)%2 != 0 {
		panic("Odd params list length for InitKeeper")
	}

	k := NewKeeper(cdc, key)

	for i := 0; i < len(params); i += 2 {
		k.set(ctx, params[i].(string), params[i+1])
	}

	return k
}

// get automatically unmarshalls parameter to pointer
func (k Keeper) get(ctx sdk.Context, key string, ptr interface{}) error {
	store := ctx.KVStore(k.key)
	bz := store.Get([]byte(key))
	return k.cdc.UnmarshalBinary(bz, ptr)
}

// getRaw returns raw byte slice
func (k Keeper) getRaw(ctx sdk.Context, key string) []byte {
	store := ctx.KVStore(k.key)
	return store.Get([]byte(key))
}

// set automatically marshalls and type check parameter
func (k Keeper) set(ctx sdk.Context, key string, param interface{}) error {
	store := ctx.KVStore(k.key)
	bz := store.Get([]byte(key))
	if bz != nil {
		ptrty := reflect.PtrTo(reflect.TypeOf(param))
		ptr := reflect.New(ptrty).Interface()

		if k.cdc.UnmarshalBinary(bz, ptr) != nil {
			return fmt.Errorf("Type mismatch with stored param and provided param")
		}
	}

	bz, err := k.cdc.MarshalBinary(param)
	if err != nil {
		return err
	}
	store.Set([]byte(key), bz)

	return nil
}

// setRaw sets raw byte slice
func (k Keeper) setRaw(ctx sdk.Context, key string, param []byte) {
	store := ctx.KVStore(k.key)
	store.Set([]byte(key), param)
}

//// Getter returns readonly struct
//func (k Keeper) Getter() Getter {
//	return Getter{k}
//}
//
//// Setter returns read/write struct
//func (k Keeper) Setter() Setter {
//	return Setter{Getter{k}}
//}



