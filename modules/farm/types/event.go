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

	AttributeValuePoolId  = "pool_id"
	AttributeValueCreator = "creator"
	AttributeValueAmount  = "amount"
	AttributeValueReward  = "reward"
)
