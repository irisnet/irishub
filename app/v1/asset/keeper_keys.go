package asset

import (
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

var (
	KeyNextGatewayID = []byte("newGatewayID") // key for the next gateway ID
	KeyGatewayCount  = []byte("gatewayCount") // key for the total number of all the gateways
)

// KeyGateway returns the key of the specified gateway id
func KeyGateway(gatewayID uint64) []byte {
	return []byte(fmt.Sprintf("gateways:%d", gatewayID))
}

// KeyOwnerGatewayID returns the key of the specifed gateway owner and ID. Intended for querying all gateway ids of a owner
func KeyOwnerGatewayID(owner sdk.AccAddress, gatewayID uint64) []byte {
	return []byte(fmt.Sprintf("gateways:%d:%d", owner, gatewayID))
}

// KeyMoniker returns the key of the specified gateway moniker
func KeyMoniker(moniker string) []byte {
	return []byte(fmt.Sprintf("gateways:%s", moniker))
}

// KeyOwnerGatewayCount returns the key which means the count of the gateways with the specified owner
func KeyOwnerGatewayCount(owner sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("ownerGatewayCount:%d", owner))
}
