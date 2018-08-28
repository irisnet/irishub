package iparams

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

type Getter interface {
	Get(ctx sdk.Context, key string, ptr interface{}) error
	GetRaw(ctx sdk.Context, key string) []byte
	GetString(ctx sdk.Context, key string) (res string, err error)
	GetBoolWithDefault(ctx sdk.Context, key string, def bool) (res bool)
}

// Getter exposes methods related with only getting params
type GlobalGetter struct {
	k Keeper
}

type GovGetter struct {
	k Keeper
}

func getGlobalStoreKey(key string) string {
	if strings.HasPrefix(key, Global+"/") {
		return key
	}
	return fmt.Sprintf("%s/%s", Global, key)
}

func getGovStoreKey(key string) string {
	if strings.HasPrefix(key, Gov+"/") {
		return key
	}
	return fmt.Sprintf("%s/%s", Gov, key)
}

// Get exposes get
func (k GlobalGetter) Get(ctx sdk.Context, key string, ptr interface{}) error {
	key = getGlobalStoreKey(key)
	return k.k.get(ctx, key, ptr)
}

// GetRaw exposes getRaw
func (k GlobalGetter) GetRaw(ctx sdk.Context, key string) []byte {
	key = getGlobalStoreKey(key)
	return k.k.getRaw(ctx, key)
}

// GetString is helper function for string params
func (k GlobalGetter) GetString(ctx sdk.Context, key string) (res string, err error) {
	key = getGlobalStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	err = k.k.cdc.UnmarshalBinary(bz, &res)
	return
}

// GetBool is helper function for bool params
func (k GlobalGetter) GetBool(ctx sdk.Context, key string) (res bool, err error) {
	key = getGlobalStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	err = k.k.cdc.UnmarshalBinary(bz, &res)
	return
}

// GetInt16 is helper function for int16 params
func (k GlobalGetter) GetInt16(ctx sdk.Context, key string) (res int16, err error) {
	key = getGlobalStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	err = k.k.cdc.UnmarshalBinary(bz, &res)
	return
}

// GetInt32 is helper function for int32 params
func (k GlobalGetter) GetInt32(ctx sdk.Context, key string) (res int32, err error) {
	key = getGlobalStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	err = k.k.cdc.UnmarshalBinary(bz, &res)
	return
}

// GetInt64 is helper function for int64 params
func (k GlobalGetter) GetInt64(ctx sdk.Context, key string) (res int64, err error) {
	key = getGlobalStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	err = k.k.cdc.UnmarshalBinary(bz, &res)
	return
}

// GetUint16 is helper function for uint16 params
func (k GlobalGetter) GetUint16(ctx sdk.Context, key string) (res uint16, err error) {
	key = getGlobalStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	err = k.k.cdc.UnmarshalBinary(bz, &res)
	return
}

// GetUint32 is helper function for uint32 params
func (k GlobalGetter) GetUint32(ctx sdk.Context, key string) (res uint32, err error) {
	key = getGlobalStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	err = k.k.cdc.UnmarshalBinary(bz, &res)
	return
}

// GetUint64 is helper function for uint64 params
func (k GlobalGetter) GetUint64(ctx sdk.Context, key string) (res uint64, err error) {
	key = getGlobalStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	err = k.k.cdc.UnmarshalBinary(bz, &res)
	return
}

// GetInt is helper function for sdk.Int params
func (k GlobalGetter) GetInt(ctx sdk.Context, key string) (res sdk.Int, err error) {
	key = getGlobalStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	err = k.k.cdc.UnmarshalBinary(bz, &res)
	return
}

// GetUint is helper function for sdk.Uint params
func (k GlobalGetter) GetUint(ctx sdk.Context, key string) (res sdk.Uint, err error) {
	key = getGlobalStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	err = k.k.cdc.UnmarshalBinary(bz, &res)
	return
}

// GetRat is helper function for rat params
func (k GlobalGetter) GetRat(ctx sdk.Context, key string) (res sdk.Rat, err error) {
	key = getGlobalStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	err = k.k.cdc.UnmarshalBinary(bz, &res)
	return
}

// GetStringWithDefault is helper function for string params with default value
func (k GlobalGetter) GetStringWithDefault(ctx sdk.Context, key string, def string) (res string) {
	key = getGlobalStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	if bz == nil {
		return def
	}
	k.k.cdc.MustUnmarshalBinary(bz, &res)
	return
}

// GetBoolWithDefault is helper function for bool params with default value
func (k GlobalGetter) GetBoolWithDefault(ctx sdk.Context, key string, def bool) (res bool) {
	key = getGlobalStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	if bz == nil {
		return def
	}
	k.k.cdc.MustUnmarshalBinary(bz, &res)
	return
}

