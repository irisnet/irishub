package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/irisnet/irishub/modules/asset/internal/types"
)

// Keeper defines the module module Keeper
type Keeper struct {
	storeKey  sdk.StoreKey
	cdc       *codec.Codec
	codespace sdk.CodespaceType

	// params subspace
	paramSpace params.Subspace
	// The supplyKeeper to reduce the supply of the network
	supplyKeeper types.SupplyKeeper

	feeCollectorName string
}

// NewKeeper returns a asset keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace,
	codespace sdk.CodespaceType, supplyKeeper types.SupplyKeeper, feeCollectorName string,
) Keeper {
	// ensure asset module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		storeKey:         key,
		cdc:              cdc,
		codespace:        codespace,
		paramSpace:       paramSpace.WithKeyTable(ParamKeyTable()),
		supplyKeeper:     supplyKeeper,
		feeCollectorName: feeCollectorName,
	}
}

// Codespace returns the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

// IssueToken issues a new token
func (k Keeper) IssueToken(ctx sdk.Context, token types.FungibleToken) sdk.Error {
	token, owner, err := k.AddToken(ctx, token)
	if err != nil {
		return err
	}

	initialSupply := sdk.NewCoin(token.GetDenom(), token.GetInitSupply())
	// for native and gateway tokens
	if owner != nil {
		// mint coins
		mintCoins := sdk.NewCoins(initialSupply)
		if err := k.supplyKeeper.MintCoins(ctx, types.ModuleName, mintCoins); err != nil {
			return err
		}

		// sent coins to owner's account
		if err = k.supplyKeeper.SendCoinsFromModuleToAccount(
			ctx, types.ModuleName, owner, mintCoins,
		); err != nil {
			return err
		}
	}

	return nil
}

// AddToken adds a new token to keystore
func (k Keeper) AddToken(ctx sdk.Context, token types.FungibleToken) (types.FungibleToken, sdk.AccAddress, sdk.Error) {
	token.Sanitize()
	tokenId, err := types.GetTokenID(token.GetSource(), token.GetSymbol())
	if err != nil {
		return token, nil, err
	}
	if k.HasToken(ctx, tokenId) {
		return token, nil, types.ErrAssetAlreadyExists(k.codespace, fmt.Sprintf("token already exists: %s", token.GetUniqueID()))
	}

	var owner sdk.AccAddress
	if token.GetSource() == types.NATIVE {
		owner = token.GetOwner()
		token.CanonicalSymbol = ""
	}

	if err = k.setToken(ctx, token); err != nil {
		return token, nil, err
	}

	// Set token to be prefixed with owner and source
	if token.GetSource() == types.NATIVE {
		if err = k.setTokens(ctx, owner, token); err != nil {
			return token, nil, err
		}
	}

	// Set token to be prefixed with source
	if err = k.setTokens(ctx, sdk.AccAddress{}, token); err != nil {
		return token, nil, err
	}

	return token, owner, nil
}

// HasToken checks if the token exists
func (k Keeper) HasToken(ctx sdk.Context, tokenId string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(KeyToken(tokenId))
}

// save token
func (k Keeper) setToken(ctx sdk.Context, token types.FungibleToken) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(token)

	tokenId, err := types.GetTokenID(token.GetSource(), token.GetSymbol())
	if err != nil {
		return err
	}

	store.Set(KeyToken(tokenId), bz)
	return nil
}

// save tokens' owner
func (k Keeper) setTokens(ctx sdk.Context, owner sdk.AccAddress, token types.FungibleToken) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	tokenId, err := types.GetTokenID(token.GetSource(), token.GetSymbol())
	if err != nil {
		return err
	}

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(tokenId)

	store.Set(KeyTokens(owner, tokenId), bz)
	return nil
}

// GetToken returns token by specified tokenID
func (k Keeper) GetToken(ctx sdk.Context, tokenID string) (token types.FungibleToken, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyToken(tokenID))
	if bz == nil {
		return token, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &token)
	return token, true
}

