package types

import (
	sdk "github.com/irisnet/irishub/types"
	"strings"
)

var (
	_ sdk.Msg = MsgSwapOrder{}
	_ sdk.Msg = MsgAddLiquidity{}
	_ sdk.Msg = MsgRemoveLiquidity{}
)

const (
	FormatUniABSPrefix = sdk.FormatUniABSPrefix
	FormatUniId        = FormatUniABSPrefix + "%s"
)

/* --------------------------------------------------------------------------- */
// MsgSwapOrder
/* --------------------------------------------------------------------------- */

// MsgSwapOrder - struct for swapping a coin
// Input and Output can either be exact or calculated.
// An exact coin has the senders desired buy or sell amount.
// A calculated coin has the desired denomination and bounded amount
// the sender is willing to buy or sell in this order.
type Input struct {
	Address sdk.AccAddress `json:"address"`
	Coin    sdk.Coin       `json:"coin"`
}

type Output struct {
	Address sdk.AccAddress `json:"address"`
	Coin    sdk.Coin       `json:"coin"`
}

type MsgSwapOrder struct {
	Input      Input  `json:"input"`        // the amount the sender is trading
	Output     Output `json:"output"`       // the amount the sender is receiving
	Deadline   int64  `json:"deadline"`     // deadline for the transaction to still be considered valid
	IsBuyOrder bool   `json:"is_buy_order"` // boolean indicating whether the order should be treated as a buy or sell
}

// NewMsgSwapOrder creates a new MsgSwapOrder object.
func NewMsgSwapOrder(
	input Input, output Output, deadline int64, isBuyOrder bool,
) MsgSwapOrder {

	return MsgSwapOrder{
		Input:      input,
		Output:     output,
		Deadline:   deadline,
		IsBuyOrder: isBuyOrder,
	}
}

// Route Implements Msg.
func (msg MsgSwapOrder) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgSwapOrder) Type() string { return MsgTypeSwapOrder }

// ValidateBasic Implements Msg.
func (msg MsgSwapOrder) ValidateBasic() sdk.Error {
	if !(msg.Input.Coin.IsValid() && msg.Input.Coin.IsPositive()) {
		return sdk.ErrInvalidCoins("input coin is invalid: " + msg.Input.Coin.String())
	}
	if strings.HasPrefix(msg.Input.Coin.Denom, FormatUniABSPrefix) {
		return sdk.ErrInvalidCoins("unsupported input coin type: " + msg.Input.Coin.String())
	}
	if !(msg.Output.Coin.IsValid() && msg.Output.Coin.IsPositive()) {
		return sdk.ErrInvalidCoins("output coin is invalid: " + msg.Output.Coin.String())
	}
	if strings.HasPrefix(msg.Output.Coin.Denom, FormatUniABSPrefix) {
		return sdk.ErrInvalidCoins("unsupported output coin type: " + msg.Output.Coin.String())
	}
	if msg.Input.Coin.Denom == msg.Output.Coin.Denom {
		return ErrEqualDenom("")
	}
	if msg.Deadline <= 0 {
		return ErrInvalidDeadline("deadline for MsgSwapOrder not initialized")
	}
	if msg.Input.Address.Empty() {
		return sdk.ErrInvalidAddress("invalid input address")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSwapOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgSwapOrder) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Input.Address}
}

/* --------------------------------------------------------------------------- */
// MsgAddLiquidity
/* --------------------------------------------------------------------------- */

// MsgAddLiquidity - struct for adding liquidity to a reserve pool
type MsgAddLiquidity struct {
	MaxToken     sdk.Coin       `json:"max_token"`      // coin to be deposited as liquidity with an upper bound for its amount
	ExactIrisAmt sdk.Int        `json:"exact_iris_amt"` // exact amount of native asset being add to the liquidity pool
	MinLiquidity sdk.Int        `json:"min_liquidity"`  // lower bound UNI sender is willing to accept for deposited coins
	Deadline     int64          `json:"deadline"`
	Sender       sdk.AccAddress `json:"sender"`
}

