package asset

import (
	"github.com/irisnet/irishub/app/v1/bank"
	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/guardian"
	sdk "github.com/irisnet/irishub/types"
)

var (
	KeyNextGatewayID = []byte("newGatewayID") // key for the next gateway ID
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
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, ck bank.Keeper, gk guardian.Keeper, codespace sdk.CodespaceType, paramSpace params.Subspace) Keeper {
	return Keeper{
		storeKey:   key,
		cdc:        cdc,
		ck:         ck,
		gk:         gk,
		codespace:  codespace,
		paramSpace: paramSpace.WithTypeTable(ParamTypeTable()),
	}
}

// return the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

// CreateGateway creates a gateway
func (k Keeper) CreateGateway(ctx sdk.Context, msg MsgCreateGateway) {
	// TODO
}

// EditGateway edits the specified gateway by moniker
func (k Keeper) EditGateway(ctx sdk.Context, msg MsgEditGateway) {
	// TODO
}

// getGateway retrieves the gateway of the given moniker
func (k Keeper) getGateway(moniker string) {

}

// getNewGatewayID gets the next available gateway ID and increments it
func (k Keeper) getNewGatewayID(ctx sdk.Context) (gatewayID uint64, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyNextGatewayID)
	if bz == nil {
		return 0, ErrInvalidGenesis(k.codespace, "Initial gateway ID never set")
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &gatewayID)
	bz = k.cdc.MustMarshalBinaryLengthPrefixed(gatewayID + 1)
	store.Set(KeyNextGatewayID, bz)

	return gatewayID, nil
}
