package asset

import (
	"fmt"
	"math"
	"regexp"

	sdk "github.com/irisnet/irishub/types"
)

const (
	// MsgRoute identifies transaction types
	MsgRoute = "asset"
)

var (
	MaximumGatewayMonikerSize = uint32(8)   // limitation for the length of the gateway's moniker
	MaximumGatewayDetailsSize = uint32(280) // limitation for the length of the gateway's details
	MaximumGatewayWebsiteSize = uint32(128) // limitation for the length of the gateway's website

	MaximumAssetInitSupply  = uint64(1e+12)
	MaximumAssetTotalSupply = math.MaxUint64
	MaximumAssetDecimal     = uint8(18)

	// 00 - fungible; 01 - non-fungible
	MsgIssueFamily = map[string]bool{"00": true, "01": true}
	// Reserved - 00 (native); 01 (external); Gateway IDs
	MsgIssueSource = map[string]bool{"00": true, "01": true}
)

var _, _, _ sdk.Msg = &MsgCreateGateway{}, &MsgEditGateway{}, &MsgIssueAsset{}

// MsgCreateGateway for creating the gateway
type MsgCreateGateway struct {
	Identity   string           `json:"identity"`       //  Identity of the gateway
	Moniker    string           `json:"moniker"`        //  Moniker of the gateway
	Details    string           `json:"details"`        //  Details of the gateway
	Website    string           `json:"website"`        //  Website of the gateway
	RedeemAddr sdk.AccAddress   `json:"redeem_address"` //  Redeem address of the gateway
	Owner      sdk.AccAddress   `json:"owner"`          //  Owner address of the gateway
	Operators  []sdk.AccAddress `json:"operators"`      //  Operators approved by the gateway
}

// NewMsgCreateGateway creates a MsgCreateGateway
func NewMsgCreateGateway(identity, moniker, details, website string, redeemAddr, owner sdk.AccAddress, operators []sdk.AccAddress) MsgCreateGateway {
	return MsgCreateGateway{
		Identity:   identity,
		Moniker:    moniker,
		Details:    details,
		Website:    website,
		RedeemAddr: redeemAddr,
		Owner:      owner,
		Operators:  operators,
	}
}

// Route implements Msg
func (msg MsgCreateGateway) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgCreateGateway) Type() string { return "create_gateway" }

// ValidateBasic implements Msg
func (msg MsgCreateGateway) ValidateBasic() sdk.Error {
	// check the moniker
	if len(msg.Moniker) == 0 || uint32(len(msg.Moniker)) > MaximumGatewayMonikerSize {
		return ErrInvalidMoniker(DefaultCodespace, fmt.Sprintf("the length of the moniker must be (0,%d]", MaximumGatewayMonikerSize))
	}

	// check the details
	if uint32(len(msg.Details)) > MaximumGatewayDetailsSize {
		return ErrInvalidDetails(DefaultCodespace, fmt.Sprintf("the length of the details must be [0,%d]", MaximumGatewayDetailsSize))
	}

	// check the website
	if uint32(len(msg.Website)) > MaximumGatewayWebsiteSize {
		return ErrInvalidDetails(DefaultCodespace, fmt.Sprintf("the length of the website must be [0,%d]", MaximumGatewayWebsiteSize))
	}

	// check if the owner is included in operators
	for _, op := range msg.Operators {
		if op.Equals(msg.Owner) {
			return ErrInvalidOperator(DefaultCodespace, "the owner can not be an operator")
		}
	}

	return nil
}

// String returns the representation of the msg
func (msg MsgCreateGateway) String() string {
	return fmt.Sprintf("MsgCreateGateway{%s, %s, %s, %s, %s, %s, %v}", msg.Identity, msg.Moniker, msg.Details, msg.Website, msg.Owner, msg.RedeemAddr, msg.Operators)
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
	Owner      sdk.AccAddress   `json:"owner"`          //  Owner of the gateway
	Moniker    string           `json:"moniker"`        //  Moniker of the gateway
	Identity   string           `json:"identity"`       //  Identity of the gateway
	Details    string           `json:"details"`        //  Details of the gateway
	Website    string           `json:"website"`        //  Website of the gateway
	RedeemAddr sdk.AccAddress   `json:"redeem_address"` //  Redeem address of the gateway
	Operators  []sdk.AccAddress `json:"operators"`      //  Operators approved by the gateway
}

