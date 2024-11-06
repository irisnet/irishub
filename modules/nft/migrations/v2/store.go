package v2

import (
	"time"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"mods.irisnet.org/modules/nft/types"
)

// Migrate is used to migrate nft data from irismod/nft to x/nft
func Migrate(ctx sdk.Context,
	storeService store.KVStoreService,
	cdc codec.BinaryCodec,
	logger log.Logger,
	saveDenom SaveDenom,
) error {
	logger.Info("migrate store data from version 1 to 2")
	startTime := time.Now()

	store := runtime.KVStoreAdapter(storeService.OpenKVStore(ctx))

	iterator := storetypes.KVStorePrefixIterator(store, KeyDenom(""))
	defer iterator.Close()

	k := keeper{
		storeService: storeService,
		cdc:          cdc,
	}

	var (
		denomNum int64
		tokenNum int64
	)
	for ; iterator.Valid(); iterator.Next() {
		var denom types.Denom
		cdc.MustUnmarshal(iterator.Value(), &denom)

		// delete unused key
		store.Delete(KeyDenom(denom.Id))
		store.Delete(KeyDenomName(denom.Name))
		store.Delete(KeyCollection(denom.Id))

		creator, err := sdk.AccAddressFromBech32(denom.Creator)
		if err != nil {
			return err
		}

		if err := saveDenom(ctx, denom.Id,
			denom.Name,
			denom.Schema,
			denom.Symbol,
			creator,
			denom.MintRestricted,
			denom.UpdateRestricted,
			denom.Description,
			denom.Uri,
			denom.UriHash,
			denom.Data,
		); err != nil {
			return err
		}

		tokenInDenom, err := migrateToken(ctx, k, logger, denom.Id)
		if err != nil {
			return err
		}
		denomNum++
		tokenNum += tokenInDenom

	}
	logger.Info("migrate store data success",
		"denomTotalNum", denomNum,
		"tokenTotalNum", tokenNum,
		"consume", time.Since(startTime).String(),
	)
	return nil
}

func migrateToken(
	ctx sdk.Context,
	k keeper,
	logger log.Logger,
	denomID string,
) (int64, error) {
	var iterator storetypes.Iterator
	defer func() {
		if iterator != nil {
			_ = iterator.Close()
		}
	}()

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	total := int64(0)
	iterator = storetypes.KVStorePrefixIterator(store, KeyNFT(denomID, ""))
	for ; iterator.Valid(); iterator.Next() {
		var baseNFT types.BaseNFT
		k.cdc.MustUnmarshal(iterator.Value(), &baseNFT)

		owner, err := sdk.AccAddressFromBech32(baseNFT.Owner)
		if err != nil {
			return 0, err
		}

		// delete unused key
		store.Delete(KeyNFT(denomID, baseNFT.Id))
		store.Delete(KeyOwner(owner, denomID, baseNFT.Id))

		if err := k.saveNFT(ctx, denomID,
			baseNFT.Id,
			baseNFT.Name,
			baseNFT.URI,
			baseNFT.UriHash,
			baseNFT.Data,
			owner,
		); err != nil {
			return 0, err
		}
		total++
	}
	logger.Info("migrate nft success", "denomID", denomID, "nftNum", total)
	return total, nil
}
