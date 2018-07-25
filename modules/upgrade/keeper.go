package upgrade

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

type ModuleLifeTime struct {
	Start		int64
	End			int64
	Handler 	sdk.Handler
	store		sdk.KVStoreKey
}

type Version struct {
	Id			int		// should be equal with corresponding upgradeProposalID
	Start		int64
	ModuleList	[]ModuleLifeTime
}

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

