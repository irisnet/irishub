package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the HTLC module
	ModuleName = "htlc"

	// StoreKey is the string store representation
	StoreKey string = ModuleName

	// QuerierRoute is the querier route for the HTLC module
	QuerierRoute string = ModuleName

	// RouterKey is the msg router key for the HTLC module
	RouterKey string = ModuleName
)

var (
	// Keys for store prefixes
	HTLCKey             = []byte{0x01} // prefix for HTLC
	HTLCExpiredQueueKey = []byte{0x02} // prefix for the HTLC expiration queue
)

// GetHTLCKey returns the key for the HTLC with the specified hash lock
// VALUE: htlc/HTLC
func GetHTLCKey(hashLock []byte) []byte {
	return append(HTLCKey, hashLock...)
}

// GetHTLCExpiredQueueKey returns the key for the HTLC expiration queue by the specified height and hash lock
// VALUE: []byte{}
func GetHTLCExpiredQueueKey(expirationHeight uint64, hashLock []byte) []byte {
	return append(append(HTLCExpiredQueueKey, sdk.Uint64ToBigEndian(expirationHeight)...), hashLock...)
}

// GetHTLCExpiredQueueSubspace returns the key prefix for the HTLC expiration queue by the given height
func GetHTLCExpiredQueueSubspace(expirationHeight uint64) []byte {
	return append(HTLCExpiredQueueKey, sdk.Uint64ToBigEndian(expirationHeight)...)
}
