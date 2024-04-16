package v1

import (
	fmt "fmt"
	"regexp"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	tokentypes "github.com/irisnet/irismod/modules/token/types"
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
	_ sdk.Msg = &MsgSwapFeeToken{}
	_ sdk.Msg = &MsgUpdateParams{}
	_ sdk.Msg = &MsgDeployERC20{}
	_ sdk.Msg = &MsgSwapFromERC20{}
	_ sdk.Msg = &MsgSwapToERC20{}

	regexpERC20Fmt = fmt.Sprintf("^[a-z][a-z0-9/]{%d,%d}$", tokentypes.MinimumSymbolLen-1, tokentypes.MaximumSymbolLen-1)
	regexpERC20    = regexp.MustCompile(regexpERC20Fmt).MatchString
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
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	return NewToken(
		msg.Symbol,
		msg.Name,
		msg.MinUnit,
		msg.Scale,
		msg.InitialSupply,
		msg.MaxSupply,
		msg.Mintable,
		owner,
	).Validate()
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
		return errorsmod.Wrapf(
			sdkerrors.ErrInvalidAddress,
			"invalid source owner address (%s)",
			err,
		)
	}

	dstOwner, err := sdk.AccAddressFromBech32(msg.DstOwner)
	if err != nil {
		return errorsmod.Wrapf(
			sdkerrors.ErrInvalidAddress,
			"invalid destination owner address (%s)",
			err,
		)
	}

	// check if the `DstOwner` is same as the original owner
	if srcOwner.Equals(dstOwner) {
		return tokentypes.ErrInvalidToAddress
	}

	// check the symbol
	if err := tokentypes.ValidateSymbol(msg.Symbol); err != nil {
		return err
	}

	return nil
}

// Route implements Msg
func (msg MsgTransferTokenOwner) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgTransferTokenOwner) Type() string { return TypeMsgTransferTokenOwner }

// NewMsgEditToken creates a MsgEditToken
func NewMsgEditToken(
	name, symbol string,
	maxSupply uint64,
	mintable tokentypes.Bool,
	owner string,
) *MsgEditToken {
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
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	if err := tokentypes.ValidateName(msg.Name); err != nil {
		return err
	}
	// check symbol
	return tokentypes.ValidateSymbol(msg.Symbol)
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
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	// check the reception
	if len(msg.To) > 0 {
		if _, err := sdk.AccAddressFromBech32(msg.To); err != nil {
			return errorsmod.Wrapf(
				sdkerrors.ErrInvalidAddress,
				"invalid mint reception address (%s)",
				err,
			)
		}
	}

	return tokentypes.ValidateCoin(msg.Coin)
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
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	return tokentypes.ValidateCoin(msg.Coin)
}

// GetSigners implements Msg
func (msg MsgSwapFeeToken) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// GetSignBytes Implements Msg.
func (msg MsgSwapFeeToken) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg
func (msg MsgSwapFeeToken) ValidateBasic() error {
	// check the Sender
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if len(msg.Recipient) != 0 {
		if _, err := sdk.AccAddressFromBech32(msg.Recipient); err != nil {
			return errorsmod.Wrapf(
				sdkerrors.ErrInvalidAddress,
				"invalid recipient address (%s)",
				err,
			)
		}
	}

	return tokentypes.ValidateCoin(msg.FeePaid)
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

// ValidateBasic implements Msg
func (m *MsgDeployERC20) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	if err := tokentypes.ValidateName(m.Name); err != nil {
		return err
	}

	if err := tokentypes.ValidateScale(m.Scale); err != nil {
		return err
	}

	if err := ValidateERC20(m.MinUnit); err != nil {
		return err
	}
	return ValidateERC20(m.Symbol)
}

// GetSigners returns the expected signers for a MsgDeployERC20 message
func (m *MsgDeployERC20) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic implements Msg
func (m *MsgSwapFromERC20) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if _, err := sdk.AccAddressFromBech32(m.Receiver); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receiver address (%s)", err)
	}

	if !m.WantedAmount.IsValid() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, m.WantedAmount.String())
	}

	if !m.WantedAmount.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, m.WantedAmount.String())
	}
	return nil
}

// GetSigners returns the expected signers for a MsgSwapFromERC20 message
func (m *MsgSwapFromERC20) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Sender)
	return []sdk.AccAddress{addr}
}

// ValidateBasic implements Msg
func (m *MsgSwapToERC20) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if tokentypes.IsValidEthAddress(m.Receiver) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "expecting a hex address of 0x, got %s", m.Receiver)
	}

	if !m.Amount.IsValid() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, m.Amount.String())
	}

	if !m.Amount.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, m.Amount.String())
	}
	return nil
}

// GetSigners returns the expected signers for a MsgSwapToERC20 message
func (m *MsgSwapToERC20) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Sender)
	return []sdk.AccAddress{addr}
}

// ValidateERC20 validates ERC20 symbol or name
func ValidateERC20(params string) error {
	if !regexpERC20(params) {
		return errorsmod.Wrapf(
			tokentypes.ErrInvalidSymbol,
			"invalid symbol or name: %s, only accepts english lowercase letters, numbers or slash, length [%d, %d], and begin with an english letter, regexp: %s",
			params,
			tokentypes.MinimumSymbolLen,
			tokentypes.MaximumSymbolLen,
			regexpERC20Fmt,
		)
	}
	return nil
}