// GetTokens returns tokens by specified owner
func (k Keeper) GetTokens(ctx sdk.Context, owner sdk.AccAddress, nonSymbolTokenId string) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, KeyTokens(owner, nonSymbolTokenId))
}

// EditToken edits the specified token
func (k Keeper) EditToken(ctx sdk.Context, msg types.MsgEditToken) sdk.Error {
	// get the destination token
	token, exist := k.GetToken(ctx, msg.TokenID)
	if !exist {
		return types.ErrAssetNotExists(k.codespace, fmt.Sprintf("token %s does not exist", msg.TokenID))
	}

	if !msg.Owner.Equals(token.Owner) {
		return types.ErrInvalidOwner(k.codespace, fmt.Sprintf("the address %d is not the owner of the token %s", msg.Owner, msg.TokenID))
	}

	hasIssuedAmt := k.AssetTokenSupply(ctx, token.GetDenom())

	maxSupply := sdk.NewIntWithDecimal(int64(msg.MaxSupply), int(token.Decimal))
	if maxSupply.GT(sdk.ZeroInt()) && (maxSupply.LT(hasIssuedAmt) || maxSupply.GT(token.MaxSupply)) {
		return types.ErrInvalidAssetMaxSupply(k.codespace, fmt.Sprintf("max supply must not be less than %s and greater than %s", hasIssuedAmt.String(), token.MaxSupply.String()))
	}
	if msg.Name != types.DoNotModify {
		token.Name = msg.Name
	}
	if msg.CanonicalSymbol != types.DoNotModify && token.Source != types.NATIVE {
		token.CanonicalSymbol = msg.CanonicalSymbol
	}
	if msg.MinUnitAlias != types.DoNotModify {
		token.MinUnitAlias = msg.MinUnitAlias
	}
	if maxSupply.GT(sdk.ZeroInt()) {
		token.MaxSupply = maxSupply
	}
	if msg.Mintable != types.Nil {
		token.Mintable = msg.Mintable.ToBool()
	}

	return k.setToken(ctx, token)
}

// IterateTokens iterates through all existing tokens
func (k Keeper) IterateTokens(ctx sdk.Context, op func(token types.FungibleToken) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, PrefixToken)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var token types.FungibleToken
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &token)

		if stop := op(token); stop {
			break
		}
	}
}

// TODO: delete
// Init
func (k Keeper) Init(ctx sdk.Context) {
	ctx = ctx.WithLogger(ctx.Logger().With("handler", "Init").With("module", "iris/asset"))

	//Initialize external tokens BTC and ETH
	maxSupply := sdk.NewIntWithDecimal(int64(types.MaximumAssetMaxSupply), 8)
	btc := types.NewFungibleToken(types.EXTERNAL, "BTC", "Bitcoin", 8, "BTC", "satoshi", sdk.ZeroInt(), maxSupply, true, nil)
	if err := k.IssueToken(ctx, btc); err != nil {
		ctx.Logger().Error(fmt.Sprintf("initialize external tokens BTC failed:%s", err.Error()))
	}

	maxSupply = sdk.NewIntWithDecimal(int64(types.MaximumAssetMaxSupply), 18)
	eth := types.NewFungibleToken(types.EXTERNAL, "ETH", "Ethereum", 18, "ETH", "wei", sdk.ZeroInt(), maxSupply, true, nil)
	if err := k.IssueToken(ctx, eth); err != nil {
		ctx.Logger().Error(fmt.Sprintf("initialize external tokens ETH failed:%s", err.Error()))
	}
}

