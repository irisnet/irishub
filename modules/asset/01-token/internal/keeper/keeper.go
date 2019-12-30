package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/irisnet/irishub/modules/asset/01-token/internal/types"
)

// Keeper defines the module module Keeper
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec

	// params subspace
	paramSpace params.Subspace
	// The supplyKeeper to reduce the supply of the network
	supplyKeeper types.SupplyKeeper

	feeCollectorName string
}

// NewKeeper returns a asset keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace,
	supplyKeeper types.SupplyKeeper, feeCollectorName string,
) Keeper {
	// ensure asset module account is set
	if addr := supplyKeeper.GetModuleAddress(types.SubModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.SubModuleName))
	}

	return Keeper{
		storeKey:         key,
		cdc:              cdc,
		paramSpace:       paramSpace.WithKeyTable(ParamKeyTable()),
		supplyKeeper:     supplyKeeper,
		feeCollectorName: feeCollectorName,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("%s", types.SubModuleName))
}

// IssueToken issues a new token
func (k Keeper) IssueToken(ctx sdk.Context, token types.FungibleToken) error {
	if err := k.AddToken(ctx, token); err != nil {
		return err
	}

	if err := IssueTokenFeeHandler(ctx, k, token.Owner, token.Symbol); err != nil {
		return err
	}

	// mint coins
	mintAmt := sdk.NewIntWithDecimal(token.GetInitSupply().Int64(), int(token.Scale))
	mintCoins := sdk.NewCoins(sdk.NewCoin(token.GetMinUnit(), mintAmt))
	if err := k.supplyKeeper.MintCoins(ctx, types.SubModuleName, mintCoins); err != nil {
		return err
	}

	// sent coins to owner's account
	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(
		ctx, types.SubModuleName, token.Owner, mintCoins,
	); err != nil {
		return err
	}

	return nil
}

// EditToken edits the specified token
func (k Keeper) EditToken(ctx sdk.Context, msg types.MsgEditToken) error {
	// get the destination token
	token, exist := k.GetToken(ctx, msg.Symbol)
	if !exist {
		return sdkerrors.Wrap(types.ErrTokenNotExists, msg.Symbol)
	}

	if !msg.Owner.Equals(token.Owner) {
		return sdkerrors.Wrap(types.ErrInvalidOwner, fmt.Sprintf("the address %d is not the owner of the token %s", msg.Owner, msg.Symbol))
	}

	// check the maxSupply
	if msg.MaxSupply > 0 {
		hasIssuedAmt := k.getTokenSupply(ctx, token.GetMinUnit())
		maxSupply := int64(msg.MaxSupply)
		maxSupplyInt := sdk.NewIntWithDecimal(maxSupply, int(token.Scale))
		if maxSupplyInt.LT(hasIssuedAmt) {
			return sdkerrors.Wrap(types.ErrInvalidAssetMaxSupply, fmt.Sprintf("max supply must not be less than %s", hasIssuedAmt.String()))
		}
		token.MaxSupply = sdk.NewInt(maxSupply)
	}

	if msg.Name != types.DoNotModify {
		token.Name = msg.Name
	}

	if msg.Mintable != types.Nil {
		token.Mintable = msg.Mintable.ToBool()
	}

	k.setToken(ctx, token)

	return nil
}

// TransferToken transfers the owner of the specified token to a new one
func (k Keeper) TransferToken(ctx sdk.Context, msg types.MsgTransferToken) error {
	// get the destination token
	token, exist := k.GetToken(ctx, msg.Symbol)
	if !exist {
		return types.ErrTokenNotExists
	}

	if !msg.SrcOwner.Equals(token.Owner) {
		return sdkerrors.Wrap(types.ErrInvalidOwner, fmt.Sprintf("the address %s is not the owner of the token %s", msg.SrcOwner.String(), msg.Symbol))
	}

	token.Owner = msg.DstOwner

	// update token information
	k.setToken(ctx, token)

	// reset all store key for token
	k.resetTokenKey(ctx, msg.SrcOwner, token)

	return nil
}

// MintToken handles MsgMintToken
func (k Keeper) MintToken(ctx sdk.Context, msg types.MsgMintToken) (sdk.Coins, error) {
	token, exist := k.GetToken(ctx, msg.Symbol)
	if !exist {
		return nil, types.ErrTokenNotExists
	}

	if !msg.Owner.Equals(token.Owner) {
		return nil, sdkerrors.Wrap(types.ErrInvalidOwner, fmt.Sprintf("the address %s is not the owner of the token %s", msg.Owner.String(), msg.Symbol))
	}

	if !token.Mintable {
		return nil, types.ErrAssetNotMintable
	}

	hasIssuedAmt := k.getTokenSupply(ctx, token.GetMinUnit())
	mintAmt := sdk.NewIntWithDecimal(int64(msg.Amount), int(token.Scale))
	maxSupply := sdk.NewIntWithDecimal(token.MaxSupply.Int64(), int(token.Scale))
	if mintAmt.Add(hasIssuedAmt).GT(maxSupply) {
		exp := sdk.NewIntWithDecimal(1, int(token.Scale))
		canAmt := token.MaxSupply.Sub(hasIssuedAmt).Quo(exp)
		return nil, sdkerrors.Wrap(types.ErrInvalidAssetMaxSupply, fmt.Sprintf("The amount that the minting tokens plus the total issued amount has exceeded the maximum supply, only accepts amount (0, %s]", canAmt.String()))
	}

	if err := MintTokenFeeHandler(ctx, k, msg.Owner, token.Symbol); err != nil {
		return nil, err
	}

	mintCoins := sdk.NewCoins(sdk.NewCoin(token.GetMinUnit(), mintAmt))

	mintAcc := msg.To
	if mintAcc.Empty() {
		mintAcc = token.Owner
	}

	// mint coins
	if err := k.supplyKeeper.MintCoins(ctx, types.SubModuleName, mintCoins); err != nil {
		return nil, err
	}

	// sent coins to owner's account
	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.SubModuleName, mintAcc, mintCoins); err != nil {
		return nil, err
	}
	return mintCoins, nil
}

