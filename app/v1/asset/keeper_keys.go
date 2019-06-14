package asset

import (
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

var (
	PrefixGateway = []byte("gateways:") // prefix for the gateway store
)

// KeyAsset returns the key of the specified asset source and id
func KeyAsset(id string) []byte {
	return []byte(fmt.Sprintf("assets:%s", id))
}

// KeyGateway returns the key of the specified moniker
func KeyGateway(moniker string) []byte {
	return []byte(fmt.Sprintf("gateways:%s", moniker))
}

// KeyOwnerGateway returns the key of the specifed owner and moniker. Intended for querying all gateways of an owner
func KeyOwnerGateway(owner sdk.AccAddress, moniker string) []byte {
	return []byte(fmt.Sprintf("ownerGateways:%d:%s", owner, moniker))
}

// KeyGatewaysSubspace returns the key prefix for iterating on all gateways of an owner
func KeyGatewaysSubspace(owner sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("ownerGateways:%d:", owner))
}
