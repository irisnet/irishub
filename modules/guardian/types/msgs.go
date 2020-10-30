package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgAddSuper    = "add_super"    // type for MsgAddSuper
	TypeMsgDeleteSuper = "delete_super" // type for MsgDeleteSuper
)

var (
	_ sdk.Msg = &MsgAddSuper{}
	_ sdk.Msg = &MsgDeleteSuper{}
)

// NewMsgAddSuper constructs a MsgAddSuper
func NewMsgAddSuper(description string, address, addedBy sdk.AccAddress) *MsgAddSuper {
	return &MsgAddSuper{
		Description: description,
		Address:     address.String(),
		AddedBy:     addedBy.String(),
	}
}

// Route implements Msg.
func (msg MsgAddSuper) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgAddSuper) Type() string { return TypeMsgAddSuper }

// GetSignBytes implements Msg.
func (msg MsgAddSuper) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgAddSuper) ValidateBasic() error {
	if len(msg.Description) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "description missing")
	}
	if len(msg.Address) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "added address missing")
	}
	if len(msg.AddedBy) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "operator address missing")
	}
	if err := msg.EnsureLength(); err != nil {
		return err
	}
	return nil
}

// GetSigners implements Msg.
func (msg MsgAddSuper) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.AddedBy)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ______________________________________________________________________

// NewMsgDeleteSuper constructs a MsgDeleteSuper
func NewMsgDeleteSuper(address, deletedBy sdk.AccAddress) *MsgDeleteSuper {
	return &MsgDeleteSuper{
		Address:   address.String(),
		DeletedBy: deletedBy.String(),
	}
}

// Route implements Msg.
func (msg MsgDeleteSuper) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgDeleteSuper) Type() string { return TypeMsgDeleteSuper }

// GetSignBytes implements Msg.
func (msg MsgDeleteSuper) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// RoValidateBasicute implements Msg.
func (msg MsgDeleteSuper) ValidateBasic() error {
	if len(msg.Address) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "deleted address missing")
	}
	if len(msg.DeletedBy) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "operator address missing")
	}
	return nil
}

// GetSigners implements Msg.
func (msg MsgDeleteSuper) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.DeletedBy)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// EnsureLength validate the length of AddGuardian
func (msg MsgAddSuper) EnsureLength() error {
	if len(msg.Description) > 70 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid website length; got: %d, max: %d", len(msg.Description), 70)
	}
	return nil
}
