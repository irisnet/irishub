package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgSwapOrder{}
	_ sdk.Msg = &MsgAddLiquidity{}
	_ sdk.Msg = &MsgRemoveLiquidity{}
)

const (
	FormatUniABSPrefix = "uni:"   // format uni ABS Prefix
	FormatUniDenom     = "uni:%s" // format uni denom

	TypeMsgAddLiquidity    = "add_liquidity"    // type for MsgAddLiquidity
	TypeMsgRemoveLiquidity = "remove_liquidity" // type for MsgRemoveLiquidity
	TypeMsgSwapOrder       = "swap_order"       // type for MsgSwapOrder
)

/* --------------------------------------------------------------------------- */
// MsgSwapOrder
/* --------------------------------------------------------------------------- */

// MsgSwapOrder - struct for swapping a coin
// Input and Output can either be exact or calculated.
// An exact coin has the senders desired buy or sell amount.
// A calculated coin has the desired denomination and bounded amount
// the sender is willing to buy or sell in this order.

// Input defines the properties of order's input
type Input struct {
	Address sdk.AccAddress `json:"address" yaml:"address"` // address of input
	Coin    sdk.Coin       `json:"coin" yaml:"coin"`       // coins of input
}

// Output defines the properties of order's output
type Output struct {
	Address sdk.AccAddress `json:"address" yaml:"address"` // address of output
	Coin    sdk.Coin       `json:"coin" yaml:"coin"`       // coins of output
}

// MsgSwapOrder represents a msg for swap order
type MsgSwapOrder struct {
	Input      Input  `json:"input" yaml:"input"`               // the amount the sender is trading
	Output     Output `json:"output" yaml:"output"`             // the amount the sender is receiving
	Deadline   int64  `json:"deadline" yaml:"deadline"`         // deadline for the transaction to still be considered valid
	IsBuyOrder bool   `json:"is_buy_order" yaml:"is_buy_order"` // boolean indicating whether the order should be treated as a buy or sell
}

// NewMsgSwapOrder creates a new MsgSwapOrder object.
func NewMsgSwapOrder(
	input Input,
	output Output,
	deadline int64,
	isBuyOrder bool,
) MsgSwapOrder {
	return MsgSwapOrder{
		Input:      input,
		Output:     output,
		Deadline:   deadline,
		IsBuyOrder: isBuyOrder,
	}
}

// Route implements Msg.
func (msg MsgSwapOrder) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgSwapOrder) Type() string { return TypeMsgSwapOrder }

// ValidateBasic implements Msg.
func (msg MsgSwapOrder) ValidateBasic() error {
	if !(msg.Input.Coin.IsValid() && msg.Input.Coin.IsPositive()) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, fmt.Sprintf("input coin is invalid: %s", msg.Input.Coin.String()))
	}
	if strings.HasPrefix(msg.Input.Coin.Denom, FormatUniABSPrefix) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("unsupported input coin type: %s", msg.Input.Coin.String()))
	}
	if !(msg.Output.Coin.IsValid() && msg.Output.Coin.IsPositive()) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, fmt.Sprintf("output coin is invalid: %s", msg.Output.Coin.String()))
	}
	if strings.HasPrefix(msg.Output.Coin.Denom, FormatUniABSPrefix) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("unsupported output coin type: %s", msg.Output.Coin.String()))
	}
	if msg.Input.Coin.Denom == msg.Output.Coin.Denom {
		return ErrEqualDenom
	}
	if msg.Deadline <= 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("deadline %d must be greater than 0", msg.Deadline))
	}
	if msg.Input.Address.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "input address missing")
	}
	return nil
}

// GetSignBytes implements Msg.
func (msg MsgSwapOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements Msg.
func (msg MsgSwapOrder) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Input.Address}
}

/* --------------------------------------------------------------------------- */
// MsgAddLiquidity
/* --------------------------------------------------------------------------- */

// MsgAddLiquidity - struct for adding liquidity to a reserve pool
type MsgAddLiquidity struct {
	MaxToken         sdk.Coin       `json:"max_token" yaml:"max_token"`                   // coin to be deposited as liquidity with an upper bound for its amount
	ExactStandardAmt sdk.Int        `json:"exact_standard_amt" yaml:"exact_standard_amt"` // exact amount of native asset being add to the liquidity pool
	MinLiquidity     sdk.Int        `json:"min_liquidity" yaml:"min_liquidity"`           // lower bound UNI sender is willing to accept for deposited coins
	Deadline         int64          `json:"deadline" yaml:"deadline"`                     // deadline of tx
	Sender           sdk.AccAddress `json:"sender" yaml:"sender"`                         // msg sender
}

