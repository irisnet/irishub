package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	asset "github.com/irisnet/irishub/modules/asset/types"
)

const (
	TypeMsgIssueToken = "issue_token" // type for MsgIssueToken

	// constant used to indicate that some field should not be updated
	DoNotModify = "[do-not-modify]"
)

var (
	_ sdk.Msg = &MsgIssueToken{}
	_ sdk.Msg = &MsgEditToken{}
	_ sdk.Msg = &MsgMintToken{}
	_ sdk.Msg = &MsgTransferToken{}
	_ sdk.Msg = &MsgBurnToken{}
)

// MsgIssueToken for issuing token
type MsgIssueToken struct {
	Symbol        string         `json:"symbol" yaml:"symbol"`                 //globally unique token identifier
	Name          string         `json:"name" yaml:"name"`                     //the name of the token
	Scale         uint8          `json:"scale" yaml:"scale"`                   //maximum number of decimals supported by this token
	MinUnit       string         `json:"min_unit" yaml:"min_unit"`             //the smallest unit name of the token
	InitialSupply uint64         `json:"initial_supply" yaml:"initial_supply"` //initial Token Issuance
	MaxSupply     uint64         `json:"max_supply" yaml:"max_supply"`         //maximum Token Issuance
	Mintable      bool           `json:"mintable" yaml:"mintable"`             //is it possible to issue additional shares after the token is issued?
	Owner         sdk.AccAddress `json:"owner" yaml:"owner"`                   //the actual controller of the token
}

// NewMsgIssueToken - construct token issue msg
func NewMsgIssueToken(symbol, name string, scale uint8,
	minUnit string, initialSupply uint64, maxSupply uint64, mintable bool, owner sdk.AccAddress,
) MsgIssueToken {
	return MsgIssueToken{
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

// Route implements Msg.
func (msg MsgIssueToken) Route() string { return asset.RouterKey }

// Type implements Msg.
func (msg MsgIssueToken) Type() string { return TypeMsgIssueToken }

// ValidateMsgIssueToken - validate msg
func ValidateMsgIssueToken(msg *MsgIssueToken) error {
	msg.Symbol = strings.ToLower(strings.TrimSpace(msg.Symbol))
	msg.MinUnit = strings.ToLower(strings.TrimSpace(msg.MinUnit))
	msg.Name = strings.TrimSpace(msg.Name)

	if msg.MaxSupply == 0 {
		if msg.Mintable {
			msg.MaxSupply = MaximumTokenMaxSupply
		} else {
			msg.MaxSupply = msg.InitialSupply
		}
	}

	if err := ValidateName(msg.Name); err != nil {
		return err
	}

	if err := ValidateSymbol(msg.Symbol); err != nil {
		return err
	}

	if err := ValidateMinUnit(msg.MinUnit); err != nil {
		return err
	}

	if err := ValidateSupply(msg.InitialSupply, msg.MaxSupply); err != nil {
		return err
	}
	return ValidateScale(msg.Scale)
}

// ValidateBasic implements Msg.
func (msg MsgIssueToken) ValidateBasic() error {
	return ValidateMsgIssueToken(&msg)
}

// GetSignBytes implements Msg.
func (msg MsgIssueToken) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg.
func (msg MsgIssueToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgTransferToken for transferring the token owner
type MsgTransferToken struct {
	Symbol   string         `json:"symbol" yaml:"symbol"`       //the token symbol
	SrcOwner sdk.AccAddress `json:"src_owner" yaml:"src_owner"` // the current owner address of the token
	DstOwner sdk.AccAddress `json:"dst_owner" yaml:"dst_owner"` // the new owner
}

// NewMsgTransferToken - construct token transfer msg
func NewMsgTransferToken(srcOwner, dstOwner sdk.AccAddress, symbol string) MsgTransferToken {
	return MsgTransferToken{
		SrcOwner: srcOwner,
		DstOwner: dstOwner,
		Symbol:   strings.TrimSpace(symbol),
	}
}

// GetSignBytes implements Msg.
func (msg MsgTransferToken) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg.
func (msg MsgTransferToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.SrcOwner}
}

func (msg MsgTransferToken) ValidateBasic() error {
	// check the SrcOwner
	if len(msg.SrcOwner) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "the owner of the token missing")
	}

	// check if the `DstOwner` is empty
	if len(msg.DstOwner) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "the new owner of the token missing")
	}

	// check if the `DstOwner` is same as the original owner
	if msg.SrcOwner.Equals(msg.DstOwner) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "the new owner must not be same as the original owner")
	}

	//check the Symbol
	return ValidateSymbol(msg.Symbol)
}

