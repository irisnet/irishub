package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "mt"

	// StoreKey is the default store key for MT
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the MT store.
	QuerierRoute = ModuleName

	// RouterKey is the message route for the MT module
	RouterKey = ModuleName

	// KeyNextDenomSequence is the key used to store the next denom sequence in the keeper
	KeyNextDenomSequence = "nextDenomSequence"

	// KeyNextMTSequence is the key used to store the next MT sequence in the keeper
	KeyNextMTSequence = "nextMTSequence"
)

var (
	PrefixDenom   = []byte{0x01}
	PrefixMT      = []byte{0x02}
	PrefixBalance = []byte{0x03}
	PrefixSupply  = []byte{0x04}

	delimiter = []byte("/")
)

// KeyDenom gets the storeKey by the denom
func KeyDenom(id string) []byte {
	key := append(PrefixDenom, delimiter...)
	return append(key, []byte(id)...)
}

// KeyMT gets the key of MT stored by an denom and MT
func KeyMT(denomID, mtID string) []byte {
	key := append(PrefixMT, delimiter...)
	if len(denomID) > 0 {
		key = append(key, []byte(denomID)...)
		key = append(key, delimiter...)
	}

	if len(denomID) > 0 && len(mtID) > 0 {
		key = append(key, []byte(mtID)...)
	}
	return key
}

// KeyBalance gets the key of a balance owned by an account address
func KeyBalance(address sdk.AccAddress, denomID, mtID string) []byte {
	key := append(PrefixBalance, delimiter...)
	if address != nil {
		key = append(key, []byte(address.String())...)
		key = append(key, delimiter...)
	}

	if address != nil && len(denomID) > 0 {
		key = append(key, []byte(denomID)...)
		key = append(key, delimiter...)
	}

	if address != nil && len(denomID) > 0 && len(mtID) > 0 {
		key = append(key, []byte(mtID)...)
	}
	return key
}

// KeySupply gets the key of supply of a denom or MT
func KeySupply(denomID, mtID string) []byte {
	key := append(PrefixSupply, delimiter...)

	if len(denomID) > 0 {
		key = append(key, []byte(denomID)...)
		key = append(key, delimiter...)
	}

	if len(denomID) > 0 && len(mtID) > 0 {
		key = append(key, []byte(mtID)...)
	}
	return key
}
