// nolint
package types

// farm module event types
const (
	EventTypeCreatePool   = "create_pool"
	EventTypeDestroyPool  = "destroy_pool"
	EventTypeAppendReward = "append_reward"
	EventTypeStake        = "stake"
	EventTypeUnstake      = "unstake"
	EventTypeHarvest      = "harvest"

	AttributeValueCategory = ModuleName

	AttributeValuePoolName = "pool_name"
	AttributeValueCreator  = "creator"
	AttributeValueAmount   = "amount"
	AttributeValueReward   = "reward"
)
