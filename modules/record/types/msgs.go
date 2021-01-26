package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateRecord = "create_record" // type for TypeMsgCreateRecord
)

var _ sdk.Msg = &MsgCreateRecord{}

// NewMsgCreateRecord constructs a MsgCreateRecord
func NewMsgCreateRecord(contents []Content, Creator string) *MsgCreateRecord {
	return &MsgCreateRecord{
		Contents: contents,
		Creator:  Creator,
	}
}

// Route implements Msg.
func (msg MsgCreateRecord) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgCreateRecord) Type() string { return TypeMsgCreateRecord }

// GetSignBytes implements Msg.
func (msg MsgCreateRecord) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgCreateRecord) ValidateBasic() error {
	if len(msg.Contents) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "contents missing")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return ValidateContents(msg.Contents...)
}

// GetSigners implements Msg.
func (msg MsgCreateRecord) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}
