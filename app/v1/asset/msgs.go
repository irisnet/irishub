package asset

import (
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

const (
	// name to idetify transaction types
	MsgRoute = "asset"
)

var (
	MaximumGatewayMonikerSize = uint32(8)   // limitation for the length of the gateway's moniker
	MaximumGatewayDetailsSize = uint32(280) // limitation for the length of the gateway's details
	MaximumGatewayWebsiteSize = uint32(128) // limitation for the length of the gateway's website
)

var _, _ sdk.Msg = &MsgCreateGateway{}, &MsgEditGateway{}

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
	// TODO
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
