package asset

import (
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

var (
	KeyNextGatewayID = []byte("newGatewayID") // key for the next gateway ID
)

// KeyAsset returns the key of the specified asset symbol
func KeyAsset(symbol string) []byte {
	return []byte(fmt.Sprintf("assets:%s", symbol))
}

// KeyGateway returns the key of the specified moniker
func KeyGateway(moniker string) []byte {
	return []byte(fmt.Sprintf("gateways:%s", moniker))
}

// KeyOwnerGateway returns the key of the specifed owner and moniker. Intended for querying all gateways of an owner
func KeyOwnerGateway(owner sdk.AccAddress, moniker string) []byte {
	return []byte(fmt.Sprintf("gateways:%d:%s", owner, moniker))
}

// KeyGatewaysSubspace returns the key prefix for iterating on all gateways of an owner
func KeyGatewaysSubspace(owner sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("gateways:%d:", owner))
}
