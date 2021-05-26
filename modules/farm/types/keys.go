package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// nolint
const (
	// module name
	ModuleName = "farm"

	// StoreKey is the default store key for farm
	StoreKey = ModuleName

	// RouterKey is the message route for farm
	RouterKey = ModuleName

	// QuerierRoute is the querier route for the farm store.
	QuerierRoute = StoreKey

	// Query endpoints supported by the farm querier
	QueryRecord = "farm"
)

var (
	FarmPoolKey       = []byte{0x01} // key for farm pool
	FarmPoolRuleKey   = []byte{0x02} // key for farm pool reward rule
	FarmerKey         = []byte{0x03} // key for farmer
	ActiveFarmPoolKey = []byte{0x04} // key for active farm pool
)

func KeyFarmPool(poolName string) []byte {
	return append(FarmPoolKey, []byte(poolName)...)
}

func KeyRewardRule(poolName, reward string) []byte {
	return append(append(FarmPoolRuleKey, []byte(poolName)...), []byte(reward)...)
}

func PrefixRewardRule(poolName string) []byte {
	return append(FarmPoolRuleKey, []byte(poolName)...)
}

func KeyFarmInfo(address, poolName string) []byte {
	return append(append(FarmerKey, []byte(address)...), []byte(poolName)...)
}

func PrefixFarmInfo(address string) []byte {
	return append(FarmerKey, []byte(address)...)
}

func KeyActiveFarmPool(expiredHeight uint64, poolName string) []byte {
	return append(append(ActiveFarmPoolKey, sdk.Uint64ToBigEndian(expiredHeight)...), []byte(poolName)...)
}

func PrefixActiveFarmPool(expiredHeight uint64) []byte {
	return append(ActiveFarmPoolKey, sdk.Uint64ToBigEndian(expiredHeight)...)
}
