package asset

import (
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

var (
	KeyNextGatewayID = []byte("newGatewayID") // key for the next gateway ID
)

// KeyGateway returns the key of the specified gateway id
func KeyGateway(gatewayID uint8) []byte {
	return []byte(fmt.Sprintf("gateways:%d", gatewayID))
}

// KeyOwnerGatewayID returns the key of the specifed gateway owner and ID. Intended for querying all gateway ids of a owner
func KeyOwnerGatewayID(owner sdk.AccAddress, gatewayID uint8) []byte {
	return []byte(fmt.Sprintf("gateways:%d:%d", owner, gatewayID))
}

// KeyMoniker returns the key of the specified gateway moniker
func KeyMoniker(moniker string) []byte {
	return []byte(fmt.Sprintf("gateways:%s", moniker))
}
