package asset

import (
	"fmt"
	"regexp"
	"strings"

	sdk "github.com/irisnet/irishub/types"
)

const (
	// MsgRoute identifies transaction types
	MsgRoute          = "asset"
	MsgTypeIssueToken = "issue_token"
	// constant used in flags to indicate that description field should not be updated
	DoNotModifyDesc = "[do-not-modify]"
)

var (
	MaximumAssetMaxSupply          = uint64(1000000000000) // maximal limitation for asset max supply，1000 billion
	MaximumAssetInitSupply         = uint64(100000000000)  // maximal limitation for asset initial supply，100 billion
	MaximumAssetDecimal            = uint8(18)             // maximal limitation for asset decimal
	MinimumAssetSymbolSize         = 3                     // minimal limitation for the length of the asset's symbol / symbol_at_source
	MaximumAssetSymbolSize         = 8                     // maximal limitation for the length of the asset's symbol / symbol_at_source
	MinimumAssetSymbolMinAliasSize = 3                     // minimal limitation for the length of the asset's symbol_min_alias
	MaximumAssetSymbolMinAliasSize = 10                    // maximal limitation for the length of the asset's symbol_min_alias
	MaximumAssetNameSize           = 32                    // maximal limitation for the length of the asset's name

	MinimumGatewayMonikerSize  = 3   // minimal limitation for the length of the gateway's moniker
	MaximumGatewayMonikerSize  = 8   // maximal limitation for the length of the gateway's moniker
	MaximumGatewayIdentitySize = 128 // maximal limitation for the length of the gateway's identity
	MaximumGatewayDetailsSize  = 280 // maximal limitation for the length of the gateway's details
	MaximumGatewayWebsiteSize  = 128 // maximal limitation for the length of the gateway's website

	IsAlpha            = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString
	IsAlphaNumeric     = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString   // only accepts alphanumeric characters
	IsAlphaNumericDash = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString // only accepts alphanumeric characters, _ and -
	IsBeginWithAlpha   = regexp.MustCompile(`^[a-zA-Z].*`).MatchString
)

var _, _, _, _ sdk.Msg = &MsgIssueToken{}, &MsgCreateGateway{}, &MsgEditGateway{}, &MsgTransferGatewayOwner{}

// MsgIssueToken
type MsgIssueToken struct {
	Family         AssetFamily    `json:"family"`
	Source         AssetSource    `json:"source"`
	Gateway        string         `json:"gateway"`
	Symbol         string         `json:"symbol"`
	SymbolAtSource string         `json:"symbol_at_source"`
	Name           string         `json:"name"`
	Decimal        uint8          `json:"decimal"`
	SymbolMinAlias string         `json:"symbol_min_alias"`
	InitialSupply  uint64         `json:"initial_supply"`
	MaxSupply      uint64         `json:"max_supply"`
	Mintable       bool           `json:"mintable"`
	Owner          sdk.AccAddress `json:"owner"`
	Fee            sdk.Coins      `json:"fee"`
}

// NewMsgIssueToken - construct asset issue msg.
func NewMsgIssueToken(family AssetFamily, source AssetSource, gateway string, symbol string, symbolAtSource string, name string, decimal uint8, alias string, initialSupply uint64, maxSupply uint64, mintable bool, owner sdk.AccAddress, fee sdk.Coins) MsgIssueToken {
	return MsgIssueToken{
		Family:         family,
		Source:         source,
		Gateway:        gateway,
		Symbol:         symbol,
		SymbolAtSource: symbolAtSource,
		Name:           name,
		Decimal:        decimal,
		SymbolMinAlias: alias,
		InitialSupply:  initialSupply,
		MaxSupply:      maxSupply,
		Mintable:       mintable,
		Owner:          owner,
		Fee:            fee,
	}
}

// Implements Msg.
func (msg MsgIssueToken) Route() string { return MsgRoute }
func (msg MsgIssueToken) Type() string  { return MsgTypeIssueToken }

