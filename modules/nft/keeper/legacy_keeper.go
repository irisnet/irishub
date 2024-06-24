package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"mods.irisnet.org/nft/exported"
	"mods.irisnet.org/nft/types"
)

type LegacyKeeper struct {
	nk Keeper
}

func NewLegacyKeeper(nk Keeper) LegacyKeeper {
	return LegacyKeeper{nk}
}

func (n LegacyKeeper) IssueDenom(
	ctx sdk.Context,
	id,
	name,
	schema,
	symbol string,
	creator sdk.AccAddress,
	mintRestricted,
	updateRestricted bool,
) error {
	return n.nk.SaveDenom(
		ctx,
		id,
		name,
		schema,
		symbol,
		creator,
		mintRestricted,
		updateRestricted,
		types.DoNotModify,
		types.DoNotModify,
		types.DoNotModify,
		types.DoNotModify,
	)
}

func (n LegacyKeeper) MintNFT(
	ctx sdk.Context,
	denomID,
	tokenID,
	tokenNm,
	tokenURI,
	tokenData string,
	owner sdk.AccAddress,
) error {
	return n.nk.SaveNFT(
		ctx,
		denomID,
		tokenID,
		tokenNm,
		tokenURI,
		"",
		tokenData,
		owner,
	)
}

func (n LegacyKeeper) TransferOwner(
	ctx sdk.Context,
	denomID,
	tokenID,
	tokenNm,
	tokenURI,
	tokenData string,
	srcOwner,
	dstOwner sdk.AccAddress,
) error {
	return n.nk.TransferOwnership(
		ctx,
		denomID,
		tokenID,
		tokenNm,
		tokenURI,
		types.DoNotModify,
		tokenData,
		srcOwner,
		dstOwner,
	)
}

func (n LegacyKeeper) BurnNFT(ctx sdk.Context, denomID, tokenID string, owner sdk.AccAddress) error {
	return n.nk.RemoveNFT(ctx, denomID, tokenID, owner)
}

func (n LegacyKeeper) GetNFT(ctx sdk.Context, denomID, tokenID string) (nft exported.NFT, err error) {
	return n.nk.GetNFT(ctx, denomID, tokenID)
}

func (n LegacyKeeper) GetDenom(ctx sdk.Context, id string) (denom types.Denom, found bool) {
	d, err := n.nk.GetDenomInfo(ctx, id)
	if err != nil {
		return denom, false
	}
	return *d, true
}
