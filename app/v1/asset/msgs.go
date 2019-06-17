package asset

import (
	"fmt"
	sdk "github.com/irisnet/irishub/types"
	"regexp"
	"strings"
)

const (
	// MsgRoute identifies transaction types
	MsgRoute          = "asset"
	MsgTypeIssueAsset = "issue_asset"
)

var (
	MaximumAssetMaxSupply          = uint64(1000000000000) // maximal limitation for asset max supply，1000 billion
	MaximumAssetInitSupply         = uint64(10000000000)   // maximal limitation for asset initial supply，100 billion
	MaximumAssetDecimal            = uint8(18)             // maximal limitation for asset decimal
	MinimumAssetSymbolSize         = 3                     // minimal limitation for the length of the asset's symbol / symbol_at_source
	MaximumAssetSymbolSize         = 8                     // maximal limitation for the length of the asset's symbol / symbol_at_source
	MinimumAssetSymbolMinAliasSize = 3                     // minimal limitation for the length of the asset's symbol_min_alias
	MaximumAssetSymbolMinAliasSize = 10                    // maximal limitation for the length of the asset's symbol_min_alias
	MaximumAssetNameSize           = 32                    // maximal limitation for the length of the asset's name

	MinimumGatewayMonikerSize = 3   // minimal limitation for the length of the gateway's moniker
	MaximumGatewayMonikerSize = 8   // maximal limitation for the length of the gateway's moniker
	MaximumGatewayDetailsSize = 280 // maximal limitation for the length of the gateway's details
	MaximumGatewayWebsiteSize = 128 // maximal limitation for the length of the gateway's website

	IsAlpha            = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString
	IsAlphaNumeric     = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString   // only accepts alphanumeric characters
	IsAlphaNumericDash = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString // only accepts alphanumeric characters, _ and -
	IsBeginWithAlpha   = regexp.MustCompile(`^[a-zA-Z].*`).MatchString
)

var _, _, _ sdk.Msg = &MsgIssueAsset{}, &MsgCreateGateway{}, &MsgEditGateway{}

