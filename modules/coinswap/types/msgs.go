package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgSwapOrder{}
	_ sdk.Msg = &MsgAddLiquidity{}
	_ sdk.Msg = &MsgRemoveLiquidity{}
)

const (
	// FormatUniABSPrefix defines the prefix of liquidity token
	FormatUniABSPrefix = "swap"
	// FormatUniDenom defines the name of liquidity token
	FormatUniDenom = "swap%s"

	// TypeMsgAddLiquidity defines the type of MsgAddLiquidity
	TypeMsgAddLiquidity = "add_liquidity"
	// TypeMsgRemoveLiquidity defines the type of MsgRemoveLiquidity
	TypeMsgRemoveLiquidity = "remove_liquidity"
	// TypeMsgSwapOrder defines the type of MsgSwapOrder
	TypeMsgSwapOrder = "swap_order"
)

/* --------------------------------------------------------------------------- */
// MsgSwapOrder
/* --------------------------------------------------------------------------- */

// MsgSwapOrder - struct for swapping a coin
// Input and Output can either be exact or calculated.
// An exact coin has the senders desired buy or sell amount.
// A calculated coin has the desired denomination and bounded amount
// the sender is willing to buy or sell in this order.

// NewMsgSwapOrder creates a new MsgSwapOrder object.
func NewMsgSwapOrder(
	input Input,
	output Output,
	deadline int64,
	isBuyOrder bool,
) *MsgSwapOrder {
	return &MsgSwapOrder{
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
	if err := ValidateInput(msg.Input); err != nil {
		return err
	}

	if err := ValidateOutput(msg.Output); err != nil {
		return err
	}

	if msg.Input.Coin.Denom == msg.Output.Coin.Denom {
		return sdkerrors.Wrap(ErrEqualDenom, "invalid swap")
	}

	return ValidateDeadline(msg.Deadline)
}

// GetSignBytes implements Msg.
func (msg MsgSwapOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements Msg.
func (msg MsgSwapOrder) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Input.Address)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

/* --------------------------------------------------------------------------- */
// MsgAddLiquidity
/* --------------------------------------------------------------------------- */

// NewMsgAddLiquidity creates a new MsgAddLiquidity object.
func NewMsgAddLiquidity(
	maxToken sdk.Coin,
	exactStandardAmt sdk.Int,
	minLiquidity sdk.Int,
	deadline int64,
	sender string,
) *MsgAddLiquidity {
	return &MsgAddLiquidity{
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
	if err := ValidateMaxToken(msg.MaxToken); err != nil {
		return err
	}

	if err := ValidateExactStandardAmt(msg.ExactStandardAmt); err != nil {
		return err
	}

	if err := ValidateMinLiquidity(msg.MinLiquidity); err != nil {
		return err
	}

	if err := ValidateDeadline(msg.Deadline); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}

// GetSignBytes implements Msg.
func (msg MsgAddLiquidity) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements Msg.
func (msg MsgAddLiquidity) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

/* --------------------------------------------------------------------------- */
// MsgRemoveLiquidity
/* --------------------------------------------------------------------------- */

// NewMsgRemoveLiquidity creates a new MsgRemoveLiquidity object
func NewMsgRemoveLiquidity(
	minToken sdk.Int,
	withdrawLiquidity sdk.Coin,
	minStandardAmt sdk.Int,
	deadline int64,
	sender string,
) *MsgRemoveLiquidity {
	return &MsgRemoveLiquidity{
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
	if err := ValidateMinToken(msg.MinToken); err != nil {
		return err
	}

	if err := ValidateWithdrawLiquidity(msg.WithdrawLiquidity); err != nil {
		return err
	}

	if err := ValidateMinStandardAmt(msg.MinStandardAmt); err != nil {
		return err
	}

	if err := ValidateDeadline(msg.Deadline); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}

// GetSignBytes implements Msg.
func (msg MsgRemoveLiquidity) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements Msg.
func (msg MsgRemoveLiquidity) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}
