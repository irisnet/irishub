package types

import (
	"fmt"
	"regexp"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// MsgRoute identifies transaction types
	MsgRoute          = "asset"
	MsgTypeIssueToken = "issue_token"

	// constant used to indicate that some field should not be updated
	DoNotModify = "[do-not-modify]"
)

var (
	MaximumAssetMaxSupply        = uint64(1000000000000) // maximal limitation for asset max supply，1000 billion
	MaximumAssetInitSupply       = uint64(100000000000)  // maximal limitation for asset initial supply，100 billion
	MaximumAssetDecimal          = uint8(18)             // maximal limitation for asset decimal
	MinimumAssetSymbolSize       = 3                     // minimal limitation for the length of the asset's symbol / canonical_symbol
	MaximumAssetSymbolSize       = 8                     // maximal limitation for the length of the asset's symbol / canonical_symbol
	MinimumAssetMinUnitAliasSize = 3                     // minimal limitation for the length of the asset's min_unit_alias
	MaximumAssetMinUnitAliasSize = 10                    // maximal limitation for the length of the asset's min_unit_alias
	MaximumAssetNameSize         = 32                    // maximal limitation for the length of the asset's name

	MinimumGatewayMonikerSize  = 3   // minimal limitation for the length of the gateway's moniker
	MaximumGatewayMonikerSize  = 8   // maximal limitation for the length of the gateway's moniker
	MaximumGatewayIdentitySize = 128 // maximal limitation for the length of the gateway's identity
	MaximumGatewayDetailsSize  = 280 // maximal limitation for the length of the gateway's details
	MaximumGatewayWebsiteSize  = 128 // maximal limitation for the length of the gateway's website

	IsAlphaNumeric     = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString   // only accepts alphanumeric characters
	IsAlphaNumericDash = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString // only accepts alphanumeric characters, _ and -
	IsBeginWithAlpha   = regexp.MustCompile(`^[a-zA-Z].*`).MatchString
)

var _, _, _, _ sdk.Msg = &MsgIssueToken{}, &MsgEditToken{}, &MsgMintToken{}, &MsgTransferTokenOwner{}

// MsgIssueToken for issuing token
type MsgIssueToken struct {
	Family          AssetFamily    `json:"family"`
	Source          AssetSource    `json:"source"`
	Symbol          string         `json:"symbol"`
	CanonicalSymbol string         `json:"canonical_symbol"`
	Name            string         `json:"name"`
	Decimal         uint8          `json:"decimal"`
	MinUnitAlias    string         `json:"min_unit_alias"`
	InitialSupply   uint64         `json:"initial_supply"`
	MaxSupply       uint64         `json:"max_supply"`
	Mintable        bool           `json:"mintable"`
	Owner           sdk.AccAddress `json:"owner"`
}

// NewMsgIssueToken - construct token issue msg
func NewMsgIssueToken(family AssetFamily, source AssetSource, symbol string, canonicalSymbol string, name string,
	decimal uint8, alias string, initialSupply uint64, maxSupply uint64, mintable bool, owner sdk.AccAddress,
) MsgIssueToken {
	return MsgIssueToken{
		Family:          family,
		Source:          source,
		Symbol:          symbol,
		CanonicalSymbol: canonicalSymbol,
		Name:            name,
		Decimal:         decimal,
		MinUnitAlias:    alias,
		InitialSupply:   initialSupply,
		MaxSupply:       maxSupply,
		Mintable:        mintable,
		Owner:           owner,
	}
}

// Route Implements Msg
func (msg MsgIssueToken) Route() string { return MsgRoute }

// Type Implements Msg
func (msg MsgIssueToken) Type() string { return MsgTypeIssueToken }

