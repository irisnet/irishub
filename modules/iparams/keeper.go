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


// Setter exposes all methods including Set
type Setter struct {
	Getter
}

// Set exposes set
func (k Setter) Set(ctx sdk.Context, key string, param interface{}) error {
	return k.k.set(ctx, key, param)
}

// SetRaw exposes setRaw
func (k Setter) SetRaw(ctx sdk.Context, key string, param []byte) {
	k.k.setRaw(ctx, key, param)
}

// SetString is helper function for string params
func (k Setter) SetString(ctx sdk.Context, key string, param string) {
	if err := k.k.set(ctx, key, param); err != nil {
		panic(err)
	}
}

// SetBool is helper function for bool params
func (k Setter) SetBool(ctx sdk.Context, key string, param bool) {
	if err := k.k.set(ctx, key, param); err != nil {
		panic(err)
	}
}

// SetInt16 is helper function for int16 params
func (k Setter) SetInt16(ctx sdk.Context, key string, param int16) {
	if err := k.k.set(ctx, key, param); err != nil {
		panic(err)
	}
}

// SetInt32 is helper function for int32 params
func (k Setter) SetInt32(ctx sdk.Context, key string, param int32) {
	if err := k.k.set(ctx, key, param); err != nil {
		panic(err)
	}
}

// SetInt64 is helper function for int64 params
func (k Setter) SetInt64(ctx sdk.Context, key string, param int64) {
	if err := k.k.set(ctx, key, param); err != nil {
		panic(err)
	}
}

// SetUint16 is helper function for uint16 params
func (k Setter) SetUint16(ctx sdk.Context, key string, param uint16) {
	if err := k.k.set(ctx, key, param); err != nil {
		panic(err)
	}
}

// SetUint32 is helper function for uint32 params
func (k Setter) SetUint32(ctx sdk.Context, key string, param uint32) {
	if err := k.k.set(ctx, key, param); err != nil {
		panic(err)
	}
}

// SetUint64 is helper function for uint64 params
func (k Setter) SetUint64(ctx sdk.Context, key string, param uint64) {
	if err := k.k.set(ctx, key, param); err != nil {
		panic(err)
	}
}

// SetInt is helper function for sdk.Int params
func (k Setter) SetInt(ctx sdk.Context, key string, param sdk.Int) {
	if err := k.k.set(ctx, key, param); err != nil {
		panic(err)
	}
}

// SetUint is helper function for sdk.Uint params
func (k Setter) SetUint(ctx sdk.Context, key string, param sdk.Uint) {
	if err := k.k.set(ctx, key, param); err != nil {
		panic(err)
	}
}

// SetRat is helper function for rat params
func (k Setter) SetRat(ctx sdk.Context, key string, param sdk.Rat) {
	if err := k.k.set(ctx, key, param); err != nil {
		panic(err)
	}
}