// NewMsgAddLiquidity creates a new MsgAddLiquidity object.
func NewMsgAddLiquidity(
	maxToken sdk.Coin,
	exactStandardAmt sdk.Int,
	minLiquidity sdk.Int,
	deadline int64,
	sender sdk.AccAddress,
) MsgAddLiquidity {
	return MsgAddLiquidity{
		MaxToken:         maxToken,
		ExactStandardAmt: exactStandardAmt,
		MinLiquidity:     minLiquidity,
		Deadline:         deadline,
		Sender:           sender,
	}
}

// Route implements Msg.
func (msg MsgAddLiquidity) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgAddLiquidity) Type() string { return TypeMsgAddLiquidity }

// ValidateBasic implements Msg.
func (msg MsgAddLiquidity) ValidateBasic() error {
	if !(msg.MaxToken.IsValid() && msg.MaxToken.IsPositive()) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, fmt.Sprintf("max token is invalid: %s", msg.MaxToken.String()))
	}
	if msg.MaxToken.Denom == StandardDenom {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("max token must not be standard token: %s", StandardDenom))
	}
	if strings.HasPrefix(msg.MaxToken.Denom, FormatUniABSPrefix) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "max token must be non-liquidity token")
	}
	if !msg.ExactStandardAmt.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "standard token amount must be positive")
	}
	if msg.MinLiquidity.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "minimum liquidity can not be negative")
	}
	if msg.Deadline <= 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("deadline %d must be greater than 0", msg.Deadline))
	}
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "sender address missing")
	}
	return nil
}

// GetSignBytes implements Msg.
func (msg MsgAddLiquidity) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements Msg.
func (msg MsgAddLiquidity) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

/* --------------------------------------------------------------------------- */
// MsgRemoveLiquidity
/* --------------------------------------------------------------------------- */

// MsgRemoveLiquidity - struct for removing liquidity from a reserve pool
type MsgRemoveLiquidity struct {
	MinToken          sdk.Int        `json:"min_token" yaml:"min_token"`                   // coin to be withdrawn with a lower bound for its amount
	WithdrawLiquidity sdk.Coin       `json:"withdraw_liquidity" yaml:"withdraw_liquidity"` // amount of UNI to be burned to withdraw liquidity from a reserve pool
	MinStandardAmt    sdk.Int        `json:"min_standard_amt" yaml:"min_standard_amt"`     // minimum amount of the native asset the sender is willing to accept
	Deadline          int64          `json:"deadline" yaml:"deadline"`                     // deadline of tx
	Sender            sdk.AccAddress `json:"sender" yaml:"sender"`                         // msg sender
}

// NewMsgRemoveLiquidity creates a new MsgRemoveLiquidity object
func NewMsgRemoveLiquidity(
	minToken sdk.Int,
	withdrawLiquidity sdk.Coin,
	minStandardAmt sdk.Int,
	deadline int64,
	sender sdk.AccAddress,
) MsgRemoveLiquidity {
	return MsgRemoveLiquidity{
		MinToken:          minToken,
		WithdrawLiquidity: withdrawLiquidity,
		MinStandardAmt:    minStandardAmt,
		Deadline:          deadline,
		Sender:            sender,
	}
}

// Route implements Msg.
func (msg MsgRemoveLiquidity) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgRemoveLiquidity) Type() string { return TypeMsgRemoveLiquidity }

// ValidateBasic implements Msg.
func (msg MsgRemoveLiquidity) ValidateBasic() error {
	if msg.MinToken.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "minimum token amount can not be negative")
	}
	if !msg.WithdrawLiquidity.IsValid() || !msg.WithdrawLiquidity.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, fmt.Sprintf("withdraw liquidity %s is not valid", msg.WithdrawLiquidity.String()))
	}
	if err := CheckUniDenom(msg.WithdrawLiquidity.Denom); err != nil {
		return err
	}
	if msg.MinStandardAmt.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("minimum standard token amount %s can not be negative", msg.MinStandardAmt.String()))
	}
	if msg.Deadline <= 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("deadline %d must be greater than 0", msg.Deadline))
	}
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "sender address missing")
	}
	return nil
}

// GetSignBytes implements Msg.
func (msg MsgRemoveLiquidity) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements Msg.
func (msg MsgRemoveLiquidity) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
