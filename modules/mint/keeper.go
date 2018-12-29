package mint

import (
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/params"
	"github.com/irisnet/irishub/modules/bank"
)

// keeper of the stake store
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
	paramSpace params.Subspace
	bk         bank.Keeper
	fck        FeeCollectionKeeper
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey,
	paramSpace params.Subspace, bk bank.Keeper, fck FeeCollectionKeeper) Keeper {

	keeper := Keeper{
		storeKey:   key,
		cdc:        cdc,
		paramSpace: paramSpace.WithTypeTable(ParamTypeTable()),
		bk:         bk,
		fck:        fck,
	}
	return keeper
}

//____________________________________________________________________
// Keys

var (
	minterKey = []byte{0x00} // the one key to use for the keeper store
)

//______________________________________________________________________

// get the minter
func (k Keeper) GetMinter(ctx sdk.Context) (minter Minter) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(minterKey)
	if b == nil {
		panic("Stored fee pool should not have been nil")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &minter)
	return
}

// set the minter
func (k Keeper) SetMinter(ctx sdk.Context, minter Minter) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(minter)
	store.Set(minterKey, b)
}
