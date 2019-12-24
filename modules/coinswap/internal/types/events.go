// nolint
package types

// coinswap module event types
const (
	EventTypeSwap            = "swap"
	EventTypeAddLiquidity    = "add_liquidity"
	EventTypeRemoveLiquidity = "remove_liquidity"

	AttributeValueAmount     = "amount"
	AttributeValueSender     = "sender"
	AttributeValueRecipient  = "recipient"
	AttributeValueIsBuyOrder = "is_buy_order"
	AttributeValueTokenPair  = "token_pair"

	AttributeValueCategory = ModuleName
)