// ValidateMsgIssueToken - validate msg
func ValidateMsgIssueToken(msg *MsgIssueToken) sdk.Error {
	msg.Symbol = strings.ToLower(strings.TrimSpace(msg.Symbol))
	msg.CanonicalSymbol = strings.ToLower(strings.TrimSpace(msg.CanonicalSymbol))
	msg.MinUnitAlias = strings.ToLower(strings.TrimSpace(msg.MinUnitAlias))
	msg.Name = strings.TrimSpace(msg.Name)

	if msg.MaxSupply == 0 {
		if msg.Mintable {
			msg.MaxSupply = MaximumAssetMaxSupply
		} else {
			msg.MaxSupply = msg.InitialSupply
		}
	}

	switch msg.Source {
	case NATIVE:
		// require owner for native asset
		if msg.Owner.Empty() {
			return ErrNilAssetOwner(DefaultCodespace, "the owner of the asset must be specified")
		}
		// ignore CanonicalSymbol for native asset
		msg.CanonicalSymbol = ""
		break
	case EXTERNAL:
		break
	default:
		return ErrInvalidAssetSource(DefaultCodespace, fmt.Sprintf("invalid asset source type %s", msg.Source))
	}

	if _, found := AssetFamilyToStringMap[msg.Family]; !found {
		return ErrInvalidAssetFamily(DefaultCodespace, fmt.Sprintf("invalid asset family type %s", msg.Family))
	}

	nameLen := len(msg.Name)
	if nameLen == 0 || nameLen > MaximumAssetNameSize {
		return ErrInvalidAssetName(DefaultCodespace, fmt.Sprintf("invalid token name %s, only accepts length (0, %d]", msg.Name, MaximumAssetNameSize))
	}

	symbolLen := len(msg.Symbol)
	if symbolLen < MinimumAssetSymbolSize || symbolLen > MaximumAssetSymbolSize || !IsBeginWithAlpha(msg.Symbol) || !IsAlphaNumeric(msg.Symbol) {
		return ErrInvalidAssetSymbol(DefaultCodespace, fmt.Sprintf("invalid token symbol %s, only accepts alphanumeric characters, and begin with an english letter, length [%d, %d]", msg.Symbol, MinimumAssetSymbolSize, MaximumAssetSymbolSize))
	}

	//if strings.Contains(strings.ToLower(msg.Symbol), sdk.Iris) {
	//	return ErrInvalidAssetSymbol(DefaultCodespace, fmt.Sprintf("invalid token symbol %s, can not contain native token symbol %s", msg.Symbol, sdk.Iris))
	//}

	minUnitAliasLen := len(msg.MinUnitAlias)
	if minUnitAliasLen > 0 && (minUnitAliasLen < MinimumAssetMinUnitAliasSize || minUnitAliasLen > MaximumAssetMinUnitAliasSize || !IsAlphaNumeric(msg.MinUnitAlias) || !IsBeginWithAlpha(msg.MinUnitAlias)) {
		return ErrInvalidAssetMinUnitAlias(DefaultCodespace, fmt.Sprintf("invalid token min_unit_alias %s, only accepts alphanumeric characters, and begin with an english letter, length [%d, %d]", msg.MinUnitAlias, MinimumAssetMinUnitAliasSize, MaximumAssetMinUnitAliasSize))
	}

	if msg.InitialSupply > MaximumAssetInitSupply {
		return ErrInvalidAssetInitSupply(DefaultCodespace, fmt.Sprintf("invalid token initial supply %d, only accepts value [0, %d]", msg.InitialSupply, MaximumAssetInitSupply))
	}

	if msg.MaxSupply < msg.InitialSupply || msg.MaxSupply > MaximumAssetMaxSupply {
		return ErrInvalidAssetMaxSupply(DefaultCodespace, fmt.Sprintf("invalid token max supply %d, only accepts value [%d, %d]", msg.MaxSupply, msg.InitialSupply, MaximumAssetMaxSupply))
	}

	if msg.Decimal > MaximumAssetDecimal {
		return ErrInvalidAssetDecimal(DefaultCodespace, fmt.Sprintf("invalid token decimal %d, only accepts value [0, %d]", msg.Decimal, MaximumAssetDecimal))
	}
	return nil
}

// ValidateBasic Implements Msg.
func (msg MsgIssueToken) ValidateBasic() sdk.Error {
	if msg.Source == EXTERNAL {
		return ErrInvalidAssetSource(DefaultCodespace, fmt.Sprintf("invalid source type %s", msg.Source.String()))
	}
	return ValidateMsgIssueToken(&msg)
}