// GetInt16WithDefault is helper function for int16 params with default value
func (k GlobalGetter) GetInt16WithDefault(ctx sdk.Context, key string, def int16) (res int16) {
	key = getGlobalStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	if bz == nil {
		return def
	}
	k.k.cdc.MustUnmarshalBinary(bz, &res)
	return
}

// GetInt32WithDefault is helper function for int32 params with default value
func (k GlobalGetter) GetInt32WithDefault(ctx sdk.Context, key string, def int32) (res int32) {
	key = getGlobalStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	if bz == nil {
		return def
	}
	k.k.cdc.MustUnmarshalBinary(bz, &res)
	return
}

// GetInt64WithDefault is helper function for int64 params with default value
func (k GlobalGetter) GetInt64WithDefault(ctx sdk.Context, key string, def int64) (res int64) {
	key = getGlobalStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	if bz == nil {
		return def
	}
	k.k.cdc.MustUnmarshalBinary(bz, &res)
	return
}

// GetUint16WithDefault is helper function for uint16 params with default value
func (k GlobalGetter) GetUint16WithDefault(ctx sdk.Context, key string, def uint16) (res uint16) {
	key = getGlobalStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	if bz == nil {
		return def
	}
	k.k.cdc.MustUnmarshalBinary(bz, &res)
	return
}

// GetUint32WithDefault is helper function for uint32 params with default value
func (k GlobalGetter) GetUint32WithDefault(ctx sdk.Context, key string, def uint32) (res uint32) {
	key = getGlobalStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	if bz == nil {
		return def
	}
	k.k.cdc.MustUnmarshalBinary(bz, &res)
	return
}

// GetUint64WithDefault is helper function for uint64 params with default value
func (k GlobalGetter) GetUint64WithDefault(ctx sdk.Context, key string, def uint64) (res uint64) {
	key = getGlobalStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	if bz == nil {
		return def
	}
	k.k.cdc.MustUnmarshalBinary(bz, &res)
	return
}

// GetIntWithDefault is helper function for sdk.Int params with default value
func (k GlobalGetter) GetIntWithDefault(ctx sdk.Context, key string, def sdk.Int) (res sdk.Int) {
	key = getGlobalStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	if bz == nil {
		return def
	}
	k.k.cdc.MustUnmarshalBinary(bz, &res)
	return
}

// GetUintWithDefault is helper function for sdk.Uint params with default value
func (k GlobalGetter) GetUintWithDefault(ctx sdk.Context, key string, def sdk.Uint) (res sdk.Uint) {
	key = getGlobalStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	if bz == nil {
		return def
	}
	k.k.cdc.MustUnmarshalBinary(bz, &res)
	return
}

// GetRatWithDefault is helper function for sdk.Rat params with default value
func (k GlobalGetter) GetRatWithDefault(ctx sdk.Context, key string, def sdk.Rat) (res sdk.Rat) {
	key = getGlobalStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	if bz == nil {
		return def
	}
	k.k.cdc.MustUnmarshalBinary(bz, &res)
	return
}

// Get exposes get
func (k GovGetter) Get(ctx sdk.Context, key string, ptr interface{}) error {
	key = getGovStoreKey(key)
	return k.k.get(ctx, key, ptr)
}

// GetRaw exposes getRaw
func (k GovGetter) GetRaw(ctx sdk.Context, key string) []byte {
	key = getGovStoreKey(key)
	return k.k.getRaw(ctx, key)
}

// GetString is helper function for string params
func (k GovGetter) GetString(ctx sdk.Context, key string) (res string, err error) {
	key = getGovStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	err = k.k.cdc.UnmarshalBinary(bz, &res)
	return
}

// GetBool is helper function for bool params
func (k GovGetter) GetBool(ctx sdk.Context, key string) (res bool, err error) {
	key = getGovStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	err = k.k.cdc.UnmarshalBinary(bz, &res)
	return
}

// GetInt16 is helper function for int16 params
func (k GovGetter) GetInt16(ctx sdk.Context, key string) (res int16, err error) {
	key = getGovStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	err = k.k.cdc.UnmarshalBinary(bz, &res)
	return
}

// GetInt32 is helper function for int32 params
func (k GovGetter) GetInt32(ctx sdk.Context, key string) (res int32, err error) {
	key = getGovStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	err = k.k.cdc.UnmarshalBinary(bz, &res)
	return
}

// GetInt64 is helper function for int64 params
func (k GovGetter) GetInt64(ctx sdk.Context, key string) (res int64, err error) {
	key = getGovStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	err = k.k.cdc.UnmarshalBinary(bz, &res)
	return
}

// GetUint16 is helper function for uint16 params
func (k GovGetter) GetUint16(ctx sdk.Context, key string) (res uint16, err error) {
	key = getGovStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	err = k.k.cdc.UnmarshalBinary(bz, &res)
	return
}

