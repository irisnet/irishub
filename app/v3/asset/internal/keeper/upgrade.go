package keeper

import (
	"fmt"
	"strings"

	"github.com/irisnet/irishub/app/v3/asset/internal/types"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

var (
	cdc           = codec.New()
	prefixGateway = []byte("gateways:") // prefix for the gateway store
)

func init() {
	cdc.RegisterConcrete(&Gateway{}, "irishub/asset/Gateway", nil)
}

// Deprecated
type Gateway struct {
	Owner    sdk.AccAddress `json:"owner"`    //  the owner address of the gateway
	Moniker  string         `json:"moniker"`  //  the globally unique name of the gateway
	Identity string         `json:"identity"` //  the identity of the gateway
	Details  string         `json:"details"`  //  the description of the gateway
	Website  string         `json:"website"`  //  the external website of the gateway
}

// keyGateway returns the key of the specified moniker
func keyGateway(moniker string) []byte {
	return []byte(fmt.Sprintf("gateways:%s", moniker))
}

// keyOwnerGateway returns the key of the specifed owner and moniker. Intended for querying all gateways of an owner
func keyOwnerGateway(owner sdk.AccAddress, moniker string) []byte {
	return []byte(fmt.Sprintf("ownerGateways:%d:%s", owner, moniker))
}

// keyGatewaysSubspace returns the key prefix for iterating on all gateways of an owner
func keyGatewaysSubspace(owner sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("ownerGateways:%d:", owner))
}

//Init Initialize module parameters during network upgrade
func (k Keeper) Init(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)

	// delete gateway
	k.iterateGateways(ctx, func(gateway Gateway) (stop bool) {
		store.Delete(keyGateway(gateway.Moniker))
		store.Delete(keyOwnerGateway(gateway.Owner, gateway.Moniker))
		store.Delete(keyGatewaysSubspace(gateway.Owner))
		return false
	})

	// delete Gateway/External token
	k.IterateTokens(ctx, func(token types.FungibleToken) (stop bool) {
		if token.Source == 0x01 || token.Source == 0x02 {
			tokenID := getTokenID(token.Source, token.GetSymbol(), token.Gateway)
			store.Delete(KeyTokens(token.Owner, tokenID))
			store.Delete(KeyTokens(sdk.AccAddress{}, tokenID))
			store.Delete(KeyToken(tokenID))
		}
		return false
	})

	//reset params
	param := k.GetParamSet(ctx)
	k.SetParamSet(ctx, param)
}

// Deprecated
func (k Keeper) iterateGateways(ctx sdk.Context, op func(gateway Gateway) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, prefixGateway)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var gateway Gateway
		if err := cdc.UnmarshalBinaryLengthPrefixed(iterator.Value(), &gateway); err != nil {
			continue
		}

		if stop := op(gateway); stop {
			break
		}
	}
}

// Deprecated
func getTokenID(source types.AssetSource, symbol string, gateway string) string {
	switch source {
	case 0x00:
		return strings.ToLower(fmt.Sprintf("i.%s", symbol))
	case 0x01:
		return strings.ToLower(fmt.Sprintf("x.%s", symbol))
	case 0x02:
		return strings.ToLower(fmt.Sprintf("%s.%s", gateway, symbol))
	}
	return ""
}