// TransferTokenOwner transfers the owner of the specified token to a new one
func (k Keeper) TransferTokenOwner(ctx sdk.Context, msg types.MsgTransferTokenOwner) sdk.Error {
	// get the destination token
	token, exist := k.GetToken(ctx, msg.TokenID)
	if !exist {
		return types.ErrAssetNotExists(k.codespace, fmt.Sprintf("token %s does not exist", msg.TokenID))
	}

	if token.Source != types.NATIVE {
		return types.ErrInvalidAssetSource(k.codespace, fmt.Sprintf("only the token of which the source is native can be transferred,but the source of the current token is %s", token.Source.String()))
	}

	if !msg.SrcOwner.Equals(token.Owner) {
		return types.ErrInvalidOwner(k.codespace, fmt.Sprintf("the address %s is not the owner of the token %s", msg.SrcOwner.String(), msg.TokenID))
	}

	token.Owner = msg.DstOwner

	// update token information
	if err := k.setToken(ctx, token); err != nil {
		return err
	}

	// reset all index for query-token
	return k.resetStoreKeyForQueryToken(ctx, msg, token)
}

// reset all index by DstOwner of token for query-token command
func (k Keeper) resetStoreKeyForQueryToken(ctx sdk.Context, msg types.MsgTransferTokenOwner, token types.FungibleToken) sdk.Error {
	store := ctx.KVStore(k.storeKey)

	tokenId, err := types.GetTokenID(token.GetSource(), token.GetSymbol())
	if err != nil {
		return err
	}
	// delete the old key
	store.Delete(KeyTokens(msg.SrcOwner, tokenId))

	// add the new key
	return k.setTokens(ctx, msg.DstOwner, token)
}

// MintToken handles MsgMintToken
func (k Keeper) MintToken(ctx sdk.Context, msg types.MsgMintToken) sdk.Error {
	token, exist := k.GetToken(ctx, msg.TokenID)
	if !exist {
		return types.ErrAssetNotExists(k.codespace, fmt.Sprintf("token %s does not exist", msg.TokenID))
	}

	if !msg.Owner.Equals(token.Owner) {
		return types.ErrInvalidOwner(k.codespace, fmt.Sprintf("the address %s is not the owner of the token %s", msg.Owner.String(), msg.TokenID))
	}

	if !token.Mintable {
		return types.ErrAssetNotMintable(k.codespace, fmt.Sprintf("the token %s is set to be non-mintable", msg.TokenID))
	}

	hasIssuedAmt := k.AssetTokenSupply(ctx, token.GetDenom())

	mintAmt := sdk.NewIntWithDecimal(int64(msg.Amount), int(token.Decimal))
	if mintAmt.Add(hasIssuedAmt).GT(token.MaxSupply) {
		exp := sdk.NewIntWithDecimal(1, int(token.Decimal))
		canAmt := token.MaxSupply.Sub(hasIssuedAmt).Quo(exp)
		return types.ErrInvalidAssetMaxSupply(k.codespace, fmt.Sprintf("The amount of mint tokens plus the total amount of issues has exceeded the maximum issue total,only accepts amount (0, %s]", canAmt.String()))
	}

	switch token.Source {
	case types.NATIVE:
		// handle fee for native token
		if err := TokenMintFeeHandler(ctx, k, msg.Owner, token.Symbol); err != nil {
			return err
		}
		break
	default:
		break
	}

	mintCoins := sdk.NewCoins(sdk.NewCoin(token.GetDenom(), mintAmt))

	mintAcc := msg.To
	if mintAcc.Empty() {
		mintAcc = token.Owner
	}

	// mint coins
	if err := k.supplyKeeper.MintCoins(ctx, types.ModuleName, mintCoins); err != nil {
		return err
	}

	// sent coins to owner's account
	return k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, mintAcc, mintCoins)
}

// AssetTokenSupply asset tokens from the total supply
func (k Keeper) AssetTokenSupply(ctx sdk.Context, denom string) sdk.Int {
	return k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(denom)
}

// AddCollectedFees implements an alias call to the underlying supply keeper's
// AddCollectedFees to be used in BeginBlocker.
func (k Keeper) AddCollectedFees(ctx sdk.Context, fees sdk.Coins) sdk.Error {
	return k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, fees)
}
