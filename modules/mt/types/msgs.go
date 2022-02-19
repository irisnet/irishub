package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// constant used to indicate that some field should not be updated
const (
	TypeMsgIssueDenom    = "issue_denom"
	TypeMsgTransferMT    = "transfer_mt"
	TypeMsgEditMT        = "edit_mt"
	TypeMsgMintMT        = "mint_mt"
	TypeMsgBurnMT        = "burn_mt"
	TypeMsgTransferDenom = "transfer_denom"
)

var (
	_ sdk.Msg = &MsgIssueDenom{}
	_ sdk.Msg = &MsgTransferMT{}
	_ sdk.Msg = &MsgEditMT{}
	_ sdk.Msg = &MsgMintMT{}
	_ sdk.Msg = &MsgBurnMT{}
	_ sdk.Msg = &MsgTransferDenom{}
)

// NewMsgIssueDenom is a constructor function for MsgSetName
func NewMsgIssueDenom(name string, data []byte, sender string) *MsgIssueDenom {
	return &MsgIssueDenom{
		Name:   name,
		Data:   data,
		Sender: sender,
	}
}

// Route Implements Msg
func (msg MsgIssueDenom) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgIssueDenom) Type() string { return TypeMsgIssueDenom }

// ValidateBasic Implements Msg.
func (msg MsgIssueDenom) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
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

// TODO add amount
// NewMsgTransferMT is a constructor function for MsgSetName
func NewMsgTransferMT(
	tokenID, denomID, sender, recipient string, amount uint64,
) *MsgTransferMT {
	return &MsgTransferMT{
		Id:        tokenID,
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
	if err := ValidateDenomID(msg.DenomId); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Recipient); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid recipient address (%s)", err)
	}
	return ValidateTokenID(msg.Id)
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

// NewMsgEditMT is a constructor function for MsgSetName
func NewMsgEditMT(
	tokenID, denomID, tokenData, sender string,
) *MsgEditMT {
	return &MsgEditMT{
		Id:      tokenID,
		DenomId: denomID,
		Data:    []byte(tokenData),
	}
}

// Route Implements Msg
func (msg MsgEditMT) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgEditMT) Type() string { return TypeMsgEditMT }

// ValidateBasic Implements Msg.
func (msg MsgEditMT) ValidateBasic() error {
	if err := ValidateDenomID(msg.DenomId); err != nil {
		return err
	}

	return ValidateTokenID(msg.Id)
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
	tokenID, denomID string, amount uint64, tokenData, sender, recipient string,
) *MsgMintMT {
	return &MsgMintMT{
		Id:        tokenID,
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
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if _, err := sdk.AccAddressFromBech32(msg.Recipient); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receipt address (%s)", err)
	}
	if err := ValidateDenomID(msg.DenomId); err != nil {
		return err
	}
	if err := ValidateKeywords(msg.DenomId); err != nil {
		return err
	}
	return ValidateTokenID(msg.Id)
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
func NewMsgBurnMT(sender, tokenID, denomID string, amount uint64) *MsgBurnMT {
	return &MsgBurnMT{
		Sender:  sender,
		Id:      tokenID,
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
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if err := ValidateDenomID(msg.DenomId); err != nil {
		return err
	}
	return ValidateTokenID(msg.Id)
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

// NewMsgTransferDenom is a constructor function for msgTransferDenom
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
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if _, err := sdk.AccAddressFromBech32(msg.Recipient); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid recipient address (%s)", err)
	}
	if err := ValidateDenomID(msg.Id); err != nil {
		return err
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