// GetSignBytes Implements Msg.
func (msg MsgIssueToken) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgIssueToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgTransferTokenOwner for transferring the token owner
type MsgTransferTokenOwner struct {
	SrcOwner sdk.AccAddress `json:"src_owner"` // the current owner address of the token
	DstOwner sdk.AccAddress `json:"dst_owner"` // the new owner
	TokenID  string         `json:"token_id"`
}

// NewMsgTransferTokenOwner - construct token transfer msg
func NewMsgTransferTokenOwner(srcOwner, dstOwner sdk.AccAddress, tokenID string) MsgTransferTokenOwner {
	tokenID = strings.TrimSpace(tokenID)
	return MsgTransferTokenOwner{
		SrcOwner: srcOwner,
		DstOwner: dstOwner,
		TokenID:  tokenID,
	}
}

// GetSignBytes implements Msg
func (msg MsgTransferTokenOwner) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgTransferTokenOwner) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.SrcOwner}
}

func (msg MsgTransferTokenOwner) ValidateBasic() sdk.Error {
	// check the SrcOwner
	if len(msg.SrcOwner) == 0 {
		return ErrInvalidAddress(DefaultCodespace, fmt.Sprintf("the owner of the token must be specified"))
	}

	// check if the `DstOwner` is empty
	if len(msg.DstOwner) == 0 {
		return ErrInvalidAddress(DefaultCodespace, fmt.Sprintf("the new owner of the token must be specified"))
	}

	// check if the `DstOwner` is same as the original owner
	if msg.SrcOwner.Equals(msg.DstOwner) {
		return ErrInvalidToAddress(DefaultCodespace, fmt.Sprintf("the new owner must not be same as the original owner"))
	}

	// check the tokenId
	if err := CheckTokenID(msg.TokenID); err != nil {
		return err
	}

	return nil
}

// Route implements Msg
func (msg MsgTransferTokenOwner) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgTransferTokenOwner) Type() string { return "transfer_token_owner" }

// MsgEditToken for editing a specified token
type MsgEditToken struct {
	TokenID         string         `json:"token_id"`         //  id of token
	Owner           sdk.AccAddress `json:"owner"`            //  owner of token
	CanonicalSymbol string         `json:"canonical_symbol"` //  canonical_symbol of token
	MinUnitAlias    string         `json:"min_unit_alias"`   //  min_unit_alias of token
	MaxSupply       uint64         `json:"max_supply"`
	Mintable        Bool           `json:"mintable"` //  mintable of token
	Name            string         `json:"name"`
}

// NewMsgEditToken creates a MsgEditToken
func NewMsgEditToken(name, canonicalSymbol, minUnitAlias, tokenID string, maxSupply uint64, mintable Bool, owner sdk.AccAddress) MsgEditToken {
	name = strings.TrimSpace(name)
	canonicalSymbol = strings.ToLower(strings.TrimSpace(canonicalSymbol))
	minUnitAlias = strings.ToLower(strings.TrimSpace(minUnitAlias))
	return MsgEditToken{
		Name:            name,
		CanonicalSymbol: canonicalSymbol,
		MinUnitAlias:    minUnitAlias,
		TokenID:         tokenID,
		MaxSupply:       maxSupply,
		Mintable:        mintable,
		Owner:           owner,
	}
}

// Route implements Msg
func (msg MsgEditToken) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgEditToken) Type() string { return "edit_token" }

