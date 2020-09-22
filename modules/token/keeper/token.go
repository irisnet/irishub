package keeper

import (
	"fmt"

	gogotypes "github.com/gogo/protobuf/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

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
	if err := k.setToken(ctx, token); err != nil {
		return err
	}

	// Set token to be prefixed with owner
	if err := k.setWithOwner(ctx, token.Owner, token.Symbol); err != nil {
		return err
	}

	// Set token to be prefixed with min_unit
	return k.setWithMinUnit(ctx, token.MinUnit, token.Symbol)
}

// HasToken asserts a token exists
func (k Keeper) HasToken(ctx sdk.Context, denom string) bool {
	store := ctx.KVStore(k.storeKey)
	existed := store.Has(types.KeySymbol(denom))
	if existed {
		return existed
	}

	return store.Has(types.KeyMinUint(denom))
}

// GetParamSet returns token params from the global param store
func (k Keeper) GetParamSet(ctx sdk.Context) types.Params {
	var p types.Params
	k.paramSpace.GetParamSet(ctx, &p)
	return p
}

// SetParamSet set token params from the global param store
func (k Keeper) SetParamSet(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

func (k Keeper) setWithOwner(ctx sdk.Context, owner sdk.AccAddress, symbol string) error {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryBare(&gogotypes.StringValue{Value: symbol})

	store.Set(types.KeyTokens(owner, symbol), bz)
	return nil
}

func (k Keeper) setWithMinUnit(ctx sdk.Context, minUnit, symbol string) error {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryBare(&gogotypes.StringValue{Value: symbol})

	store.Set(types.KeyMinUint(minUnit), bz)
	return nil
}

func (k Keeper) setToken(ctx sdk.Context, token types.Token) error {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(&token)

	store.Set(types.KeySymbol(token.Symbol), bz)
	return nil
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

// reset all index by DstOwner of token for query-token command
func (k Keeper) resetStoreKeyForQueryToken(ctx sdk.Context, msg types.MsgTransferTokenOwner, token types.Token) error {
	store := ctx.KVStore(k.storeKey)

	// delete the old key
	store.Delete(types.KeyTokens(msg.SrcOwner, token.Symbol))

	// add the new key
	return k.setWithOwner(ctx, msg.DstOwner, token.Symbol)
}

// getTokenSupply query issued tokens supply from the total supply
func (k Keeper) getTokenSupply(ctx sdk.Context, denom string) sdk.Int {
	return k.bankKeeper.GetSupply(ctx).GetTotal().AmountOf(denom)
}