// Implements Msg.
func (msg MsgIssueToken) ValidateBasic() sdk.Error {

	msg.Gateway = strings.ToLower(strings.TrimSpace(msg.Gateway))
	msg.Symbol = strings.ToLower(strings.TrimSpace(msg.Symbol))
	msg.SymbolAtSource = strings.ToLower(strings.TrimSpace(msg.SymbolAtSource))
	msg.SymbolMinAlias = strings.ToLower(strings.TrimSpace(msg.SymbolMinAlias))
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
		// ignore SymbolAtSource for native asset
		msg.SymbolAtSource = ""

		break
	case EXTERNAL:
		return ErrInvalidAssetSource(DefaultCodespace, fmt.Sprintf("invalid source type %s", msg.Source.String()))
	case GATEWAY:
		// require gateway moniker for gateway asset
		if len(msg.Gateway) < MinimumGatewayMonikerSize || len(msg.Gateway) > MaximumGatewayMonikerSize {
			return ErrInvalidMoniker(DefaultCodespace, fmt.Sprintf("invalid gateway moniker, length [%d,%d]", MinimumGatewayMonikerSize, MaximumGatewayMonikerSize))
		}

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

	if strings.Contains(strings.ToLower(msg.Symbol), sdk.NativeTokenName) {
		return ErrInvalidAssetSymbol(DefaultCodespace, fmt.Sprintf("invalid token symbol %s, can not contain native token symbol %s", msg.Symbol, sdk.NativeTokenName))
	}

	symbolAtSourceLen := len(msg.SymbolAtSource)
	if symbolAtSourceLen > 0 && (symbolAtSourceLen < MinimumAssetSymbolSize || symbolAtSourceLen > MaximumAssetSymbolSize || !IsAlphaNumeric(msg.SymbolAtSource)) {
		return ErrInvalidAssetSymbolAtSource(DefaultCodespace, fmt.Sprintf("invalid token symbol_at_source %s, only accepts alphanumeric characters, length [%d, %d]", msg.SymbolAtSource, MinimumAssetSymbolSize, MaximumAssetSymbolSize))
	}

	symbolMinAliasLen := len(msg.SymbolMinAlias)
	if symbolMinAliasLen > 0 && (symbolMinAliasLen < MinimumAssetSymbolMinAliasSize || symbolMinAliasLen > MaximumAssetSymbolMinAliasSize || !IsAlphaNumeric(msg.SymbolMinAlias)) || !IsBeginWithAlpha(msg.SymbolMinAlias) {
		return ErrInvalidAssetSymbolMinAlias(DefaultCodespace, fmt.Sprintf("invalid token symbol_min_alias %s, only accepts alphanumeric characters, and begin with an english letter, length [%d, %d]", msg.SymbolMinAlias, MinimumAssetSymbolMinAliasSize, MaximumAssetSymbolMinAliasSize))
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

// Implements Msg.
func (msg MsgIssueToken) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgIssueToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgCreateGateway for creating a gateway
type MsgCreateGateway struct {
	Owner    sdk.AccAddress `json:"owner"`    //  the owner address of the gateway
	Moniker  string         `json:"moniker"`  //  the globally unique name of the gateway
	Identity string         `json:"identity"` //  the identity of the gateway
	Details  string         `json:"details"`  //  the description of the gateway
	Website  string         `json:"website"`  //  the external website of the gateway
}

// NewMsgCreateGateway creates a MsgCreateGateway
func NewMsgCreateGateway(owner sdk.AccAddress, moniker, identity, details, website string) MsgCreateGateway {
	return MsgCreateGateway{
		Owner:    owner,
		Moniker:  moniker,
		Identity: identity,
		Details:  details,
		Website:  website,
	}
}

// Route implements Msg
func (msg MsgCreateGateway) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgCreateGateway) Type() string { return "create_gateway" }

// ValidateBasic implements Msg
func (msg MsgCreateGateway) ValidateBasic() sdk.Error {
	// check the owner
	if len(msg.Owner) == 0 {
		return ErrInvalidAddress(DefaultCodespace, fmt.Sprintf("the owner of the gateway must be specified"))
	}

	// check the moniker
	if err := ValidateMoniker(msg.Moniker); err != nil {
		return err
	}

	// check the identity
	if len(msg.Identity) > MaximumGatewayIdentitySize {
		return ErrInvalidIdentity(DefaultCodespace, fmt.Sprintf("the length of the identity must be between [0,%d]", MaximumGatewayIdentitySize))
	}

	// check the details
	if len(msg.Details) > MaximumGatewayDetailsSize {
		return ErrInvalidDetails(DefaultCodespace, fmt.Sprintf("the length of the details must be between [0,%d]", MaximumGatewayDetailsSize))
	}

	// check the website
	if len(msg.Website) > MaximumGatewayWebsiteSize {
		return ErrInvalidWebsite(DefaultCodespace, fmt.Sprintf("the length of the website must be between [0,%d]", MaximumGatewayWebsiteSize))
	}

	return nil
}

// String returns the representation of the msg
func (msg MsgCreateGateway) String() string {
	return fmt.Sprintf(`MsgCreateGateway:
  Owner:             %s
  Moniker:           %s
  Identity:          %s
  Details:           %s
  Website:           %s`,
		msg.Owner, msg.Moniker, msg.Identity, msg.Details, msg.Website)
}

// GetSignBytes implements Msg
func (msg MsgCreateGateway) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgCreateGateway) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgEditGateway for editing a specified gateway
type MsgEditGateway struct {
	Owner    sdk.AccAddress `json:"owner"`    //  Owner of the gateway
	Moniker  string         `json:"moniker"`  //  Moniker of the gateway
	Identity *string        `json:"identity"` //  Identity of the gateway
	Details  *string        `json:"details"`  //  Details of the gateway
	Website  *string        `json:"website"`  //  Website of the gateway
}

