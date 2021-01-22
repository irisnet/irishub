package types

import (
	"strings"

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
	msg = msg.Normalize()
	if len(msg.Contents) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "contents missing")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return ValidateContents(msg.Contents...)
}

// Normalize return a string with spaces removed and lowercase
func (msg MsgCreateRecord) Normalize() MsgCreateRecord {
	for i, ctx := range msg.Contents {
		ctx.Digest = strings.TrimSpace(ctx.Digest)
		ctx.DigestAlgo = strings.TrimSpace(ctx.DigestAlgo)
		ctx.URI = strings.TrimSpace(ctx.URI)
		ctx.Meta = strings.TrimSpace(ctx.Meta)
		msg.Contents[i] = ctx
	}
	return msg
}

// GetSigners implements Msg.
func (msg MsgCreateRecord) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}
