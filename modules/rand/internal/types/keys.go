package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the rand module
	ModuleName = "rand"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the asset module
	QuerierRoute = ModuleName

	// RouterKey is the msg router key for the asset module
	RouterKey = ModuleName
)

var (
	KeyDelimiter           = []byte(":")                 // key delimiter
	PrefixRand             = []byte("rands:")            // key prefix for the random number
	PrefixRandRequestQueue = []byte("randRequestQueue:") // key prefix for the random number request queue
)

// KeyRand returns the key for a random number by the specified request id
func KeyRand(reqID []byte) []byte {
	return append(PrefixRand, reqID...)
}

// KeyRandRequestQueue returns the key for the random number request queue by the given height and request id
func KeyRandRequestQueue(height int64, reqID []byte) []byte {
	prefix := append(PrefixRandRequestQueue, sdk.Uint64ToBigEndian(uint64(height))...)
	return append(append(prefix, KeyDelimiter...), reqID...)
}

// KeyRandRequestQueueSubspace returns the key prefix for iterating through all requests at the specified height
func KeyRandRequestQueueSubspace(height int64) []byte {
	return append(append(PrefixRandRequestQueue, sdk.Uint64ToBigEndian(uint64(height))...), KeyDelimiter...)
}