// NewMsgEditGateway creates a MsgEditGateway
func NewMsgEditGateway(owner sdk.AccAddress, moniker string, identity, details, website *string) MsgEditGateway {
	return MsgEditGateway{
		Owner:    owner,
		Moniker:  moniker,
		Identity: identity,
		Details:  details,
		Website:  website,
	}
}

// Route implements Msg
func (msg MsgEditGateway) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgEditGateway) Type() string { return "edit_gateway" }

// ValidateBasic implements Msg
func (msg MsgEditGateway) ValidateBasic() sdk.Error {
	// check the owner
	if len(msg.Owner) == 0 {
		return ErrInvalidAddress(DefaultCodespace, fmt.Sprintf("the owner of the gateway must be specified"))
	}

	// check the moniker
	if err := ValidateMoniker(msg.Moniker); err != nil {
		return err
	}

	// check the identity
	if msg.Identity != nil && len(*msg.Identity) > MaximumGatewayIdentitySize {
		return ErrInvalidIdentity(DefaultCodespace, fmt.Sprintf("the length of the identity must be between [0,%d]", MaximumGatewayIdentitySize))
	}

	// check the details
	if msg.Details != nil && len(*msg.Details) > MaximumGatewayDetailsSize {
		return ErrInvalidDetails(DefaultCodespace, fmt.Sprintf("the length of the details must be between [0,%d]", MaximumGatewayDetailsSize))
	}

	// check the website
	if msg.Website != nil && len(*msg.Website) > MaximumGatewayWebsiteSize {
		return ErrInvalidWebsite(DefaultCodespace, fmt.Sprintf("the length of the website must be between [0,%d]", MaximumGatewayWebsiteSize))
	}

	// check if updates occur
	if msg.Identity == nil && msg.Details == nil && msg.Website == nil {
		return ErrNoUpdatesProvided(DefaultCodespace, fmt.Sprintf("no updated values provided"))
	}

	return nil
}

// String returns the representation of the msg
func (msg MsgEditGateway) String() string {
	return fmt.Sprintf(`MsgEditGateway:
  Owner:             %s
  Moniker:           %s
  Identity:          %s
  Details:           %s
  Website:           %s`,
		msg.Owner, msg.Moniker, *msg.Identity, *msg.Details, *msg.Website)
}

