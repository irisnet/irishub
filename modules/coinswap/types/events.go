package types

// coinswap module event types
const (
	EventTypeSwap                      = "swap"
	EventTypeAddLiquidity              = "add_liquidity"
	EventTypeRemoveLiquidity           = "remove_liquidity"
	EventTypeAddUnilateralLiquidity    = "add_unilateral_liquidity"
	EventTypeRemoveUnilateralLiquidity = "remove_unilateral_liquidity"

	AttributeValueCategory = ModuleName

	AttributeValueAmount          = "amount"
	AttributeValueSender          = "sender"
	AttributeValueRecipient       = "recipient"
	AttributeValueIsBuyOrder      = "is_buy_order"
	AttributeValueTokenPair       = "token_pair"
	AttributeValueTokenUnilateral = "token_unilateral"
	AttributeValueLptDenom        = "lpt_denom"
)
