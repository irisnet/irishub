package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// MsgRoute identifies transaction types
	MsgRoute = "token"

	TypeMsgIssueToken         = "issue_token"
	TypeMsgEditToken          = "edit_token"
	TypeMsgMintToken          = "mint_token"
	TypeMsgBurnToken          = "burn_token"
	TypeMsgTransferTokenOwner = "transfer_token_owner"

	// DoNotModify used to indicate that some field should not be updated
	DoNotModify = "[do-not-modify]"
)

var (
	_ sdk.Msg = &MsgIssueToken{}
	_ sdk.Msg = &MsgEditToken{}
	_ sdk.Msg = &MsgMintToken{}
	_ sdk.Msg = &MsgBurnToken{}
	_ sdk.Msg = &MsgTransferTokenOwner{}
)

// NewMsgIssueToken - construct token issue msg.
func NewMsgIssueToken(
	symbol string, minUnit string, name string,
	scale uint32, initialSupply, maxSupply uint64,
	mintable bool, owner string,
) *MsgIssueToken {
	return &MsgIssueToken{
		Symbol:        symbol,
		Name:          name,
		Scale:         scale,
		MinUnit:       minUnit,
		InitialSupply: initialSupply,
		MaxSupply:     maxSupply,
		Mintable:      mintable,
		Owner:         owner,
	}
}

// Route Implements Msg.
func (msg MsgIssueToken) Route() string { return MsgRoute }

// Type Implements Msg.
func (msg MsgIssueToken) Type() string { return TypeMsgIssueToken }

// ValidateBasic Implements Msg.
func (msg MsgIssueToken) ValidateBasic() error {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	return ValidateToken(
		NewToken(
			msg.Symbol,
			msg.Name,
			msg.MinUnit,
			msg.Scale,
			msg.InitialSupply,
			msg.MaxSupply,
			msg.Mintable,
			owner,
		),
	)
}

// GetSignBytes Implements Msg.
func (msg MsgIssueToken) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgIssueToken) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// NewMsgTransferTokenOwner return a instance of MsgTransferTokenOwner
func NewMsgTransferTokenOwner(srcOwner, dstOwner, symbol string) *MsgTransferTokenOwner {
	return &MsgTransferTokenOwner{
		SrcOwner: srcOwner,
		DstOwner: dstOwner,
		Symbol:   symbol,
	}
}

// GetSignBytes implements Msg
func (msg MsgTransferTokenOwner) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgTransferTokenOwner) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.SrcOwner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ValidateBasic implements Msg
func (msg MsgTransferTokenOwner) ValidateBasic() error {
	srcOwner, err := sdk.AccAddressFromBech32(msg.SrcOwner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid source owner address (%s)", err)
	}

	dstOwner, err := sdk.AccAddressFromBech32(msg.DstOwner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid destination owner address (%s)", err)
	}

	// check if the `DstOwner` is same as the original owner
	if srcOwner.Equals(dstOwner) {
		return ErrInvalidToAddress
	}

	// check the symbol
	if err := ValidateSymbol(msg.Symbol); err != nil {
		return err
	}

	return nil
}

// Route implements Msg
func (msg MsgTransferTokenOwner) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgTransferTokenOwner) Type() string { return TypeMsgTransferTokenOwner }

// NewMsgEditToken creates a MsgEditToken
func NewMsgEditToken(name, symbol string, maxSupply uint64, mintable Bool, owner string) *MsgEditToken {
	return &MsgEditToken{
		Name:      name,
		Symbol:    symbol,
		MaxSupply: maxSupply,
		Mintable:  mintable,
		Owner:     owner,
	}
}

// Route implements Msg
func (msg MsgEditToken) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgEditToken) Type() string { return TypeMsgEditToken }

// GetSignBytes implements Msg
func (msg MsgEditToken) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgEditToken) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ValidateBasic implements Msg
func (msg MsgEditToken) ValidateBasic() error {
	// check owner
	if _, err := sdk.AccAddressFromBech32(msg.Owner); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	if err := ValidateName(msg.Name); err != nil {
		return err
	}
	// check symbol
	return ValidateSymbol(msg.Symbol)
}

// NewMsgMintToken creates a MsgMintToken
func NewMsgMintToken(symbol, owner, to string, amount uint64) *MsgMintToken {
	return &MsgMintToken{
		Symbol: symbol,
		Owner:  owner,
		To:     to,
		Amount: amount,
	}
}

// Route implements Msg
func (msg MsgMintToken) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgMintToken) Type() string { return TypeMsgMintToken }

// GetSignBytes implements Msg
func (msg MsgMintToken) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgMintToken) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ValidateBasic implements Msg
func (msg MsgMintToken) ValidateBasic() error {
	// check the owner
	if _, err := sdk.AccAddressFromBech32(msg.Owner); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	// check the reception
	if len(msg.To) > 0 {
		if _, err := sdk.AccAddressFromBech32(msg.To); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid mint reception address (%s)", err)
		}
	}

	if err := ValidateAmount(msg.Amount); err != nil {
		return err
	}

	return ValidateSymbol(msg.Symbol)
}

// NewMsgBurnToken creates a MsgMintToken
func NewMsgBurnToken(symbol string, owner string, amount uint64) *MsgBurnToken {
	return &MsgBurnToken{
		Symbol: symbol,
		Amount: amount,
		Sender: owner,
	}
}

// Route implements Msg
func (msg MsgBurnToken) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgBurnToken) Type() string { return TypeMsgBurnToken }

// GetSignBytes implements Msg
func (msg MsgBurnToken) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgBurnToken) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ValidateBasic implements Msg
func (msg MsgBurnToken) ValidateBasic() error {
	// check the owner
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	if err := ValidateAmount(msg.Amount); err != nil {
		return err
	}

	return ValidateSymbol(msg.Symbol)
}
