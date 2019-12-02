package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Gateway represents a gateway
type Gateway struct {
	Owner    sdk.AccAddress `json:"owner"`    //  the owner address of the gateway
	Moniker  string         `json:"moniker"`  //  the globally unique name of the gateway
	Identity string         `json:"identity"` //  the identity of the gateway
	Details  string         `json:"details"`  //  the description of the gateway
	Website  string         `json:"website"`  //  the external website of the gateway
}

// NewGateway constructs a gateway
func NewGateway(owner sdk.AccAddress, moniker, identity, details, website string) Gateway {
	return Gateway{
		Owner:    owner,
		Moniker:  moniker,
		Identity: identity,
		Details:  details,
		Website:  website,
	}
}

// Validate checks if a gateway is valid
func (g Gateway) Validate() sdk.Error {
	// check the moniker
	if err := ValidateMoniker(g.Moniker); err != nil {
		return err
	}

	// check the description fields
	if err := validateGatewayDesc(&g.Identity, &g.Details, &g.Website); err != nil {
		return err
	}

	return nil
}

// String implements fmt.Stringer
func (g Gateway) String() string {
	return fmt.Sprintf(`Gateway:
  Owner:             %s
  Moniker:           %s
  Identity:          %s
  Details:           %s
  Website:           %s`,
		g.Owner, g.Moniker, g.Identity, g.Details, g.Website)
}

// Gateways is a set of gateways
type Gateways []Gateway

// String implements fmt.Stringer
func (gs Gateways) String() string {
	if len(gs) == 0 {
		return "[]"
	}

	var owners []string
	ownerToGateways := make(map[string][]Gateway)

	for _, g := range gs {
		owner := g.Owner.String()

		if _, ok := ownerToGateways[owner]; !ok {
			owners = append(owners, owner)
		}

		ownerToGateways[owner] = append(ownerToGateways[owner], g)
	}

	var str string
	for _, o := range owners {
		str += fmt.Sprintf("Gateways for owner %s:\n", o)

		for _, g := range ownerToGateways[o] {
			str += fmt.Sprintf("  Moniker: %s, Identity: %s, Details: %s, Website: %s\n", g.Moniker, g.Identity, g.Details, g.Website)
		}
	}

	return str
}
