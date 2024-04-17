package keeper

import (
	"fmt"

	gogotypes "github.com/cosmos/gogoproto/types"
	"github.com/ethereum/go-ethereum/common"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/irisnet/irismod/modules/token/types"
	v1 "github.com/irisnet/irismod/modules/token/types/v1"
)

// GetTokens returns all existing tokens
func (k Keeper) GetTokens(ctx sdk.Context, owner sdk.AccAddress) (tokens []v1.TokenI) {
	store := ctx.KVStore(k.storeKey)

	var it sdk.Iterator
	if owner == nil {
		it = sdk.KVStorePrefixIterator(store, types.PrefixTokenForSymbol)
		defer it.Close()

		for ; it.Valid(); it.Next() {
			var token v1.Token
			k.cdc.MustUnmarshal(it.Value(), &token)

			tokens = append(tokens, &token)
		}
		return
	}

	it = sdk.KVStorePrefixIterator(store, types.KeyTokens(owner, ""))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		var symbol gogotypes.StringValue
		k.cdc.MustUnmarshal(it.Value(), &symbol)

		token, err := k.getTokenBySymbol(ctx, symbol.Value)
		if err != nil {
			continue
		}
		tokens = append(tokens, token)
	}
	return
}

// GetToken returns the token of the specified symbol or min uint
func (k Keeper) GetToken(ctx sdk.Context, denom string) (v1.TokenI, error) {
	// query token by symbol
	if token, err := k.getTokenBySymbol(ctx, denom); err == nil {
		return &token, nil
	}

	// query token by min unit
	if token, err := k.getTokenByMinUnit(ctx, denom); err == nil {
		return &token, nil
	}

	return nil, errorsmod.Wrapf(types.ErrTokenNotExists, "token: %s does not exist", denom)
}

// AddToken saves a new token
func (k Keeper) AddToken(ctx sdk.Context, token v1.Token) error {
	if k.HasToken(ctx, token.Symbol) {
		return errorsmod.Wrapf(
			types.ErrSymbolAlreadyExists,
			"symbol already exists: %s",
			token.Symbol,
		)
	}

	if k.HasToken(ctx, token.MinUnit) {
		return errorsmod.Wrapf(
			types.ErrMinUnitAlreadyExists,
			"min-unit already exists: %s",
			token.MinUnit,
		)
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

// HasSymbol asserts a token exists by symbol
func (k Keeper) HasSymbol(ctx sdk.Context, symbol string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.KeySymbol(symbol))
}

// HasMinUint asserts a token exists by minUint
func (k Keeper) HasMinUint(ctx sdk.Context, minUint string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.KeyMinUint(minUint))
}

// HasToken asserts a token exists
func (k Keeper) HasToken(ctx sdk.Context, denom string) bool {
	if k.HasSymbol(ctx, denom) {
		return true
	}
	return k.HasMinUint(ctx, denom)
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

	bz := k.cdc.MustMarshal(&total)
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
		return sdk.Coin{}, errorsmod.Wrapf(
			types.ErrNotFoundTokenAmt,
			"not found symbol: %s",
			minUint,
		)
	}

	var coin sdk.Coin
	k.cdc.MustUnmarshal(bz, &coin)

	return coin, nil
}

// GetAllBurnCoin returns the total amount of all the burned tokens
func (k Keeper) GetAllBurnCoin(ctx sdk.Context) []sdk.Coin {
	store := ctx.KVStore(k.storeKey)

	var coins []sdk.Coin
	it := sdk.KVStorePrefixIterator(store, types.PrefixBurnTokenAmt)
	for ; it.Valid(); it.Next() {
		var coin sdk.Coin
		k.cdc.MustUnmarshal(it.Value(), &coin)
		coins = append(coins, coin)
	}

	return coins
}

func (k Keeper) setWithOwner(ctx sdk.Context, owner sdk.AccAddress, symbol string) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&gogotypes.StringValue{Value: symbol})

	store.Set(types.KeyTokens(owner, symbol), bz)
}

