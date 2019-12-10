package types

// coinswap module event types
const (
	EventSwap            = "swap"
	EventAddLiquidity    = "add_liquidity"
	EventRemoveLiquidity = "remove_liquidity"

	AttributeValueCategory = ModuleName

	AttributeValueAmount     = "amount"
	AttributeValueSender     = "sender"
	AttributeValueRecipient  = "recipient"
	AttributeValueIsBuyOrder = "is_buy_order"
	AttributeValueTokenPair  = "token_pair"
)
