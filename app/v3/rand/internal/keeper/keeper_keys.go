package keeper

import (
	"fmt"
)

var (
	KeyDelimiter           = []byte(":")                 // key delimiter
	PrefixRand             = []byte("rands:")            // key prefix for the random number
	PrefixSeed             = []byte("seeds:")            // key prefix for the oracle seed
	PrefixRandRequestQueue = []byte("randRequestQueue:") // key prefix for the random number request queue
)

// KeyRand returns the key for a random number by the specified request id
func KeyRand(reqID []byte) []byte {
	return append(PrefixRand, reqID...)
}

// KeySeed returns the key for a oracle seed by the specified request id
func KeySeed(reqID []byte) []byte {
	return append(PrefixSeed, reqID...)
}

// KeyRandRequestQueue returns the key for the random number request queue by the given height and request id
func KeyRandRequestQueue(height int64, reqID []byte) []byte {
	return append([]byte(fmt.Sprintf("randRequestQueue:%d:", height)), reqID...)
}

// KeyRandRequestQueueSubspace returns the key prefix for iterating through all requests at the specified height
func KeyRandRequestQueueSubspace(height int64) []byte {
	return []byte(fmt.Sprintf("randRequestQueue:%d:", height))
}
