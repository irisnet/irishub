package upgrade

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

type Keeper struct {
	storeKey   		sdk.StoreKey
	cdc        		*wire.Codec
	coinKeeper 		bank.Keeper
}

func NewKeeper(cdc *wire.Codec, key sdk.StoreKey, ck bank.Keeper) Keeper {
	keeper := Keeper {
		storeKey:   key,
		cdc:        cdc,
		coinKeeper: ck,
	}
	return keeper
}

var (
	defaultSwichPeriod     int64 = 200
)