package asset

import (
	"fmt"
	"regexp"

	sdk "github.com/irisnet/irishub/types"
)

const (
	// MsgRoute identifies transaction types
	MsgRoute = "asset"
)

var (
	MinimumGatewayMonikerSize = 3   // minimal limitation for the length of the gateway's moniker
	MaximumGatewayMonikerSize = 8   // maximal limitation for the length of the gateway's moniker
	MaximumGatewayDetailsSize = 280 // maximal limitation for the length of the gateway's details
	MaximumGatewayWebsiteSize = 128 // maximal limitation for the length of the gateway's website

	IsAlpha = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString
)

var _, _ sdk.Msg = &MsgCreateGateway{}, &MsgEditGateway{}

// MsgIssueAsset
type MsgIssueAsset struct {
	Asset
}

// NewMsgIssueAsset - construct asset issue msg.
func NewMsgIssueAsset(asset Asset) MsgIssueAsset {
	return MsgIssueAsset{asset}
}

// Implements Msg.
func (msg MsgIssueAsset) Route() string { return MsgRoute }
func (msg MsgIssueAsset) Type() string  { return "issue_asset" }

// Implements Msg.
func (msg MsgIssueAsset) ValidateBasic() sdk.Error {
	// only accepts alphanumeric characters, _ and -
	reg := regexp.MustCompile(`[^a-zA-Z0-9_-]`)

	baseAsset := msg.Asset.(BaseAsset)

	if baseAsset.Owner == nil {
		return ErrNilAssetOwner(DefaultCodespace)
	}

	if _, found := AssetFamilyToStringMap[baseAsset.Family]; !found {
		return ErrInvalidAssetFamily(DefaultCodespace, byte(baseAsset.Family))
	}

	if _, found := AssetSourceToStringMap[baseAsset.Source]; !found {
		return ErrInvalidAssetSource(DefaultCodespace, byte(baseAsset.Source))
	}

	if len(baseAsset.Name) == 0 || reg.Match([]byte(baseAsset.Name)) {
		return ErrInvalidAssetName(DefaultCodespace, baseAsset.Name)
	}

	if len(baseAsset.Symbol) == 0 || reg.Match([]byte(baseAsset.Symbol)) {
		return ErrInvalidAssetSymbol(DefaultCodespace, baseAsset.Symbol)
	}

	if baseAsset.InitialSupply == 0 {
		return ErrInvalidAssetInitSupply(DefaultCodespace, baseAsset.InitialSupply)
	}

	if baseAsset.MaxSupply < baseAsset.InitialSupply {
		return ErrInvalidAssetMaxSupply(DefaultCodespace, baseAsset.MaxSupply)
	}

	if baseAsset.Decimal > 18 {
		return ErrInvalidAssetDecimal(DefaultCodespace, baseAsset.Decimal)
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
	return []sdk.AccAddress{msg.Asset.(BaseAsset).Owner}
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

	// check the moniker size
	if len(msg.Moniker) < MinimumGatewayMonikerSize || len(msg.Moniker) > MaximumGatewayMonikerSize {
		return ErrInvalidMoniker(DefaultCodespace, fmt.Sprintf("the length of the moniker must be [%d,%d]", MinimumGatewayMonikerSize, MaximumGatewayMonikerSize))
	}

	// check the moniker format
	if !IsAlpha(msg.Moniker) {
		return ErrInvalidMoniker(DefaultCodespace, fmt.Sprintf("the moniker must contain only letters"))
	}

	// check the details
	if len(msg.Details) > MaximumGatewayDetailsSize {
		return ErrInvalidDetails(DefaultCodespace, fmt.Sprintf("the length of the details must be [0,%d]", MaximumGatewayDetailsSize))
	}

	// check the website
	if len(msg.Website) > MaximumGatewayWebsiteSize {
		return ErrInvalidDetails(DefaultCodespace, fmt.Sprintf("the length of the website must be [0,%d]", MaximumGatewayWebsiteSize))
	}

	return nil
}

// String returns the representation of the msg
func (msg MsgCreateGateway) String() string {
	return fmt.Sprintf("MsgCreateGateway{%s, %s, %s, %s, %s}", msg.Owner, msg.Moniker, msg.Identity, msg.Details, msg.Website)
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
		return ErrInvalidMoniker(DefaultCodespace, fmt.Sprintf("the length of the moniker must be [%d,%d]", MinimumGatewayMonikerSize, MaximumGatewayMonikerSize))
	}

	// check the moniker format
	if !IsAlpha(msg.Moniker) {
		return ErrInvalidMoniker(DefaultCodespace, fmt.Sprintf("the moniker must contain only letters"))
	}

	// check the details
	if msg.Details != nil && len(*msg.Details) > MaximumGatewayDetailsSize {
		return ErrInvalidDetails(DefaultCodespace, fmt.Sprintf("the length of the details must be [0,%d]", MaximumGatewayDetailsSize))
	}

	// check the website
	if msg.Website != nil && len(*msg.Website) > MaximumGatewayWebsiteSize {
		return ErrInvalidWebsite(DefaultCodespace, fmt.Sprintf("the length of the website must be [0,%d]", MaximumGatewayWebsiteSize))
	}

	// check if updates occur
	if msg.Identity == nil && msg.Details == nil && msg.Website == nil {
		return ErrNoUpdatesProvided(DefaultCodespace, fmt.Sprintf("no updated values provided"))
	}

	return nil
}

// String returns the representation of the msg
func (msg MsgEditGateway) String() string {
	return fmt.Sprintf("MsgEditGateway{%s, %s, %s, %s, %s}", msg.Owner, msg.Moniker, *msg.Identity, *msg.Details, *msg.Website)
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
