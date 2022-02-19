package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/mt/exported"
	"github.com/irisnet/irismod/modules/mt/types"
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

func (n LegacyKeeper) MintMT(ctx sdk.Context,
	denomID, tokenID, tokenNm, tokenURI, tokenData string,
	owner sdk.AccAddress) error {
	return n.nk.MintMT(ctx, denomID, tokenID, tokenNm, tokenURI, "", tokenData, owner)
}

func (n LegacyKeeper) TransferOwner(ctx sdk.Context,
	denomID, tokenID, tokenNm, tokenURI, tokenData string,
	srcOwner, dstOwner sdk.AccAddress) error {
	return n.nk.TransferOwner(ctx, denomID, tokenID, tokenNm, tokenURI, types.DoNotModify, tokenData, srcOwner, dstOwner)
}

func (n LegacyKeeper) BurnMT(ctx sdk.Context, denomID, tokenID string, owner sdk.AccAddress) error {
	return n.nk.BurnMT(ctx, denomID, tokenID, owner)
}

func (n LegacyKeeper) GetMT(ctx sdk.Context, denomID, tokenID string) (mt exported.MT, err error) {
	return n.nk.GetMT(ctx, denomID, tokenID)
}

func (n LegacyKeeper) GetDenom(ctx sdk.Context, id string) (denom types.Denom, found bool) {
	return n.nk.GetDenom(ctx, id)
}
