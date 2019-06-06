package asset

import (
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
