package types

import sdk "github.com/cosmos/cosmos-sdk/types"

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

	// RewardCollector is the root string for the reward distribution account address
	RewardCollector = "reward_collector"

	// EscrowCollector is the root string for the reward escrow account address
	EscrowCollector = "escrow_collector"

	// Prefix for farm pool_id
	PrefixFarmPool = "farm"
)

var (
	Delimiter         = []byte{0x00} // Separator for string key
	FarmPoolRuleKey   = []byte{0x02} // key for farm pool reward rule
	FarmerKey         = []byte{0x03} // key for farmer
	ActiveFarmPoolKey = []byte{0x04} // key for active farm pool
	FarmPoolSeq       = []byte{0x05}
	FarmPoolKey       = []byte{0x06} // key for farm pool
	EscrowInfoKey     = []byte{0x07}
	ParamsKey         = []byte{0x08}
)

func KeyFarmPool(poolId string) []byte {
	return append(FarmPoolKey, []byte(poolId)...)
}

func KeyRewardRule(poolId, reward string) []byte {
	key := append(FarmPoolRuleKey, []byte(poolId)...)
	return append(append(key, Delimiter...), []byte(reward)...)
}

func PrefixRewardRule(poolId string) []byte {
	key := append(FarmPoolRuleKey, []byte(poolId)...)
	return append(key, Delimiter...)
}

func KeyFarmInfo(address, poolId string) []byte {
	return append(append(FarmerKey, []byte(address)...), []byte(poolId)...)
}

func KeyFarmPoolSeq() []byte {
	return append(FarmPoolSeq, []byte("seq")...)
}

func PrefixFarmInfo(address string) []byte {
	return append(FarmerKey, []byte(address)...)
}

func KeyActiveFarmPool(height int64, poolId string) []byte {
	return append(
		append(ActiveFarmPoolKey, sdk.Uint64ToBigEndian(uint64(height))...),
		[]byte(poolId)...)
}

func PrefixActiveFarmPool(height int64) []byte {
	return append(ActiveFarmPoolKey, sdk.Uint64ToBigEndian(uint64(height))...)
}

func KeyEscrowInfo(proposalId uint64) []byte {
	return append(EscrowInfoKey, sdk.Uint64ToBigEndian(proposalId)...)
}
