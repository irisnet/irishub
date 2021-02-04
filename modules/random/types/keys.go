package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the random module
	ModuleName = "random"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the asset module
	QuerierRoute = ModuleName

	// RouterKey is the msg router key for the asset module
	RouterKey = ModuleName
)

var (
	RandomKey              = []byte{0x01} // key prefix for the random number
	RandomRequestQueueKey  = []byte{0x02} // key prefix for the random number request queue
	OracleRandomRequestKey = []byte{0x03} // key prefix for the oracle request
)

// KeyRandom returns the key for a random number by the specified request id
func KeyRandom(reqID []byte) []byte {
	return append(RandomKey, reqID...)
}

// KeyRandomRequestQueue returns the key for the random number request queue by the given height and request id
func KeyRandomRequestQueue(height int64, reqID []byte) []byte {
	return append(append(RandomRequestQueueKey, sdk.Uint64ToBigEndian(uint64(height))...), reqID...)
}

// KeyRandRequestQueueSubspace returns the key prefix for iterating through all requests at the specified height
func KeyRandomRequestQueueSubspace(height int64) []byte {
	return append(RandomRequestQueueKey, sdk.Uint64ToBigEndian(uint64(height))...)
}

// KeyOracleRandRequest returns the key for an OracleRandRequest by the specified requestContextID
func KeyOracleRandomRequest(requestContextID []byte) []byte {
	return append(OracleRandomRequestKey, requestContextID...)
}
