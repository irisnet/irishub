package types

import "fmt"

const (
	// ModuleName is the name of the module.
	ModuleName = "coinswap"

	// RouterKey is the message route for the coinswap module.
	RouterKey = ModuleName

	// StoreKey is the default store key for the coinswap module.
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the coinswap module.
	QuerierRoute = StoreKey

	// KeyNextPoolSequence is the key used to store the next pool sequence in
	// the keeper.
	KeyNextPoolSequence = "nextPoolSequence"

	// KeyPool is the key used to store the pool information  in
	// the keeper.
	KeyPool = "pool"

	// KeyPoolLptDenom is the key used to store the pool information  in
	// the keeper.
	KeyPoolLptDenom = "lptDenom"
)

// GetPoolKey return the stored pool key for the given pooId.
func GetPoolKey(pooId string) []byte {
	return []byte(fmt.Sprintf("%s/%s", KeyPool, pooId))
}

// GetLptDenomKey return the stored pool key for the given liquidity pool token denom.
func GetLptDenomKey(lptDenom string) []byte {
	return []byte(fmt.Sprintf("%s/%s", KeyPoolLptDenom, lptDenom))
}