// NewMsgEditGateway creates a MsgEditGateway
func NewMsgEditGateway(identity, moniker, details, website string, redeemAddr, owner sdk.AccAddress, operators []sdk.AccAddress) MsgEditGateway {
	return MsgEditGateway{
		Identity:   identity,
		Moniker:    moniker,
		Details:    details,
		Website:    website,
		RedeemAddr: redeemAddr,
		Owner:      owner,
		Operators:  operators,
	}
}

// Route implements Msg
func (msg MsgEditGateway) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgEditGateway) Type() string { return "edit_gateway" }

// ValidateBasic implements Msg
func (msg MsgEditGateway) ValidateBasic() sdk.Error {
	// TODO
	return nil
}

// String returns the representation of the msg
func (msg MsgEditGateway) String() string {
	return fmt.Sprintf("MsgEditGateway{%s, %s, %s, %s, %s, %s, %v}", msg.Owner, msg.Identity, msg.Moniker, msg.Details, msg.Website, msg.RedeemAddr, msg.Operators)
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

// ---------------------------------------------------------
// MsgIssueAsset
type MsgIssueAsset struct {
	Family     string           `json:"family"`
	Name       string           `json:"name"`
	Symbol     string           `json:"symbol"`
	Source     string           `json:"source"`
	InitSupply uint64           `json:"init_supply"`
	MaxSupply  uint64           `json:"max_supply"`
	Decimal    uint8            `json:"decimal"`
	Mintable   bool             `json:"mintable"`
	Owner      sdk.AccAddress   `json:"owner"`
	Operators  []sdk.AccAddress `json:"operators"`
}

// NewMsgIssue - construct asset issue msg.
func NewMsgIssue(family string, name string, symbol string, source string, initSupply uint64, maxSupply uint64, decimal uint8, mintable bool, owner sdk.AccAddress, operators []sdk.AccAddress) MsgIssueAsset {
	return MsgIssueAsset{Family: family, Name: name, Symbol: symbol, Source: source, InitSupply: initSupply, MaxSupply: maxSupply, Decimal: decimal, Mintable: mintable, Owner: owner, Operators: operators}
}

// Implements Msg.
// nolint
func (msg MsgIssueAsset) Route() string { return MsgRoute }
func (msg MsgIssueAsset) Type() string  { return "issue" }

// Implements Msg.
func (msg MsgIssueAsset) ValidateBasic() sdk.Error {

	// only accepts alphanumeric characters, _ and -
	reg := regexp.MustCompile(`[^a-zA-Z0-9_-]`)

	if msg.Owner == nil {
		return ErrNilAssetOwner(DefaultCodespace)
	}

	if _, found := MsgIssueFamily[msg.Family]; len(msg.Family) > 0 && !found {
		return ErrInvalidAssetFamily(DefaultCodespace, msg.Family)

	}

	if len(msg.Name) == 0 || !reg.Match([]byte(msg.Name)) {
		return ErrInvalidAssetName(DefaultCodespace, msg.Name)
	}

	if len(msg.Symbol) == 0 || !reg.Match([]byte(msg.Symbol)) {
		return ErrInvalidAssetSymbol(DefaultCodespace, msg.Symbol)
	}

	if msg.InitSupply == 0 || msg.InitSupply > MaximumAssetInitSupply {
		return ErrInvalidAssetInitSupply(DefaultCodespace, msg.InitSupply)
	}

	if msg.MaxSupply > 0 && msg.MaxSupply < msg.InitSupply {
		return ErrInvalidAssetMaxSupply(DefaultCodespace, msg.MaxSupply)
	}

	if msg.Decimal > 18 {
		return ErrInvalidAssetDecimal(DefaultCodespace, msg.Decimal)
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
// ---------------------------------------------------------