func (k Keeper) setWithMinUnit(ctx sdk.Context, minUnit, symbol string) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&gogotypes.StringValue{Value: symbol})

	store.Set(types.KeyMinUint(minUnit), bz)
}

func (k Keeper) setWithContract(ctx sdk.Context, contract, symbol string) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&gogotypes.StringValue{Value: symbol})

	store.Set(types.KeyContract(contract), bz)
}

func (k Keeper) setToken(ctx sdk.Context, token v1.Token) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&token)

	store.Set(types.KeySymbol(token.Symbol), bz)
}

func (k Keeper) getTokenBySymbol(ctx sdk.Context, symbol string) (token v1.Token, err error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeySymbol(symbol))
	if bz == nil {
		return token, errorsmod.Wrap(
			types.ErrTokenNotExists,
			fmt.Sprintf("token symbol %s does not exist", symbol),
		)
	}

	k.cdc.MustUnmarshal(bz, &token)
	return token, nil
}

func (k Keeper) getTokenByMinUnit(ctx sdk.Context, minUnit string) (token v1.Token, err error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyMinUint(minUnit))
	if bz == nil {
		return token, errorsmod.Wrap(
			types.ErrTokenNotExists,
			fmt.Sprintf("token minUnit %s does not exist", minUnit),
		)
	}

	var symbol gogotypes.StringValue
	k.cdc.MustUnmarshal(bz, &symbol)

	token, err = k.getTokenBySymbol(ctx, symbol.Value)
	if err != nil {
		return token, err
	}

	return token, nil
}

func (k Keeper) getTokenByContract(ctx sdk.Context, contract common.Address) (token v1.Token, err error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyContract(contract.Hex()))
	if bz == nil {
		return token, errorsmod.Wrap(
			types.ErrTokenNotExists,
			fmt.Sprintf("token contract %s does not exist", contract),
		)
	}

	var symbol gogotypes.StringValue
	k.cdc.MustUnmarshal(bz, &symbol)

	token, err = k.getTokenBySymbol(ctx, symbol.Value)
	if err != nil {
		return token, err
	}

	return token, nil
}

func (k Keeper) getSymbolByMinUnit(ctx sdk.Context, minUnit string) (string, error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyMinUint(minUnit))
	if bz == nil {
		return "", errorsmod.Wrap(
			types.ErrTokenNotExists,
			fmt.Sprintf("token minUnit %s does not exist", minUnit),
		)
	}

	var symbol gogotypes.StringValue
	k.cdc.MustUnmarshal(bz, &symbol)

	return symbol.Value, nil
}

// reset all indices by the new owner for token query
func (k Keeper) resetStoreKeyForQueryToken(
	ctx sdk.Context,
	symbol string,
	srcOwner, dstOwner sdk.AccAddress,
) {
	store := ctx.KVStore(k.storeKey)

	// delete the old key
	store.Delete(types.KeyTokens(srcOwner, symbol))

	// add the new key
	k.setWithOwner(ctx, dstOwner, symbol)
}

// getTokenSupply queries the token supply from the total supply
func (k Keeper) getTokenSupply(ctx sdk.Context, denom string) sdk.Int {
	return k.bankKeeper.GetSupply(ctx, denom).Amount
}


// upsertToken updates or inserts a token into the database.
//
// ctx: the context in which the token is being upserted.
// token: the token struct to be upserted.
func (k Keeper) upsertToken(ctx sdk.Context, token v1.Token) {
	// set token
	k.setToken(ctx, token)
	// set token to be prefixed with min unit
	k.setWithMinUnit(ctx, token.MinUnit, token.Symbol)
	if len(token.Owner) != 0 {
		// set token to be prefixed with owner
		k.setWithOwner(ctx, token.GetOwner(), token.Symbol)
	}
	if len(token.Contract) != 0 {
		// set token to be prefixed with owner
		k.setWithContract(ctx, token.Contract, token.Symbol)
	}
}