// GetSignBytes implements Msg
func (msg MsgEditGateway) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgEditGateway) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgTransferGatewayOwner for transferring the gateway owner
type MsgTransferGatewayOwner struct {
	Owner   sdk.AccAddress `json:"owner"`   //  the current owner address of the gateway
	Moniker string         `json:"moniker"` //  the unique name of the gateway to be transferred
	To      sdk.AccAddress `json:"to"`      // the new owner to which the gateway ownership will be transferred
}

// NewMsgTransferGatewayOwner creates a MsgTransferGatewayOwner
func NewMsgTransferGatewayOwner(owner sdk.AccAddress, moniker string, to sdk.AccAddress) MsgTransferGatewayOwner {
	return MsgTransferGatewayOwner{
		Owner:   owner,
		Moniker: moniker,
		To:      to,
	}
}

// Route implements Msg
func (msg MsgTransferGatewayOwner) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgTransferGatewayOwner) Type() string { return "transfer_gateway_owner" }

// ValidateBasic implements Msg
func (msg MsgTransferGatewayOwner) ValidateBasic() sdk.Error {
	// check the owner
	if len(msg.Owner) == 0 {
		return ErrInvalidAddress(DefaultCodespace, fmt.Sprintf("the owner of the gateway must be specified"))
	}

	// check if the `to` is empty
	if len(msg.To) == 0 {
		return ErrInvalidAddress(DefaultCodespace, fmt.Sprintf("the new owner of the gateway must be specified"))
	}

	// check if the `to` is same as the original owner
	if msg.To.Equals(msg.Owner) {
		return ErrInvalidToAddress(DefaultCodespace, fmt.Sprintf("the new owner must not be same as the original owner"))
	}

	// check the moniker
	if err := ValidateMoniker(msg.Moniker); err != nil {
		return err
	}

	return nil
}

// String returns the representation of the msg
func (msg MsgTransferGatewayOwner) String() string {
	return fmt.Sprintf(`MsgTransferGatewayOwner:
  Owner:             %s
  Moniker:           %s
  To:                %s`,
		msg.Owner, msg.Moniker, msg.To)
}

// GetSignBytes implements Msg
func (msg MsgTransferGatewayOwner) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgTransferGatewayOwner) GetSigners() []sdk.AccAddress {
	// the msg needs signatures from both the original owner and the new one
	return []sdk.AccAddress{msg.Owner, msg.To}
}

// ValidateMoniker checks if the specified moniker is valid
func ValidateMoniker(moniker string) sdk.Error {
	// check the moniker size
	if len(moniker) < MinimumGatewayMonikerSize || len(moniker) > MaximumGatewayMonikerSize {
		return ErrInvalidMoniker(DefaultCodespace, fmt.Sprintf("the length of the moniker must be between [%d,%d]", MinimumGatewayMonikerSize, MaximumGatewayMonikerSize))
	}

	// check the moniker format
	if !IsAlpha(moniker) {
		return ErrInvalidMoniker(DefaultCodespace, fmt.Sprintf("the moniker must contain only letters"))
	}

	// check if the moniker contains the native token name
	if strings.Contains(strings.ToLower(moniker), sdk.NativeTokenName) {
		return ErrInvalidMoniker(DefaultCodespace, fmt.Sprintf("the moniker must not contain the native token name"))
	}

	return nil
}

// MsgEditToken for editing a specified gateway
type MsgEditToken struct {
	TokenId        string         `json:"token_id"`         //  id of asset
	Owner          sdk.AccAddress `json:"owner"`            // owner of asset
	SymbolAtSource string         `json:"symbol_at_source"` //  symbol_at_source of asset
	SymbolMinAlias string         `json:"symbol_min_alias"` //  symbol_min_alias of asset
	MaxSupply      uint64         `json:"max_supply"`
	Mintable       *bool          `json:"mintable"` //  mintable of asset
	Name           string         `json:"name"`
}