// BurnToken handles MsgBurnToken
func (k Keeper) BurnToken(ctx sdk.Context, msg types.MsgBurnToken) error {
	if err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, msg.Sender, types.SubModuleName, msg.Amount); err != nil {
		return err
	}
	return k.supplyKeeper.BurnCoins(ctx, types.SubModuleName, msg.Amount)
}

// AddToken adds a new token to keystore
func (k Keeper) AddToken(ctx sdk.Context, token types.FungibleToken) error {
	token.Sanitize()
	if exited, msg := k.HasToken(ctx, token); exited {
		return sdkerrors.Wrap(types.ErrTokenAlreadyExists, msg)
	}
	// Set token to be prefixed with Symbol
	k.setToken(ctx, token)

	// Set token to be prefixed with owner and Symbol
	k.setTokens(ctx, token)

	// Set token to be prefixed with minUnit
	k.setTokenMinUnit(ctx, token)

	return nil
}

// HasTokenSymbol checks if the token exists
func (k Keeper) HasToken(ctx sdk.Context, token types.FungibleToken) (bool, string) {
	if exited := k.HasTokenSymbol(ctx, token.Symbol); exited {
		return exited, fmt.Sprintf("token symbol already exists: %s", token.GetSymbol())
	}
	if exited := k.HasTokenMinUnit(ctx, token.Symbol); exited {
		return exited, fmt.Sprintf("token minUnit already exists: %s", token.GetMinUnit())
	}
	return false, ""
}

// HasTokenSymbol checks if the token symbol exists
func (k Keeper) HasTokenSymbol(ctx sdk.Context, symbol string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.KeyToken(symbol))
}

// HasTokenMinUnit checks if the token minUnit exists
func (k Keeper) HasTokenMinUnit(ctx sdk.Context, minUnit string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.KeyMinUnit(minUnit))
}

// GetToken returns token by specified symbol
func (k Keeper) GetToken(ctx sdk.Context, symbol string) (token types.FungibleToken, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyToken(symbol))
	if bz == nil {
		return token, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &token)
	return token, true
}

// GetTokenByMintUint returns token by specified mintUint
func (k Keeper) GetTokenByMintUint(ctx sdk.Context, mintUint string) (types.FungibleToken, bool) {
	store := ctx.KVStore(k.storeKey)
	var symbol string
	bz := store.Get(types.KeyMinUnit(mintUint))
	if bz == nil {
		return types.FungibleToken{}, false
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &symbol)
	return k.GetToken(ctx, symbol)
}

// GetTokens returns tokens by specified owner
func (k Keeper) GetTokens(ctx sdk.Context, owner sdk.AccAddress, symbol string) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.KeyTokens(owner, symbol))
}

// GetAllTokens return all existing tokens
func (k Keeper) GetAllTokens(ctx sdk.Context) (tokens types.Tokens) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.KeyTokenPrefix())
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var token types.FungibleToken
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &token)
		tokens = append(tokens, token)
	}
	return
}

// getTokenSupply query issued tokens supply from the total supply
func (k Keeper) getTokenSupply(ctx sdk.Context, denom string) sdk.Int {
	return k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(denom)
}

// addCollectedFees implements an alias call to the underlying supply keeper's
// addCollectedFees to be used in BeginBlocker.
func (k Keeper) addCollectedFees(ctx sdk.Context, fees sdk.Coins) error {
	return k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.SubModuleName, k.feeCollectorName, fees)
}

// reset all index by DstOwner of token for query-token command
func (k Keeper) resetTokenKey(ctx sdk.Context, srcOwner sdk.AccAddress, token types.FungibleToken) {
	store := ctx.KVStore(k.storeKey)

	// delete the old key
	store.Delete(types.KeyTokens(srcOwner, token.GetSymbol()))

	// add the new key
	k.setTokens(ctx, token)
}

// save token
func (k Keeper) setToken(ctx sdk.Context, token types.FungibleToken) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(token)
	store.Set(types.KeyToken(token.Symbol), bz)
}

// save tokens' owner
func (k Keeper) setTokens(ctx sdk.Context, token types.FungibleToken) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(token.Symbol)
	store.Set(types.KeyTokens(token.Owner, token.Symbol), bz)
}

// save tokens' MinUnit
func (k Keeper) setTokenMinUnit(ctx sdk.Context, token types.FungibleToken) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(token.Symbol)
	store.Set(types.KeyMinUnit(token.MinUnit), bz)
}
