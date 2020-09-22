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
func NewMsgCreateRecord(contents []Content, Creator sdk.AccAddress) *MsgCreateRecord {
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
	if msg.Creator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "creator missing")
	}

	for i, content := range msg.Contents {
		if len(content.Digest) == 0 {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "content[%d] digest missing", i)
		}
		if len(content.DigestAlgo) == 0 {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "content[%d] digest algo missing", i)
		}
	}
	return nil
}

// GetSigners implements Msg.
func (msg MsgCreateRecord) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}
