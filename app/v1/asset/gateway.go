package asset

import (
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

// Gateway represents a gateway
type Gateway struct {
	Owner    sdk.AccAddress `json:"owner"`    //  the owner address of the gateway
	Moniker  string         `json:"moniker"`  //  the globally unique name of the gateway
	Identity string         `json:"identity"` //  the identity of the gateway
	Details  string         `json:"details"`  //  the description of the gateway
	Website  string         `json:"website"`  //  the external website of the gateway
}

// String implements fmt.Stringer
func (g Gateway) String() string {
	return fmt.Sprintf("Gateway{%s, %s, %s, %s, %s}", g.Owner, g.Moniker, g.Identity, g.Details, g.Website)
}

// Gateways is a set of gateways
type Gateways []Gateway

// String implements fmt.Stringer
func (gs Gateways) String() string {
	if len(gs) == 0 {
		return "[]"
	}

	str := fmt.Sprintf("Gateways for owner %s:", gs[0].Owner)
	for _, g := range gs {
		str += fmt.Sprintf("\n  %s: %s: %s : %s", g.Moniker, g.Identity, g.Details, g.Website)
	}
	return str
}
