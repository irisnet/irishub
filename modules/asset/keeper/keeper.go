package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	token "github.com/irisnet/irishub/modules/asset/01-token"
)

type Keeper struct {
	TokenKeeper token.Keeper
}

// NewKeeper creates a new asset Keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace,
	codespace sdk.CodespaceType, supplyKeeper token.SupplyKeeper, feeCollectorName string,
) Keeper {
	tokenKeeper := token.NewKeeper(cdc, key, paramSpace, codespace, supplyKeeper, feeCollectorName)

	return Keeper{
		TokenKeeper: tokenKeeper,
	}
}