// MsgIssueAsset
type MsgIssueAsset struct {
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

// NewMsgIssueAsset - construct asset issue msg.
func NewMsgIssueAsset(family AssetFamily, source AssetSource, gateway string, symbol string, symbolAtSource string, name string, decimal uint8, alias string, initialSupply uint64, maxSupply uint64, mintable bool, owner sdk.AccAddress, fee sdk.Coins) MsgIssueAsset {
	return MsgIssueAsset{
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
func (msg MsgIssueAsset) Route() string { return MsgRoute }
func (msg MsgIssueAsset) Type() string  { return MsgTypeIssueAsset }

// Implements Msg.
func (msg MsgIssueAsset) ValidateBasic() sdk.Error {

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
		break
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
		return ErrInvalidAssetName(DefaultCodespace, fmt.Sprintf("invalid asset name %s, only accepts length (0, %d]", msg.Name, MaximumAssetNameSize))
	}

	symbolLen := len(msg.Symbol)
	if symbolLen < MinimumAssetSymbolSize || symbolLen > MaximumAssetSymbolSize || !IsBeginWithAlpha(msg.Symbol) || !IsAlphaNumeric(msg.Symbol) {
		return ErrInvalidAssetSymbol(DefaultCodespace, fmt.Sprintf("invalid asset symbol %s, only accepts alphanumeric characters, and begin with an english letter, length [%d, %d]", msg.Symbol, MinimumAssetSymbolSize, MaximumAssetSymbolSize))
	}

	if strings.Contains(strings.ToLower(msg.Symbol), sdk.NativeTokenName) {
		return ErrInvalidAssetSymbol(DefaultCodespace, fmt.Sprintf("invalid asset symbol %s, cat not contain native token symbol %s", msg.Symbol, sdk.NativeTokenName))
	}

	symbolAtSourceLen := len(msg.SymbolAtSource)
	if symbolAtSourceLen > 0 && (symbolAtSourceLen < MinimumAssetSymbolSize || symbolAtSourceLen > MaximumAssetSymbolSize || !IsAlphaNumeric(msg.SymbolAtSource)) {
		return ErrInvalidAssetSymbolAtSource(DefaultCodespace, fmt.Sprintf("invalid asset symbol_at_source %s, only accepts alphanumeric characters, length [%d, %d]", msg.SymbolAtSource, MinimumAssetSymbolSize, MaximumAssetSymbolSize))
	}

	symbolMinAliasLen := len(msg.SymbolMinAlias)
	if symbolMinAliasLen > 0 && (symbolMinAliasLen < MinimumAssetSymbolMinAliasSize || symbolMinAliasLen > MaximumAssetSymbolMinAliasSize || !IsAlphaNumeric(msg.SymbolMinAlias)) {
		return ErrInvalidAssetSymbolMinAlias(DefaultCodespace, fmt.Sprintf("invalid asset symbol_min_alias %s, only accepts alphanumeric characters, length [%d, %d]", msg.SymbolMinAlias, MinimumAssetSymbolMinAliasSize, MaximumAssetSymbolMinAliasSize))
	}

	if msg.InitialSupply > MaximumAssetInitSupply {
		return ErrInvalidAssetInitSupply(DefaultCodespace, fmt.Sprintf("invalid asset initial supply %d, only accepts value [0, %d]", msg.InitialSupply, MaximumAssetInitSupply))
	}

	if msg.MaxSupply < msg.InitialSupply || msg.MaxSupply > MaximumAssetMaxSupply {
		return ErrInvalidAssetMaxSupply(DefaultCodespace, fmt.Sprintf("invalid asset max supply %d, only accepts value [%d, %d]", msg.MaxSupply, msg.InitialSupply, MaximumAssetMaxSupply))
	}

	if msg.Decimal > MaximumAssetDecimal {
		return ErrInvalidAssetDecimal(DefaultCodespace, fmt.Sprintf("invalid asset decimal %d, only accepts value [0, %d]", msg.Decimal, MaximumAssetDecimal))
	}

	return nil
}

// Implements Msg.
func (msg MsgIssueAsset) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgIssueAsset) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgCreateGateway for creating a gateway
type MsgCreateGateway struct {
	Owner    sdk.AccAddress `json:"owner"`    //  the owner address of the gateway
	Moniker  string         `json:"moniker"`  //  the globally unique name of the gateway
	Identity string         `json:"identity"` //  the identity of the gateway
	Details  string         `json:"details"`  //  the description of the gateway
	Website  string         `json:"website"`  //  the external website of the gateway
	Fee      sdk.Coin       `json:"fee"`      //  the fee for gateway creation
}

// NewMsgCreateGateway creates a MsgCreateGateway
func NewMsgCreateGateway(owner sdk.AccAddress, moniker, identity, details, website string, fee sdk.Coin) MsgCreateGateway {
	return MsgCreateGateway{
		Owner:    owner,
		Moniker:  moniker,
		Identity: identity,
		Details:  details,
		Website:  website,
		Fee:      fee,
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

	// check the moniker size
	if len(msg.Moniker) < MinimumGatewayMonikerSize || len(msg.Moniker) > MaximumGatewayMonikerSize {
		return ErrInvalidMoniker(DefaultCodespace, fmt.Sprintf("the length of the moniker must be between [%d,%d]", MinimumGatewayMonikerSize, MaximumGatewayMonikerSize))
	}

	// check the moniker format
	if !IsAlpha(msg.Moniker) {
		return ErrInvalidMoniker(DefaultCodespace, fmt.Sprintf("the moniker must contain only letters"))
	}

	// check the details
	if len(msg.Details) > MaximumGatewayDetailsSize {
		return ErrInvalidDetails(DefaultCodespace, fmt.Sprintf("the length of the details must be between [0,%d]", MaximumGatewayDetailsSize))
	}

	// check the website
	if len(msg.Website) > MaximumGatewayWebsiteSize {
		return ErrInvalidDetails(DefaultCodespace, fmt.Sprintf("the length of the website must be between [0,%d]", MaximumGatewayWebsiteSize))
	}

	// check the fee
	if !msg.Fee.IsNotNegative() {
		return ErrNegativeFee(DefaultCodespace, "the fee must not be negative")
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

	// check the moniker size
	if len(msg.Moniker) < MinimumGatewayMonikerSize || len(msg.Moniker) > MaximumGatewayMonikerSize {
		return ErrInvalidMoniker(DefaultCodespace, fmt.Sprintf("the length of the moniker must be between [%d,%d]", MinimumGatewayMonikerSize, MaximumGatewayMonikerSize))
	}

	// check the moniker format
	if !IsAlpha(msg.Moniker) {
		return ErrInvalidMoniker(DefaultCodespace, fmt.Sprintf("the moniker must contain only letters"))
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
