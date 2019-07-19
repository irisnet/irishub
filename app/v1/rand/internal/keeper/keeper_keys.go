package keeper

import (
	"fmt"
)

var (
	PrefixRand             = []byte("rands:")            // key prefix for the random number
	PrefixRandRequestQueue = []byte("randRequestQueue:") // key prefix for the random number request queue
)

// KeyRand returns the key for a random number by the specified request id
func KeyRand(reqID string) []byte {
	return []byte(fmt.Sprintf("rands:%s", reqID))
}

// KeyRandRequestQueue returns the key for the random number request queue by the given height and request id
func KeyRandRequestQueue(height int64, reqID string) []byte {
	return []byte(fmt.Sprintf("randRequestQueue:%d:%s", height, reqID))
}

// KeyRandRequestQueueSubspace returns the key prefix for iterating through all requests at the specified height
func KeyRandRequestQueueSubspace(height int64) []byte {
	return []byte(fmt.Sprintf("randRequestQueue:%d:", height))
}