// ValidateBasic implements Msg
func (msg MsgEditToken) ValidateBasic() sdk.Error {

	//check owner
	if msg.Owner.Empty() {
		return ErrNilAssetOwner(DefaultCodespace, "the owner of the asset must be specified")
	}

	nameLen := len(msg.Name)
	if DoNotModify != msg.Name && nameLen > MaximumAssetNameSize {
		return ErrInvalidAssetName(DefaultCodespace, fmt.Sprintf("invalid token name %s, only accepts length (0, %d]", msg.Name, MaximumAssetNameSize))
	}

	//check max_supply for fast failed
	if msg.MaxSupply > MaximumAssetMaxSupply {
		return ErrInvalidAssetMaxSupply(DefaultCodespace, fmt.Sprintf("invalid token max supply %d, must be less than %d", msg.MaxSupply, MaximumAssetMaxSupply))
	}

	//check token_id
	if err := CheckTokenID(msg.TokenID); err != nil {
		return err
	}

	//check canonical_symbol
	canonicalSymbolLen := len(msg.CanonicalSymbol)
	if DoNotModify != msg.CanonicalSymbol && (canonicalSymbolLen < MinimumAssetSymbolSize || canonicalSymbolLen > MaximumAssetSymbolSize || !IsAlphaNumeric(msg.CanonicalSymbol)) {
		return ErrInvalidAssetCanonicalSymbol(DefaultCodespace, fmt.Sprintf("invalid token canonical_symbol %s, only accepts alphanumeric characters, length [%d, %d]", msg.CanonicalSymbol, MinimumAssetSymbolSize, MaximumAssetSymbolSize))
	}

	//check min_unit_alias
	minUnitAliasLen := len(msg.MinUnitAlias)
	if DoNotModify != msg.MinUnitAlias && (minUnitAliasLen < MinimumAssetMinUnitAliasSize || minUnitAliasLen > MaximumAssetMinUnitAliasSize || !IsAlphaNumeric(msg.MinUnitAlias) || !IsBeginWithAlpha(msg.MinUnitAlias)) {
		return ErrInvalidAssetMinUnitAlias(DefaultCodespace, fmt.Sprintf("invalid token min_unit_alias %s, only accepts alphanumeric characters, and begin with an english letter, length [%d, %d]", msg.MinUnitAlias, MinimumAssetMinUnitAliasSize, MaximumAssetMinUnitAliasSize))
	}

	return nil
}

// GetSignBytes implements Msg
func (msg MsgEditToken) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgEditToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgMintToken for mint the token to a specified address
type MsgMintToken struct {
	TokenID string         `json:"token_id"` // the unique id of the token
	Owner   sdk.AccAddress `json:"owner"`    // the current owner address of the token
	To      sdk.AccAddress `json:"to"`       // address of mint token to
	Amount  uint64         `json:"amount"`   // amount of mint token
}

// NewMsgMintToken creates a MsgMintToken
func NewMsgMintToken(tokenID string, owner, to sdk.AccAddress, amount uint64) MsgMintToken {
	tokenID = strings.TrimSpace(tokenID)
	return MsgMintToken{
		TokenID: tokenID,
		Owner:   owner,
		To:      to,
		Amount:  amount,
	}
}

// Route implements Msg
func (msg MsgMintToken) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgMintToken) Type() string { return "mint_token" }

// GetSignBytes implements Msg
func (msg MsgMintToken) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgMintToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// ValidateBasic implements Msg
func (msg MsgMintToken) ValidateBasic() sdk.Error {
	// check the owner
	if len(msg.Owner) == 0 {
		return ErrInvalidAddress(DefaultCodespace, fmt.Sprintf("the owner of the token must be specified"))
	}

	if msg.Amount <= 0 || msg.Amount > MaximumAssetMaxSupply {
		return ErrInvalidAssetMaxSupply(DefaultCodespace, fmt.Sprintf("invalid token amount %d, only accepts value (0, %d]", msg.Amount, MaximumAssetMaxSupply))
	}

	return CheckTokenID(msg.TokenID)
}

// ValidateMoniker checks if the specified moniker is valid
func ValidateMoniker(moniker string) sdk.Error {
	// check the moniker size
	if len(moniker) < MinimumGatewayMonikerSize || len(moniker) > MaximumGatewayMonikerSize {
		return ErrInvalidMoniker(DefaultCodespace, fmt.Sprintf("the length of the moniker must be between [%d,%d]", MinimumGatewayMonikerSize, MaximumGatewayMonikerSize))
	}

	// check the moniker format
	if !IsBeginWithAlpha(moniker) || !IsAlphaNumeric(moniker) {
		return ErrInvalidMoniker(DefaultCodespace, fmt.Sprintf("the moniker must begin with a letter followed by alphanumeric characters"))
	}

	// check if the moniker contains the native token name
	//if strings.Contains(strings.ToLower(moniker), sdk.Iris) {
	//	return ErrInvalidMoniker(DefaultCodespace, fmt.Sprintf("the moniker must not contain the native token name"))
	//}

	return nil
}
