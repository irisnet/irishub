package keeper

import (
	"fmt"
)

var (
	PrefixHTLC            = []byte("htlcs:")           // key prefix for HTLC
	PrefixHTLCState       = []byte("htlcState:")       // key prefix for the HTLC state
	PrefixHTLCToken       = []byte("htlcTokens:")      // key prefix for the HTLC-locked token
	PrefixHTLCExpireQueue = []byte("htlcExpireQueue:") // key prefix for the HTLC expiration queue
)

// KeyHTLC returns the key for a HTLC by the specified secret hash
func KeyHTLC(secretHashLock []byte) []byte {
	return append(PrefixHTLC, secretHashLock...)
}

// KeyHTLCState returns the key for the HTLC state by the given secret hash lock
func KeyHTLCState(secretHashLock []byte) []byte {
	return append(PrefixHTLCState, secretHashLock...)
}

// KeyHTLCTokens returns the key for the HTLC-locked token by the given token id
func KeyHTLCTokens(tokenId string) []byte {
	return append(PrefixHTLCToken, []byte(tokenId)...)
}

// KeyHTLCExpireQueue returns the key prefix for HTLC expiration queue
func KeyHTLCExpireQueue(expireHeight uint64, secretHashLock []byte) []byte {
	return append([]byte(fmt.Sprintf("PrefixHTLCExpireQueue:%d:", expireHeight)), secretHashLock...)
}
