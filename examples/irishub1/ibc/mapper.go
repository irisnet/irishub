package ibc

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	codec "github.com/cosmos/cosmos-sdk/codec"
)

// IBC Mapper
type Mapper struct {
	key       sdk.StoreKey
	cdc       *codec.Codec
	codespace sdk.CodespaceType
}

// XXX: The Mapper should not take a CoinKeeper. Rather have the CoinKeeper
// take an Mapper.
func NewMapper(cdc *codec.Codec, key sdk.StoreKey, codespace sdk.CodespaceType) Mapper {
	// XXX: How are these codecs supposed to work?
	return Mapper{
		key:       key,
		cdc:       cdc,
		codespace: codespace,
	}
}

// --------------------------
// Functions for accessing the underlying KVStore.

func MarshalBinaryLengthPrefixedPanic(cdc *codec.Codec, value interface{}) []byte {
	res, err := cdc.MarshalBinaryLengthPrefixed(value)
	if err != nil {
		panic(err)
	}
	return res
}

func unMarshalBinaryLengthPrefixedPanic(cdc *codec.Codec, bz []byte, ptr interface{}) {
	err := cdc.UnmarshalBinaryLengthPrefixed(bz, ptr)
	if err != nil {
		panic(err)
	}
}

// TODO add description
func (ibcm Mapper) GetIngressSequence(ctx sdk.Context, srcChain string) int64 {
	store := ctx.KVStore(ibcm.key)
	key := IngressSequenceKey(srcChain)

	bz := store.Get(key)
	if bz == nil {
		zero := MarshalBinaryLengthPrefixedPanic(ibcm.cdc, int64(0))
		store.Set(key, zero)
		return 0
	}

	var res int64
	unMarshalBinaryLengthPrefixedPanic(ibcm.cdc, bz, &res)
	return res
}

// TODO add description
func (ibcm Mapper) SetIngressSequence(ctx sdk.Context, srcChain string, sequence int64) {
	store := ctx.KVStore(ibcm.key)
	key := IngressSequenceKey(srcChain)

	bz := MarshalBinaryLengthPrefixedPanic(ibcm.cdc, sequence)
	store.Set(key, bz)
}

// Retrieves the index of the currently stored outgoing IBC packets.
func (ibcm Mapper) getEgressLength(store sdk.KVStore, destChain string) int64 {
	bz := store.Get(EgressLengthKey(destChain))
	if bz == nil {
		zero := MarshalBinaryLengthPrefixedPanic(ibcm.cdc, int64(0))
		store.Set(EgressLengthKey(destChain), zero)
		return 0
	}
	var res int64
	unMarshalBinaryLengthPrefixedPanic(ibcm.cdc, bz, &res)
	return res
}

// Stores an outgoing IBC packet under "egress/chain_id/index".
func EgressKey(destChain string, index int64) []byte {
	return []byte(fmt.Sprintf("egress/%s/%d", destChain, index))
}

// Stores the number of outgoing IBC packets under "egress/index".
func EgressLengthKey(destChain string) []byte {
	return []byte(fmt.Sprintf("egress/%s", destChain))
}

// Stores the sequence number of incoming IBC packet under "ingress/index".
func IngressSequenceKey(srcChain string) []byte {
	return []byte(fmt.Sprintf("ingress/%s", srcChain))
}


func (ibcm Mapper) Get(ctx sdk.Context) (string, bool) {
	store := ctx.KVStore(ibcm.key)
	bz := store.Get([]byte("ibcaddr"))
	if bz == nil {
		return " ", false
	}
	var Addr string
	ibcm.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &Addr)
	return Addr, true
}

func (ibcm Mapper) Set(ctx sdk.Context,Addr string) {
	store := ctx.KVStore(ibcm.key)
	bz := ibcm.cdc.MustMarshalBinaryLengthPrefixed(Addr)
	store.Set([]byte("ibcaddr"), bz)
}