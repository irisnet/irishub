package asset

import (
	sdk "github.com/irisnet/irishub/types"
)

// Gateway represents the gateway
type Gateway struct {
	ID         uint8            `json:"id"`             //  ID of the gateway
	Owner      sdk.AccAddress   `json:"owner"`          //  Owner address of the gateway
	Identity   string           `json:"identity"`       //  Identity of the gateway
	Moniker    string           `json:"moniker"`        //  Moniker of the gateway
	Details    string           `json:"details"`        //  Details of the gateway
	Website    string           `json:"website"`        //  Website of the gateway
	RedeemAddr sdk.AccAddress   `json:"redeem_address"` //  Redeem address of the gateway
	Operators  []sdk.AccAddress `json:"operators"`      //  Operators approved by the gateway
}