// GetUint32 is helper function for uint32 params
func (k GovGetter) GetUint32(ctx sdk.Context, key string) (res uint32, err error) {
	key = getGovStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	err = k.k.cdc.UnmarshalBinary(bz, &res)
	return
}

// GetUint64 is helper function for uint64 params
func (k GovGetter) GetUint64(ctx sdk.Context, key string) (res uint64, err error) {
	key = getGovStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	err = k.k.cdc.UnmarshalBinary(bz, &res)
	return
}

// GetInt is helper function for sdk.Int params
func (k GovGetter) GetInt(ctx sdk.Context, key string) (res sdk.Int, err error) {
	key = getGovStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	err = k.k.cdc.UnmarshalBinary(bz, &res)
	return
}

// GetUint is helper function for sdk.Uint params
func (k GovGetter) GetUint(ctx sdk.Context, key string) (res sdk.Uint, err error) {
	key = getGovStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	err = k.k.cdc.UnmarshalBinary(bz, &res)
	return
}

// GetRat is helper function for rat params
func (k GovGetter) GetRat(ctx sdk.Context, key string) (res sdk.Rat, err error) {
	key = getGovStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	err = k.k.cdc.UnmarshalBinary(bz, &res)
	return
}

// GetStringWithDefault is helper function for string params with default value
func (k GovGetter) GetStringWithDefault(ctx sdk.Context, key string, def string) (res string) {
	key = getGovStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	if bz == nil {
		return def
	}
	k.k.cdc.MustUnmarshalBinary(bz, &res)
	return
}

// GetBoolWithDefault is helper function for bool params with default value
func (k GovGetter) GetBoolWithDefault(ctx sdk.Context, key string, def bool) (res bool) {
	key = getGovStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	if bz == nil {
		return def
	}
	k.k.cdc.MustUnmarshalBinary(bz, &res)
	return
}

// GetInt16WithDefault is helper function for int16 params with default value
func (k GovGetter) GetInt16WithDefault(ctx sdk.Context, key string, def int16) (res int16) {
	key = getGovStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	if bz == nil {
		return def
	}
	k.k.cdc.MustUnmarshalBinary(bz, &res)
	return
}

// GetInt32WithDefault is helper function for int32 params with default value
func (k GovGetter) GetInt32WithDefault(ctx sdk.Context, key string, def int32) (res int32) {
	key = getGovStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	if bz == nil {
		return def
	}
	k.k.cdc.MustUnmarshalBinary(bz, &res)
	return
}

// GetInt64WithDefault is helper function for int64 params with default value
func (k GovGetter) GetInt64WithDefault(ctx sdk.Context, key string, def int64) (res int64) {
	key = getGovStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	if bz == nil {
		return def
	}
	k.k.cdc.MustUnmarshalBinary(bz, &res)
	return
}

// GetUint16WithDefault is helper function for uint16 params with default value
func (k GovGetter) GetUint16WithDefault(ctx sdk.Context, key string, def uint16) (res uint16) {
	key = getGovStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	if bz == nil {
		return def
	}
	k.k.cdc.MustUnmarshalBinary(bz, &res)
	return
}

// GetUint32WithDefault is helper function for uint32 params with default value
func (k GovGetter) GetUint32WithDefault(ctx sdk.Context, key string, def uint32) (res uint32) {
	key = getGovStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	if bz == nil {
		return def
	}
	k.k.cdc.MustUnmarshalBinary(bz, &res)
	return
}

// GetUint64WithDefault is helper function for uint64 params with default value
func (k GovGetter) GetUint64WithDefault(ctx sdk.Context, key string, def uint64) (res uint64) {
	key = getGovStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	if bz == nil {
		return def
	}
	k.k.cdc.MustUnmarshalBinary(bz, &res)
	return
}

// GetIntWithDefault is helper function for sdk.Int params with default value
func (k GovGetter) GetIntWithDefault(ctx sdk.Context, key string, def sdk.Int) (res sdk.Int) {
	key = getGovStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	if bz == nil {
		return def
	}
	k.k.cdc.MustUnmarshalBinary(bz, &res)
	return
}

// GetUintWithDefault is helper function for sdk.Uint params with default value
func (k GovGetter) GetUintWithDefault(ctx sdk.Context, key string, def sdk.Uint) (res sdk.Uint) {
	key = getGovStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	if bz == nil {
		return def
	}
	k.k.cdc.MustUnmarshalBinary(bz, &res)
	return
}

// GetRatWithDefault is helper function for sdk.Rat params with default value
func (k GovGetter) GetRatWithDefault(ctx sdk.Context, key string, def sdk.Rat) (res sdk.Rat) {
	key = getGlobalStoreKey(key)
	store := ctx.KVStore(k.k.key)
	bz := store.Get([]byte(key))
	if bz == nil {
		return def
	}
	k.k.cdc.MustUnmarshalBinary(bz, &res)
	return
}
