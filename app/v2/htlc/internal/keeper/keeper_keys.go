package keeper

import (
	"fmt"
)

var (
	PrefixHTLC            = []byte("htlcs:")           // key prefix for HTLC
	PrefixHTLCExpireQueue = []byte("htlcExpireQueue:") // key prefix for the HTLC expiration queue
)

// KeyHTLC returns the key for a HTLC by the specified secret hash
func KeyHTLC(secretHashLock []byte) []byte {
	return append(PrefixHTLC, secretHashLock...)
}

// KeyHTLCExpireQueue returns the key prefix for HTLC expiration queue
func KeyHTLCExpireQueue(expireHeight uint64, secretHashLock []byte) []byte {
	return append([]byte(fmt.Sprintf("PrefixHTLCExpireQueue:%d:", expireHeight)), secretHashLock...)
}
