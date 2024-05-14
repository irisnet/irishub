package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgSwapOrder{}
	_ sdk.Msg = &MsgAddLiquidity{}
	_ sdk.Msg = &MsgAddUnilateralLiquidity{}
	_ sdk.Msg = &MsgRemoveLiquidity{}
	_ sdk.Msg = &MsgRemoveUnilateralLiquidity{}
	_ sdk.Msg = &MsgUpdateParams{}
)

const (
	// LptTokenPrefix defines the prefix of liquidity token
	LptTokenPrefix = "lpt"
	// LptTokenFormat defines the name of liquidity token
	LptTokenFormat = "lpt-%d"

	// TypeMsgAddLiquidity defines the type of MsgAddLiquidity
	TypeMsgAddLiquidity = "add_liquidity"
	// TypeMsgAddUnilateralLiquidity defines the type of MsgAddUnilateralLiquidity
	TypeMsgAddUnilateralLiquidity = "add_unilateral_liquidity"
	// TypeMsgRemoveLiquidity defines the type of MsgRemoveLiquidity
	TypeMsgRemoveLiquidity = "remove_liquidity"
	// TypeMsgRemoveUnilateralLiquidity defines the type of MsgRemoveUnilateralLiquidity
	TypeMsgRemoveUnilateralLiquidity = "remove_unilateral_liquidity"
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
		return errorsmod.Wrap(ErrEqualDenom, "invalid swap")
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
	exactStandardAmt sdkmath.Int,
	minLiquidity sdkmath.Int,
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

	if err := ValidateToken(msg.MaxToken); err != nil {
		return err
	}

	if err := ValidateExactStandardAmt(msg.ExactStandardAmt); err != nil {
		return err
	}

	if err := ValidateLiquidity(msg.MinLiquidity); err != nil {
		return err
	}

	if err := ValidateDeadline(msg.Deadline); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
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
// MsgAddUnilateralLiquidity
/* --------------------------------------------------------------------------- */
func NewMsgAddUnilateralLiquidity(
	counterpartyDenom string,
	exactToken sdk.Coin,
	minLiquidity sdkmath.Int,
	deadline int64,
	sender string,
) *MsgAddUnilateralLiquidity {
	return &MsgAddUnilateralLiquidity{
		CounterpartyDenom: counterpartyDenom,
		ExactToken:        exactToken,
		MinLiquidity:      minLiquidity,
		Deadline:          deadline,
		Sender:            sender,
	}
}

func (m MsgAddUnilateralLiquidity) Route() string { return RouterKey }

func (m MsgAddUnilateralLiquidity) Type() string { return TypeMsgAddUnilateralLiquidity }

func (m MsgAddUnilateralLiquidity) ValidateBasic() error {
	if err := ValidateCounterpartyDenom(m.CounterpartyDenom); err != nil {
		return err
	}

	if err := ValidateToken(m.ExactToken); err != nil {
		return err
	}

	if err := ValidateLiquidity(m.MinLiquidity); err != nil {
		return err
	}

	if err := ValidateDeadline(m.Deadline); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}

func (m MsgAddUnilateralLiquidity) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgAddUnilateralLiquidity) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.Sender)
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
	minToken sdkmath.Int,
	withdrawLiquidity sdk.Coin,
	minStandardAmt sdkmath.Int,
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
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
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

/* --------------------------------------------------------------------------- */
// MsgRemoveUnilateralLiquidity
/* --------------------------------------------------------------------------- */
func NewMsgRemoveUnilateralLiquidity(
	couterpartyDenom string,
	minToken sdk.Coin,
	exactLiquidity sdkmath.Int,
	deadline int64,
	sender string,
) *MsgRemoveUnilateralLiquidity {
	return &MsgRemoveUnilateralLiquidity{
		CounterpartyDenom: couterpartyDenom,
		MinToken:          minToken,
		ExactLiquidity:    exactLiquidity,
		Deadline:          deadline,
		Sender:            sender,
	}
}

func (m MsgRemoveUnilateralLiquidity) Route() string { return RouterKey }

func (m MsgRemoveUnilateralLiquidity) Type() string { return TypeMsgRemoveUnilateralLiquidity }

func (m MsgRemoveUnilateralLiquidity) ValidateBasic() error {
	if err := ValidateCounterpartyDenom(m.CounterpartyDenom); err != nil {
		return err
	}

	if err := ValidateToken(m.MinToken); err != nil {
		return err
	}

	if err := ValidateLiquidity(m.ExactLiquidity); err != nil {
		return err
	}

	if err := ValidateDeadline(m.Deadline); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}

func (m MsgRemoveUnilateralLiquidity) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgRemoveUnilateralLiquidity) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// GetSignBytes returns the raw bytes for a MsgUpdateParams message that
// the expected signer needs to sign.
func (m *MsgUpdateParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic executes sanity validation on the provided data
func (m *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}
	return m.Params.Validate()
}

// GetSigners returns the expected signers for a MsgUpdateParams message
func (m *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}
