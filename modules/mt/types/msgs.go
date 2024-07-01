package types

import (
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// constant used to indicate that some field should not be updated
const (
	TypeMsgIssueDenom    = "issue_denom"
	TypeMsgTransferDenom = "transfer_denom"

	TypeMsgMintMT     = "mint_mt"
	TypeMsgTransferMT = "transfer_mt"
	TypeMsgEditMT     = "edit_mt"
	TypeMsgBurnMT     = "burn_mt"
)

var (
	_ sdk.Msg = &MsgIssueDenom{}
	_ sdk.Msg = &MsgTransferDenom{}

	_ sdk.Msg = &MsgMintMT{}
	_ sdk.Msg = &MsgTransferMT{}
	_ sdk.Msg = &MsgEditMT{}
	_ sdk.Msg = &MsgBurnMT{}
)

// NewMsgIssueDenom is a constructor function for MsgIssueDenom
func NewMsgIssueDenom(name, data, sender string) *MsgIssueDenom {
	return &MsgIssueDenom{
		Name:   name,
		Data:   []byte(data),
		Sender: sender,
	}
}

// Route Implements Msg
func (msg MsgIssueDenom) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgIssueDenom) Type() string { return TypeMsgIssueDenom }

// ValidateBasic Implements Msg.
func (msg MsgIssueDenom) ValidateBasic() error {
	if len(strings.TrimSpace(msg.Name)) <= 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "name is required")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgIssueDenom) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgIssueDenom) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// NewMsgTransferMT is a constructor function for MsgTransferMT
func NewMsgTransferMT(
	mtID, denomID, sender, recipient string, amount uint64,
) *MsgTransferMT {
	return &MsgTransferMT{
		Id:        mtID,
		DenomId:   denomID,
		Amount:    amount,
		Sender:    sender,
		Recipient: recipient,
	}
}

// Route Implements Msg
func (msg MsgTransferMT) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgTransferMT) Type() string { return TypeMsgTransferMT }

// ValidateBasic Implements Msg.
func (msg MsgTransferMT) ValidateBasic() error {
	if len(strings.TrimSpace(msg.Id)) <= 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "mt id is required")
	}

	if len(strings.TrimSpace(msg.DenomId)) <= 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "denom id is required")
	}

	if msg.Amount <= 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "amount is required")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Recipient); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid recipient address (%s)", err)
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgTransferMT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgTransferMT) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// NewMsgEditMT is a constructor function for MsgEditMT
func NewMsgEditMT(
	mtID, denomID, tokenData, sender string,
) *MsgEditMT {
	return &MsgEditMT{
		Id:      mtID,
		DenomId: denomID,
		Data:    []byte(tokenData),
		Sender:  sender,
	}
}

// Route Implements Msg
func (msg MsgEditMT) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgEditMT) Type() string { return TypeMsgEditMT }

// ValidateBasic Implements Msg.
func (msg MsgEditMT) ValidateBasic() error {
	if len(strings.TrimSpace(msg.Id)) <= 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "mt id is required")
	}

	if len(strings.TrimSpace(msg.DenomId)) <= 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "denom id is required")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgEditMT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgEditMT) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// NewMsgMintMT is a constructor function for MsgMintMT
func NewMsgMintMT(
	mtID, denomID string, amount uint64, tokenData, sender, recipient string,
) *MsgMintMT {
	return &MsgMintMT{
		Id:        mtID,
		DenomId:   denomID,
		Amount:    amount,
		Data:      []byte(tokenData),
		Sender:    sender,
		Recipient: recipient,
	}
}

// Route Implements Msg
func (msg MsgMintMT) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgMintMT) Type() string { return TypeMsgMintMT }

// ValidateBasic Implements Msg.
func (msg MsgMintMT) ValidateBasic() error {
	if len(strings.TrimSpace(msg.DenomId)) <= 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "denom id is required")
	}

	if msg.Amount <= 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "amount is required")
	}

	if len(strings.TrimSpace(msg.Id)) > 0 && len(msg.Data) > 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "metadata can not be accepted while minting, use 'edit mt' instead")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if len(strings.TrimSpace(msg.Recipient)) > 0 {
		if _, err := sdk.AccAddressFromBech32(msg.Recipient); err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receipt address (%s)", err)
		}
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgMintMT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgMintMT) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// NewMsgBurnMT is a constructor function for MsgBurnMT
func NewMsgBurnMT(sender, mtID, denomID string, amount uint64) *MsgBurnMT {
	return &MsgBurnMT{
		Sender:  sender,
		Id:      mtID,
		DenomId: denomID,
		Amount:  amount,
	}
}

// Route Implements Msg
func (msg MsgBurnMT) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgBurnMT) Type() string { return TypeMsgBurnMT }

// ValidateBasic Implements Msg.
func (msg MsgBurnMT) ValidateBasic() error {
	if len(strings.TrimSpace(msg.Id)) <= 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "mt id is required")
	}

	if len(strings.TrimSpace(msg.DenomId)) <= 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "denom id is required")
	}

	if msg.Amount <= 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "amount is required")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgBurnMT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgBurnMT) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// NewMsgTransferDenom is a constructor function for MsgTransferDenom
func NewMsgTransferDenom(denomId, sender, recipient string) *MsgTransferDenom {
	return &MsgTransferDenom{
		Id:        denomId,
		Sender:    sender,
		Recipient: recipient,
	}
}

// Route Implements Msg
func (msg MsgTransferDenom) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgTransferDenom) Type() string { return TypeMsgTransferDenom }

// ValidateBasic Implements Msg.
func (msg MsgTransferDenom) ValidateBasic() error {
	if len(strings.TrimSpace(msg.Id)) <= 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "denom id is required")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Recipient); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid recipient address (%s)", err)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgTransferDenom) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgTransferDenom) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}
