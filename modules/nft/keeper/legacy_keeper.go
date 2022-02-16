package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/nft/exported"
	"github.com/irisnet/irismod/modules/nft/types"
)

type LegacyKeeper struct {
	nk Keeper
}

func NewLegacyKeeper(nk Keeper) LegacyKeeper {
	return LegacyKeeper{nk}
}

func (n LegacyKeeper) IssueDenom(ctx sdk.Context,
	id, name, schema, symbol string,
	creator sdk.AccAddress,
	mintRestricted, updateRestricted bool) error {
	return n.nk.IssueDenom(ctx, id, name, schema, symbol, creator, mintRestricted, updateRestricted, types.DoNotModify, types.DoNotModify, types.DoNotModify, types.DoNotModify)
}

func (n LegacyKeeper) MintNFT(ctx sdk.Context,
	denomID, tokenID, tokenNm, tokenURI, tokenData string,
	owner sdk.AccAddress) error {
	return n.nk.MintNFT(ctx, denomID, tokenID, tokenNm, tokenURI, "", tokenData, owner)
}

func (n LegacyKeeper) TransferOwner(ctx sdk.Context,
	denomID, tokenID, tokenNm, tokenURI, tokenData string,
	srcOwner, dstOwner sdk.AccAddress) error {
	return n.nk.TransferOwner(ctx, denomID, tokenID, tokenNm, tokenURI, types.DoNotModify, tokenData, srcOwner, dstOwner)
}

func (n LegacyKeeper) BurnNFT(ctx sdk.Context, denomID, tokenID string, owner sdk.AccAddress) error {
	return n.nk.BurnNFT(ctx, denomID, tokenID, owner)
}

func (n LegacyKeeper) GetNFT(ctx sdk.Context, denomID, tokenID string) (nft exported.NFT, err error) {
	return n.nk.GetNFT(ctx, denomID, tokenID)
}

func (n LegacyKeeper) GetDenom(ctx sdk.Context, id string) (denom types.Denom, found bool) {
	return n.nk.GetDenom(ctx, id)
}
