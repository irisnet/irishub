package iservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *wire.Codec
}

func NewKeeper(cdc *wire.Codec, key sdk.StoreKey) Keeper {
	keeper := Keeper{
		storeKey: key,
		cdc:      cdc,
	}
	return keeper
}
