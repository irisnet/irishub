package asset

import (
	"github.com/irisnet/irishub/app/v1/bank"
	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/guardian"
	sdk "github.com/irisnet/irishub/types"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
	ck       bank.Keeper
	gk       guardian.Keeper

	// codespace
	codespace sdk.CodespaceType
	// params subspace
	paramSpace params.Subspace
	// metrics
	metrics *Metrics
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, ck bank.Keeper, gk guardian.Keeper, codespace sdk.CodespaceType, paramSpace params.Subspace, metrics *Metrics) Keeper {
	keeper := Keeper{
		storeKey:   key,
		cdc:        cdc,
		ck:         ck,
		gk:         gk,
		codespace:  codespace,
		paramSpace: paramSpace.WithTypeTable(ParamTypeTable()),
		metrics:    metrics,
	}
	return keeper
}

// return the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}