// NewMsgEditToken creates a MsgEditAsset
func NewMsgEditToken(name, symbolAtSource, symbolMinAlias, tokenId string, maxSupply uint64, mintable *bool, owner sdk.AccAddress) MsgEditToken {
	name = strings.TrimSpace(name)
	symbolAtSource = strings.ToLower(strings.TrimSpace(symbolAtSource))
	symbolMinAlias = strings.ToLower(strings.TrimSpace(symbolMinAlias))
	return MsgEditToken{
		Name:           name,
		SymbolAtSource: symbolAtSource,
		SymbolMinAlias: symbolMinAlias,
		TokenId:        tokenId,
		MaxSupply:      maxSupply,
		Mintable:       mintable,
		Owner:          owner,
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
	if DoNotModifyDesc != msg.Name && nameLen > MaximumAssetNameSize {
		return ErrInvalidAssetName(DefaultCodespace, fmt.Sprintf("invalid token name %s, only accepts length (0, %d]", msg.Name, MaximumAssetNameSize))
	}

	//check max_supply for fast failed
	if msg.MaxSupply > MaximumAssetMaxSupply {
		return ErrInvalidAssetMaxSupply(DefaultCodespace, fmt.Sprintf("invalid token max supply %d, must be less than %d", msg.MaxSupply, MaximumAssetMaxSupply))
	}

	//check token_id
	if err := CheckAssetID(msg.TokenId); err != nil {
		return err
	}

	//check symbol_at_source
	symbolAtSourceLen := len(msg.SymbolAtSource)
	if DoNotModifyDesc != msg.SymbolAtSource && (symbolAtSourceLen < MinimumAssetSymbolSize || symbolAtSourceLen > MaximumAssetSymbolSize || !IsAlphaNumeric(msg.SymbolAtSource)) {
		return ErrInvalidAssetSymbolAtSource(DefaultCodespace, fmt.Sprintf("invalid token symbol_at_source %s, only accepts alphanumeric characters, length [%d, %d]", msg.SymbolAtSource, MinimumAssetSymbolSize, MaximumAssetSymbolSize))
	}

	//check symbol_min_alias
	symbolMinAliasLen := len(msg.SymbolMinAlias)
	if DoNotModifyDesc != msg.SymbolMinAlias && (symbolMinAliasLen < MinimumAssetSymbolMinAliasSize || symbolMinAliasLen > MaximumAssetSymbolMinAliasSize || !IsAlphaNumeric(msg.SymbolMinAlias) || !IsBeginWithAlpha(msg.SymbolMinAlias)) {
		return ErrInvalidAssetSymbolMinAlias(DefaultCodespace, fmt.Sprintf("invalid token symbol_min_alias %s, only accepts alphanumeric characters, and begin with an english letter, length [%d, %d]", msg.SymbolMinAlias, MinimumAssetSymbolMinAliasSize, MaximumAssetSymbolMinAliasSize))
	}

	return nil
}

// GetSignBytes implements Msg
func (msg MsgEditToken) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
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
	TokenId string         `json:"moniker"` // the unique id of the token
	Owner   sdk.AccAddress `json:"owner"`   // the current owner address of the token
	To      sdk.AccAddress `json:"to"`      // address of mint token to
	Amount  uint64         `json:"amount"`  // amount of mint token
}

// NewMsgMintToken creates a MsgMintToken
func NewMsgMintToken(tokenId string, owner, to sdk.AccAddress, amount uint64) MsgMintToken {
	tokenId = strings.TrimSpace(tokenId)
	return MsgMintToken{
		TokenId: tokenId,
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
	b, err := msgCdc.MarshalJSON(msg)
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
		return ErrInvalidAddress(DefaultCodespace, fmt.Sprintf("the owner of the gateway must be specified"))
	}

	if msg.Amount <= 0 || msg.Amount > MaximumAssetMaxSupply {
		return ErrInvalidAssetMaxSupply(DefaultCodespace, fmt.Sprintf("invalid token max supply %d, only accepts value (0, %d]", msg.Amount, MaximumAssetMaxSupply))
	}

	return CheckAssetID(msg.TokenId)
}
