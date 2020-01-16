package keeper

import (
	"fmt"
	"strings"

	"github.com/irisnet/irishub/app/v3/asset/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

var (
	prefixGateway      = []byte("gateways:")
	prefixOwnerGateway = []byte("ownerGateways:")
)

//Init Initialize module parameters during network upgrade
func (k Keeper) Init(ctx sdk.Context) {
	logger := k.Logger(ctx).With("handler", "Init")

	logger.Info("Begin execute upgrade method")
	store := ctx.KVStore(k.storeKey)
	// delete gateway
	k.iterateGateways(ctx, prefixGateway, func(key []byte) {
		logger.Info("Delete gateway information", "key", string(key))
		store.Delete(key)
	})

	// delete gateway owner
	k.iterateGateways(ctx, prefixOwnerGateway, func(key []byte) {
		logger.Info("Delete gateway owner", "key", string(key))
		store.Delete(key)
	})

	// delete Gateway/External token
	k.IterateTokens(ctx, func(token types.FungibleToken) (stop bool) {
		if token.Source == 0x01 || token.Source == 0x02 {
			tokenID := getTokenID(token.Source, token.GetSymbol(), token.Gateway)
			logger.Info("Delete token", "tokenID", tokenID)
			store.Delete(KeyTokens(token.Owner, tokenID))
			store.Delete(KeyTokens(sdk.AccAddress{}, tokenID))
			store.Delete(KeyToken(tokenID))
		}
		return false
	})

	//reset params
	param := k.GetParamSet(ctx)
	k.SetParamSet(ctx, param)
	logger.Info("End execute upgrade method")
}

func (k Keeper) iterateGateways(ctx sdk.Context, prefix []byte, op func(key []byte)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		op(iterator.Key())
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
