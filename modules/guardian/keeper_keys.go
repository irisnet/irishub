package guardian

import (
	sdk "github.com/irisnet/irishub/types"
)

var (
	profilerKey = []byte{0x00}
	trusteeKey  = []byte{0x01}
)

func GetProfilerKey(addr sdk.AccAddress) []byte {
	return append(profilerKey, addr.Bytes()...)
}

func GetTrusteeKey(addr sdk.AccAddress) []byte {
	return append(trusteeKey, addr.Bytes()...)
}

// Key for getting all profilers from the store
func GetProfilersSubspaceKey() []byte {
	return profilerKey
}

// Key for getting all profilers from the store
func GetTrusteesSubspaceKey() []byte {
	return trusteeKey
}
