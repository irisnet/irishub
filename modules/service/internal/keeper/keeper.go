package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/service/internal/types"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
	bk       types.BankKeeper
	gk       types.GuardianKeeper

	// codespace
	codespace sdk.CodespaceType
	// params subspace
	paramSpace params.Subspace
	// metrics
	metrics *types.Metrics
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, bk types.BankKeeper, gk types.GuardianKeeper, codespace sdk.CodespaceType, paramSpace params.Subspace, metrics *types.Metrics) Keeper {
	keeper := Keeper{
		storeKey:   key,
		cdc:        cdc,
		bk:         bk,
		gk:         gk,
		codespace:  codespace,
		paramSpace: paramSpace.WithKeyTable(ParamKeyTable()),
		metrics:    metrics,
	}

	return keeper
}

// return the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}
