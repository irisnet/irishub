package profiling

import (
	sdk "github.com/irisnet/irishub/types"
)

var (
	profilerKey = []byte{0x00}
)

func GetProfilerKey(addr sdk.AccAddress) []byte {
	return append(profilerKey, addr.Bytes()...)
}

// Key for getting all profilers from the store
func GetProfilersSubspaceKey() []byte {
	return profilerKey
}
