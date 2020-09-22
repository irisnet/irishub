package keeper

import (
	"encoding/binary"
	"fmt"

	gogotypes "github.com/gogo/protobuf/types"

	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/record/types"
)

// Keeper of the record store
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      codec.Marshaler
}

// NewKeeper returns a record keeper
func NewKeeper(cdc codec.Marshaler, key sdk.StoreKey) Keeper {
	keeper := Keeper{
		storeKey: key,
		cdc:      cdc,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("irismod/%s", types.ModuleName))
}

// AddRecord add a record
func (k Keeper) AddRecord(ctx sdk.Context, record types.Record) []byte {
	store := ctx.KVStore(k.storeKey)
	recordBz := k.cdc.MustMarshalBinaryBare(&record)
	intraTxCounter := k.GetIntraTxCounter(ctx)

	bz := make([]byte, 4+len(recordBz))
	copy(bz[:len(recordBz)], recordBz[:])
	binary.BigEndian.PutUint32(bz[len(recordBz):], intraTxCounter)

	recordID := getRecordID(bz)
	store.Set(types.GetRecordKey(recordID), recordBz)

	// update intraTxCounter + 1
	k.SetIntraTxCounter(ctx, intraTxCounter+1)
	return recordID
}

// GetRecord retrieves the record by specified recordID
func (k Keeper) GetRecord(ctx sdk.Context, recordID []byte) (record types.Record, found bool) {
	store := ctx.KVStore(k.storeKey)
	if bz := store.Get(types.GetRecordKey(recordID)); bz != nil {
		k.cdc.MustUnmarshalBinaryBare(bz, &record)
		return record, true
	}
	return record, false
}

// RecordsIterator gets all records
func (k Keeper) RecordsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.RecordKey)
}

// GetIntraTxCounter gets the current in-block request operation counter
func (k Keeper) GetIntraTxCounter(ctx sdk.Context) uint32 {
	store := ctx.KVStore(k.storeKey)

	b := store.Get(types.IntraTxCounterKey)
	if b == nil {
		return 0
	}

	var counter gogotypes.UInt32Value
	k.cdc.MustUnmarshalBinaryBare(b, &counter)

	return counter.Value
}

// SetIntraTxCounter sets the current in-block request counter
func (k Keeper) SetIntraTxCounter(ctx sdk.Context, counter uint32) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryBare(&gogotypes.UInt32Value{Value: counter})
	store.Set(types.IntraTxCounterKey, bz)
}

func getRecordID(bz []byte) []byte {
	return tmhash.Sum(bz)
}
