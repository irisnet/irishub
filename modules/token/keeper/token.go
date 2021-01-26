package keeper

import (
	"fmt"

	gogotypes "github.com/gogo/protobuf/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/irisnet/irismod/modules/token/types"
)

// GetTokens returns all existing tokens
func (k Keeper) GetTokens(ctx sdk.Context, owner sdk.AccAddress) (tokens []types.TokenI) {
	store := ctx.KVStore(k.storeKey)

	var it sdk.Iterator
	if owner == nil {
		it = sdk.KVStorePrefixIterator(store, types.PrefixTokenForSymbol)
		defer it.Close()

		for ; it.Valid(); it.Next() {
			var token types.Token
			k.cdc.MustUnmarshalBinaryBare(it.Value(), &token)

			tokens = append(tokens, &token)
		}
		return
	}

	it = sdk.KVStorePrefixIterator(store, types.KeyTokens(owner, ""))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		var symbol gogotypes.StringValue
		k.cdc.MustUnmarshalBinaryBare(it.Value(), &symbol)

		token, err := k.GetToken(ctx, symbol.Value)
		if err != nil {
			continue
		}
		tokens = append(tokens, token)
	}
	return
}

// GetToken returns the token of the specified symbol or minUint
func (k Keeper) GetToken(ctx sdk.Context, denom string) (types.TokenI, error) {
	store := ctx.KVStore(k.storeKey)

	if token, err := k.getToken(ctx, denom); err == nil {
		return &token, nil
	}

	bz := store.Get(types.KeyMinUint(denom))
	if bz == nil {
		return nil, sdkerrors.Wrap(types.ErrTokenNotExists, fmt.Sprintf("token %s does not exist", denom))
	}

	var symbol gogotypes.StringValue
	k.cdc.MustUnmarshalBinaryBare(bz, &symbol)

	token, err := k.getToken(ctx, symbol.Value)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

// AddToken saves a new token
func (k Keeper) AddToken(ctx sdk.Context, token types.Token) error {
	if k.HasToken(ctx, token.Symbol) {
		return sdkerrors.Wrapf(types.ErrSymbolAlreadyExists, "symbol already exists: %s", token.Symbol)
	}

	if k.HasToken(ctx, token.MinUnit) {
		return sdkerrors.Wrapf(types.ErrMinUnitAlreadyExists, "min-unit already exists: %s", token.MinUnit)
	}

	// set token
	k.setToken(ctx, token)

	// set token to be prefixed with min unit
	k.setWithMinUnit(ctx, token.MinUnit, token.Symbol)

	if len(token.Owner) != 0 {
		// set token to be prefixed with owner
		k.setWithOwner(ctx, token.GetOwner(), token.Symbol)
	}

	denomMetaData := banktypes.Metadata{
		Description: token.Name,
		Base:        token.MinUnit,
		Display:     token.Symbol,
		DenomUnits: []*banktypes.DenomUnit{
			{Denom: token.MinUnit, Exponent: 0},
			{Denom: token.Symbol, Exponent: token.Scale},
		},
	}
	k.bankKeeper.SetDenomMetaData(ctx, denomMetaData)

	// Set token to be prefixed with min_unit
	k.setWithMinUnit(ctx, token.MinUnit, token.Symbol)

	return nil
}

// HasToken asserts a token exists
func (k Keeper) HasToken(ctx sdk.Context, denom string) bool {
	store := ctx.KVStore(k.storeKey)
	if store.Has(types.KeySymbol(denom)) {
		return true
	}

	return store.Has(types.KeyMinUint(denom))
}

// GetOwner returns the owner of the specified token
func (k Keeper) GetOwner(ctx sdk.Context, denom string) (sdk.AccAddress, error) {
	token, err := k.GetToken(ctx, denom)
	if err != nil {
		return nil, err
	}

	return token.GetOwner(), nil
}

// AddBurnCoin saves the total amount of the burned tokens
func (k Keeper) AddBurnCoin(ctx sdk.Context, coin sdk.Coin) {
	var total = coin
	if hasCoin, err := k.GetBurnCoin(ctx, coin.Denom); err == nil {
		total = total.Add(hasCoin)
	}

	bz := k.cdc.MustMarshalBinaryBare(&total)
	key := types.KeyBurnTokenAmt(coin.Denom)

	store := ctx.KVStore(k.storeKey)
	store.Set(key, bz)
}

// GetBurnCoin returns the total amount of the burned tokens
func (k Keeper) GetBurnCoin(ctx sdk.Context, minUint string) (sdk.Coin, error) {
	key := types.KeyBurnTokenAmt(minUint)
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(key)

	if len(bz) == 0 {
		return sdk.Coin{}, sdkerrors.Wrapf(types.ErrNotFoundTokenAmt, "not found symbol:%s", minUint)
	}

	var coin sdk.Coin
	k.cdc.MustUnmarshalBinaryBare(bz, &coin)

	return coin, nil
}

// GetAllBurnCoin returns the total amount of all the burned tokens
func (k Keeper) GetAllBurnCoin(ctx sdk.Context) []sdk.Coin {
	store := ctx.KVStore(k.storeKey)

	var coins []sdk.Coin
	it := sdk.KVStorePrefixIterator(store, types.PeffixBurnTokenAmt)
	for ; it.Valid(); it.Next() {
		var coin sdk.Coin
		k.cdc.MustUnmarshalBinaryBare(it.Value(), &coin)
		coins = append(coins, coin)
	}

	return coins
}

// GetParamSet returns token params from the global param store
func (k Keeper) GetParamSet(ctx sdk.Context) types.Params {
	var p types.Params
	k.paramSpace.GetParamSet(ctx, &p)
	return p
}

// SetParamSet sets token params to the global param store
func (k Keeper) SetParamSet(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

func (k Keeper) setWithOwner(ctx sdk.Context, owner sdk.AccAddress, symbol string) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryBare(&gogotypes.StringValue{Value: symbol})

	store.Set(types.KeyTokens(owner, symbol), bz)
}

func (k Keeper) setWithMinUnit(ctx sdk.Context, minUnit, symbol string) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryBare(&gogotypes.StringValue{Value: symbol})

	store.Set(types.KeyMinUint(minUnit), bz)
}

func (k Keeper) setToken(ctx sdk.Context, token types.Token) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(&token)

	store.Set(types.KeySymbol(token.Symbol), bz)
}

func (k Keeper) getToken(ctx sdk.Context, symbol string) (token types.Token, err error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeySymbol(symbol))
	if bz == nil {
		return token, sdkerrors.Wrap(types.ErrTokenNotExists, fmt.Sprintf("token %s does not exist", symbol))
	}

	k.cdc.MustUnmarshalBinaryBare(bz, &token)
	return token, nil
}

// reset all indices by the new owner for token query
func (k Keeper) resetStoreKeyForQueryToken(ctx sdk.Context, symbol string, srcOwner, dstOwner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	// delete the old key
	store.Delete(types.KeyTokens(srcOwner, symbol))

	// add the new key
	k.setWithOwner(ctx, dstOwner, symbol)
}

// getTokenSupply queries the token supply from the total supply
func (k Keeper) getTokenSupply(ctx sdk.Context, denom string) sdk.Int {
	return k.bankKeeper.GetSupply(ctx).GetTotal().AmountOf(denom)
}