// NewMsgAddLiquidity creates a new MsgAddLiquidity object.
func NewMsgAddLiquidity(
	maxToken sdk.Coin, exactIrisAmt, minLiquidity sdk.Int,
	deadline int64, sender sdk.AccAddress,
) MsgAddLiquidity {

	return MsgAddLiquidity{
		MaxToken:     maxToken,
		ExactIrisAmt: exactIrisAmt,
		MinLiquidity: minLiquidity,
		Deadline:     deadline,
		Sender:       sender,
	}
}

// Route Implements Msg.
func (msg MsgAddLiquidity) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgAddLiquidity) Type() string { return MsgTypeAddLiquidity }

// ValidateBasic Implements Msg.
func (msg MsgAddLiquidity) ValidateBasic() sdk.Error {
	if !(msg.MaxToken.IsValid() && msg.MaxToken.IsPositive()) {
		return sdk.ErrInvalidCoins("max token is invalid: " + msg.MaxToken.String())
	}
	if msg.MaxToken.Denom == sdk.IrisAtto {
		return sdk.ErrInvalidCoins("max token must be non-iris token")
	}
	if strings.HasPrefix(msg.MaxToken.Denom, FormatUniABSPrefix) {
		return sdk.ErrInvalidCoins("max token must be non-liquidity token")
	}
	if msg.ExactIrisAmt.IsNil() || !msg.ExactIrisAmt.IsPositive() {
		return ErrNotPositive("iris amount must be positive")
	}
	if msg.MinLiquidity.IsNil() || msg.MinLiquidity.IsNegative() {
		return ErrNotPositive("minimum liquidity can not be negative")
	}
	if msg.Deadline <= 0 {
		return ErrInvalidDeadline("deadline for MsgAddLiquidity not initialized")
	}
	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress("invalid sender address")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgAddLiquidity) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgAddLiquidity) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

/* --------------------------------------------------------------------------- */
// MsgRemoveLiquidity
/* --------------------------------------------------------------------------- */

// MsgRemoveLiquidity - struct for removing liquidity from a reserve pool
type MsgRemoveLiquidity struct {
	MinToken          sdk.Int        `json:"min_token"`          // coin to be withdrawn with a lower bound for its amount
	WithdrawLiquidity sdk.Coin       `json:"withdraw_liquidity"` // amount of UNI to be burned to withdraw liquidity from a reserve pool
	MinIrisAmt        sdk.Int        `json:"min_iris_amt"`       // minimum amount of the native asset the sender is willing to accept
	Deadline          int64          `json:"deadline"`
	Sender            sdk.AccAddress `json:"sender"`
}

// NewMsgRemoveLiquidity creates a new MsgRemoveLiquidity object
func NewMsgRemoveLiquidity(
	minToken sdk.Int, withdrawLiquidity sdk.Coin, minIrisAmt sdk.Int,
	deadline int64, sender sdk.AccAddress,
) MsgRemoveLiquidity {

	return MsgRemoveLiquidity{
		MinToken:          minToken,
		WithdrawLiquidity: withdrawLiquidity,
		MinIrisAmt:        minIrisAmt,
		Deadline:          deadline,
		Sender:            sender,
	}
}

// Route Implements Msg.
func (msg MsgRemoveLiquidity) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgRemoveLiquidity) Type() string { return MsgTypeRemoveLiquidity }

// ValidateBasic Implements Msg.
func (msg MsgRemoveLiquidity) ValidateBasic() sdk.Error {
	if msg.MinToken.IsNil() || msg.MinToken.IsNegative() {
		return sdk.ErrInvalidCoins("minimum token amount can not be negative")
	}
	if !msg.WithdrawLiquidity.IsValid() || !msg.WithdrawLiquidity.IsPositive() {
		return ErrNotPositive("withdraw liquidity is not valid: " + msg.WithdrawLiquidity.String())
	}
	if err := CheckUniDenom(msg.WithdrawLiquidity.Denom); err != nil {
		return err
	}
	if msg.MinIrisAmt.IsNil() || msg.MinIrisAmt.IsNegative() {
		return ErrNotPositive("minimum iris amount can not be negative")
	}
	if msg.Deadline <= 0 {
		return ErrInvalidDeadline("deadline for MsgRemoveLiquidity not initialized")
	}
	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress("invalid sender address")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgRemoveLiquidity) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgRemoveLiquidity) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
