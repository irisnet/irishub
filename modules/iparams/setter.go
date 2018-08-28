package iparams

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Setter interface {
	Set(ctx sdk.Context, key string, ptr interface{}) error
	SetRaw(ctx sdk.Context, key string, param []byte)
	SetString(ctx sdk.Context, key string, param string)
}

// Setter exposes all methods including Set
type GlobalSetter struct {
	GlobalGetter
}

// Setter exposes all methods including Set
type GovSetter struct {
	GovGetter
}

// Set exposes set
func (k GlobalSetter) Set(ctx sdk.Context, key string, param interface{}) error {
	key = getGlobalStoreKey(key)
	return k.k.set(ctx, key, param)
}

// SetRaw exposes setRaw
func (k GlobalSetter) SetRaw(ctx sdk.Context, key string, param []byte) {
	key = getGlobalStoreKey(key)
	k.k.setRaw(ctx, key, param)
}

// SetString is helper function for string params
func (k GlobalSetter) SetString(ctx sdk.Context, key string, param string) {
	key = getGlobalStoreKey(key)
	if err := k.k.set(ctx, key, param); err != nil {
		panic(err)
	}
}

// SetBool is helper function for bool params
func (k GlobalSetter) SetBool(ctx sdk.Context, key string, param bool) {
	key = getGlobalStoreKey(key)
	if err := k.k.set(ctx, key, param); err != nil {
		panic(err)
	}
}

// SetInt16 is helper function for int16 params
func (k GlobalSetter) SetInt16(ctx sdk.Context, key string, param int16) {
	key = getGlobalStoreKey(key)
	if err := k.k.set(ctx, key, param); err != nil {
		panic(err)
	}
}

// SetInt32 is helper function for int32 params
func (k GlobalSetter) SetInt32(ctx sdk.Context, key string, param int32) {
	key = getGlobalStoreKey(key)
	if err := k.k.set(ctx, key, param); err != nil {
		panic(err)
	}
}

// SetInt64 is helper function for int64 params
func (k GlobalSetter) SetInt64(ctx sdk.Context, key string, param int64) {
	key = getGlobalStoreKey(key)
	if err := k.k.set(ctx, key, param); err != nil {
		panic(err)
	}
}

// SetUint16 is helper function for uint16 params
func (k GlobalSetter) SetUint16(ctx sdk.Context, key string, param uint16) {
	key = getGlobalStoreKey(key)
	if err := k.k.set(ctx, key, param); err != nil {
		panic(err)
	}
}

// SetUint32 is helper function for uint32 params
func (k GlobalSetter) SetUint32(ctx sdk.Context, key string, param uint32) {
	key = getGlobalStoreKey(key)
	if err := k.k.set(ctx, key, param); err != nil {
		panic(err)
	}
}

// SetUint64 is helper function for uint64 params
func (k GlobalSetter) SetUint64(ctx sdk.Context, key string, param uint64) {
	key = getGlobalStoreKey(key)
	if err := k.k.set(ctx, key, param); err != nil {
		panic(err)
	}
}

// SetInt is helper function for sdk.Int params
func (k GlobalSetter) SetInt(ctx sdk.Context, key string, param sdk.Int) {
	key = getGlobalStoreKey(key)
	if err := k.k.set(ctx, key, param); err != nil {
		panic(err)
	}
}

// SetUint is helper function for sdk.Uint params
func (k GlobalSetter) SetUint(ctx sdk.Context, key string, param sdk.Uint) {
	key = getGlobalStoreKey(key)
	if err := k.k.set(ctx, key, param); err != nil {
		panic(err)
	}
}

// SetRat is helper function for rat params
func (k GlobalSetter) SetRat(ctx sdk.Context, key string, param sdk.Rat) {
	key = getGlobalStoreKey(key)
	if err := k.k.set(ctx, key, param); err != nil {
		panic(err)
	}
}


// Set exposes set
func (k GovSetter) Set(ctx sdk.Context, key string, param interface{}) error {
	key = getGovStoreKey(key)
	return k.k.set(ctx, key, param)
}

// SetRaw exposes setRaw
func (k GovSetter) SetRaw(ctx sdk.Context, key string, param []byte) {
	key = getGovStoreKey(key)
	k.k.setRaw(ctx, key, param)
}

// SetString is helper function for string params
func (k GovSetter) SetString(ctx sdk.Context, key string, param string) {
	key = getGovStoreKey(key)
	if err := k.k.set(ctx, key, param); err != nil {
		panic(err)
	}
}

// SetBool is helper function for bool params
func (k GovSetter) SetBool(ctx sdk.Context, key string, param bool) {
	key = getGovStoreKey(key)
	if err := k.k.set(ctx, key, param); err != nil {
		panic(err)
	}
}

// SetInt16 is helper function for int16 params
func (k GovSetter) SetInt16(ctx sdk.Context, key string, param int16) {
	key = getGovStoreKey(key)
	if err := k.k.set(ctx, key, param); err != nil {
		panic(err)
	}
}

// SetInt32 is helper function for int32 params
func (k GovSetter) SetInt32(ctx sdk.Context, key string, param int32) {
	key = getGovStoreKey(key)
	if err := k.k.set(ctx, key, param); err != nil {
		panic(err)
	}
}

// SetInt64 is helper function for int64 params
func (k GovSetter) SetInt64(ctx sdk.Context, key string, param int64) {
	key = getGovStoreKey(key)
	if err := k.k.set(ctx, key, param); err != nil {
		panic(err)
	}
}

// SetUint16 is helper function for uint16 params
func (k GovSetter) SetUint16(ctx sdk.Context, key string, param uint16) {
	key = getGovStoreKey(key)
	if err := k.k.set(ctx, key, param); err != nil {
		panic(err)
	}
}

// SetUint32 is helper function for uint32 params
func (k GovSetter) SetUint32(ctx sdk.Context, key string, param uint32) {
	key = getGovStoreKey(key)
	if err := k.k.set(ctx, key, param); err != nil {
		panic(err)
	}
}

// SetUint64 is helper function for uint64 params
func (k GovSetter) SetUint64(ctx sdk.Context, key string, param uint64) {
	key = getGovStoreKey(key)
	if err := k.k.set(ctx, key, param); err != nil {
		panic(err)
	}
}

// SetInt is helper function for sdk.Int params
func (k GovSetter) SetInt(ctx sdk.Context, key string, param sdk.Int) {
	key = getGovStoreKey(key)
	if err := k.k.set(ctx, key, param); err != nil {
		panic(err)
	}
}

// SetUint is helper function for sdk.Uint params
func (k GovSetter) SetUint(ctx sdk.Context, key string, param sdk.Uint) {
	key = getGovStoreKey(key)
	if err := k.k.set(ctx, key, param); err != nil {
		panic(err)
	}
}

// SetRat is helper function for rat params
func (k GovSetter) SetRat(ctx sdk.Context, key string, param sdk.Rat) {
	key = getGovStoreKey(key)
	if err := k.k.set(ctx, key, param); err != nil {
		panic(err)
	}
}
