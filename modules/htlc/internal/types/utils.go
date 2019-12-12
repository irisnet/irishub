package types

import (
	"crypto/sha256"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SHA256 wraps sha256.Sum256 with result converted to slice
func SHA256(data []byte) []byte {
	sum := sha256.Sum256(data)
	return sum[:]
}

// GetHashLock calculates the hash lock from the given secret and timestamp
func GetHashLock(secret HTLCSecret, timestamp uint64) HTLCHashLock {
	if timestamp > 0 {
		return HTLCHashLock(SHA256(append(secret, sdk.Uint64ToBigEndian(timestamp)...)))
	}
	return HTLCHashLock(SHA256(secret))
}
