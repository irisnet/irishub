package keeper

import (
	"fmt"

	"github.com/irisnet/irishub/app/v3/asset/internal/types"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

var (
	cdc           = codec.New()
	PrefixGateway = []byte("gateways:") // prefix for the gateway store
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

// Deprecated
// KeyGateway returns the key of the specified moniker
func KeyGateway(moniker string) []byte {
	return []byte(fmt.Sprintf("gateways:%s", moniker))
}

// Deprecated
// KeyOwnerGateway returns the key of the specifed owner and moniker. Intended for querying all gateways of an owner
func KeyOwnerGateway(owner sdk.AccAddress, moniker string) []byte {
	return []byte(fmt.Sprintf("ownerGateways:%d:%s", owner, moniker))
}

// Deprecated
// KeyGatewaysSubspace returns the key prefix for iterating on all gateways of an owner
func KeyGatewaysSubspace(owner sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("ownerGateways:%d:", owner))
}

//Init Initialize module parameters during network upgrade
func (k Keeper) Init(ctx sdk.Context) error {
	store := ctx.KVStore(k.storeKey)

	// delete gateway
	k.IterateGateways(ctx, func(gateway Gateway) (stop bool) {
		store.Delete(KeyGateway(gateway.Moniker))
		store.Delete(KeyOwnerGateway(gateway.Owner, gateway.Moniker))
		store.Delete(KeyGatewaysSubspace(gateway.Owner))
		return false
	})

	// delete Gateway/External token
	k.IterateTokens(ctx, func(token types.FungibleToken) (stop bool) {
		if token.Family == 0x01 || token.Family == 0x02 {
			tokenID := token.GetSymbol()
			store.Delete(KeyTokens(token.Owner, tokenID))
			store.Delete(KeyTokens(sdk.AccAddress{}, tokenID))
			store.Delete(KeyToken(tokenID))
		}
		return false
	})

	//reset params
	param := k.GetParamSet(ctx)
	k.SetParamSet(ctx, param)
	return nil
}

// Deprecated
func (k Keeper) IterateGateways(ctx sdk.Context, op func(gateway Gateway) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, PrefixGateway)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var gateway Gateway
		cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &gateway)

		if stop := op(gateway); stop {
			break
		}
	}
}