// Route implements Msg.
func (msg MsgTransferToken) Route() string { return asset.RouterKey }

// Type implements Msg.
func (msg MsgTransferToken) Type() string { return "transfer_token" }

// MsgEditToken for editing a specified token
type MsgEditToken struct {
	Symbol    string         `json:"symbol" yaml:"symbol"`         //the token symbol
	Name      string         `json:"name" yaml:"name"`             // token name
	MaxSupply uint64         `json:"max_supply" yaml:"max_supply"` // max supply
	Mintable  Bool           `json:"mintable" yaml:"mintable"`     // mintable of token
	Owner     sdk.AccAddress `json:"owner" yaml:"owner"`           // owner of token
}

// NewMsgEditToken creates a MsgEditToken
func NewMsgEditToken(name, symbol string, maxSupply uint64, mintable Bool, owner sdk.AccAddress) MsgEditToken {
	return MsgEditToken{
		Name:      strings.TrimSpace(name),
		Symbol:    strings.ToLower(strings.TrimSpace(symbol)),
		MaxSupply: maxSupply,
		Mintable:  mintable,
		Owner:     owner,
	}
}

// Route implements Msg.
func (msg MsgEditToken) Route() string { return asset.RouterKey }

// Type implements Msg.
func (msg MsgEditToken) Type() string { return "edit_token" }

// ValidateBasic implements Msg.
func (msg MsgEditToken) ValidateBasic() error {
	//check owner
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "the owner of the token missing")
	}

	//check max_supply for fast failed
	if err := ValidateMaxSupply(msg.MaxSupply); err != nil {
		return err
	}

	if err := ValidateName(msg.Name); DoNotModify != msg.Name && err != nil {
		return err
	}

	return ValidateSymbol(msg.Symbol)
}

// GetSignBytes implements Msg.
func (msg MsgEditToken) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg.
func (msg MsgEditToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgMintToken for mint the token to a specified address
type MsgMintToken struct {
	Symbol string         `json:"symbol" yaml:"symbol"` //the token symbol
	Owner  sdk.AccAddress `json:"owner" yaml:"owner"`   // the current owner address of the token
	To     sdk.AccAddress `json:"to" yaml:"to"`         // address of mint token to
	Amount uint64         `json:"amount" yaml:"amount"` // amount of mint token
}

// NewMsgMintToken creates a MsgMintToken
func NewMsgMintToken(symbol string, owner, to sdk.AccAddress, amount uint64) MsgMintToken {
	return MsgMintToken{
		Symbol: strings.TrimSpace(symbol),
		Owner:  owner,
		To:     to,
		Amount: amount,
	}
}

// Route implements Msg.
func (msg MsgMintToken) Route() string { return asset.RouterKey }

// Type implements Msg.
func (msg MsgMintToken) Type() string { return "mint_token" }

// GetSignBytes implements Msg.
func (msg MsgMintToken) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg.
func (msg MsgMintToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// ValidateBasic implements Msg.
func (msg MsgMintToken) ValidateBasic() error {
	// check the owner
	if len(msg.Owner) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "the owner of the token missing")
	}

	if err := ValidateSymbol(msg.Symbol); err != nil {
		return err
	}

	if msg.Amount <= 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("the amount of the token must be great than zero: %d", msg.Amount))
	}

	return ValidateMaxSupply(msg.Amount)
}

type MsgBurnToken struct {
	Sender sdk.AccAddress `json:"sender" yaml:"sender"`
	Amount sdk.Coins      `json:"amount" yaml:"amount"`
}

// NewMsgMintToken creates a MsgMintToken
func NewMsgBurnToken(sender sdk.AccAddress, amount sdk.Coins) MsgBurnToken {
	return MsgBurnToken{
		Sender: sender,
		Amount: amount,
	}
}

// Route implements Msg
func (msg MsgBurnToken) Route() string { return asset.RouterKey }

// Type implements Msg
func (msg MsgBurnToken) Type() string { return "burn_token" }

// GetSignBytes implements Msg
func (msg MsgBurnToken) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgBurnToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// ValidateBasic implements Msg
func (msg MsgBurnToken) ValidateBasic() error {
	// check the Sender
	if len(msg.Sender) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "the sender of the token missing")
	}

	if msg.Amount.Empty() || !msg.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	return nil
}
