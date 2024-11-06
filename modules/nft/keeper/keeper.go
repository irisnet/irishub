package keeper

import (
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/x/nft"
	nftkeeper "cosmossdk.io/x/nft/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"mods.irisnet.org/modules/nft/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeService store.KVStoreService
	cdc          codec.Codec
	nk           nftkeeper.Keeper
}

// NewKeeper creates a new instance of the NFT Keeper
func NewKeeper(cdc codec.Codec,
	storeService store.KVStoreService,
	ak nft.AccountKeeper,
	bk nft.BankKeeper,
) Keeper {
	return Keeper{
		storeService: storeService,
		cdc:          cdc,
		nk:           nftkeeper.NewKeeper(storeService, cdc, ak, bk),
	}
}

// NFTkeeper returns a cosmos-sdk nftkeeper.Keeper.
func (k Keeper) NFTkeeper() nftkeeper.Keeper {
	return k.nk
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("irismod/%s", types.ModuleName))
}
