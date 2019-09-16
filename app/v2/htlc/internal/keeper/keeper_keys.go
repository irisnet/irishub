package keeper

import (
	sdk "github.com/irisnet/irishub/types"
)

var (
	KeyDelimiter          = []byte(":")                // key separator
	PrefixHTLC            = []byte("htlcs:")           // key prefix for HTLC
	PrefixHTLCExpireQueue = []byte("htlcExpireQueue:") // key prefix for the HTLC expiration queue
)

// KeyHTLC returns the key for an HTLC by the specified hash lock
func KeyHTLC(hashLock []byte) []byte {
	return append(PrefixHTLC, hashLock...)
}

// KeyHTLCExpireQueue returns the key for HTLC expiration queue by the specified height and hash lock
func KeyHTLCExpireQueue(expireHeight uint64, hashLock []byte) []byte {
	prefix := append(PrefixHTLCExpireQueue, sdk.Uint64ToBigEndian(expireHeight)...)
	return append(append(prefix, KeyDelimiter...), hashLock...)
}

// KeyHTLCExpireQueueSubspace returns the key prefix for HTLC expiration queue by the given height
func KeyHTLCExpireQueueSubspace(expireHeight uint64) []byte {
	return append(append(PrefixHTLCExpireQueue, sdk.Uint64ToBigEndian(expireHeight)...), KeyDelimiter...)
}
