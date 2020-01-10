package keeper

import (
	"fmt"

	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/app/v3/asset/internal/types"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
	bk       types.BankKeeper

	// codespace
	codespace sdk.CodespaceType
	// params subspace
	paramSpace params.Subspace
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, bk types.BankKeeper, codespace sdk.CodespaceType, paramSpace params.Subspace) Keeper {
	return Keeper{
		storeKey:   key,
		cdc:        cdc,
		bk:         bk,
		codespace:  codespace,
		paramSpace: paramSpace.WithTypeTable(types.ParamTypeTable()),
	}
}

// Codespace return the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

// IssueToken issue a new token
func (k Keeper) IssueToken(ctx sdk.Context, token types.FungibleToken) (sdk.Tags, sdk.Error) {
	if err := k.AddToken(ctx, token); err != nil {
		return nil, err
	}

	initialSupply := sdk.NewCoin(token.GetDenom(), token.GetInitSupply())
	// Add coins into owner's account
	if _, _, err := k.bk.AddCoins(ctx, token.Owner, sdk.Coins{initialSupply}); err != nil {
		return nil, err
	}

	// Set total supply
	k.bk.SetTotalSupply(ctx, initialSupply)
	if initialSupply.Amount.GT(sdk.ZeroInt()) {
		ctx.CoinFlowTags().AppendCoinFlowTag(ctx, token.Owner.String(), token.Owner.String(), initialSupply.String(), sdk.IssueTokenFlow, "")
	}

	createTags := sdk.NewTags(
		types.TagId, []byte(token.GetUniqueID()),
		types.TagDenom, []byte(token.GetDenom()),
		types.TagOwner, []byte(token.GetOwner().String()),
	)

	return createTags, nil
}

// EditToken edits the specified token
func (k Keeper) EditToken(ctx sdk.Context, msg types.MsgEditToken) (sdk.Tags, sdk.Error) {
	// get the destination token
	token, err := k.getToken(ctx, msg.TokenId)
	if err != nil {
		return nil, err
	}

	if !msg.Owner.Equals(token.Owner) {
		return nil, types.ErrInvalidOwner(k.codespace, fmt.Sprintf("the address %d is not the owner of the token %s", msg.Owner, msg.TokenId))
	}

	hasIssuedAmt, found := k.bk.GetTotalSupply(ctx, token.GetDenom())
	if !found {
		return nil, types.ErrAssetNotExists(k.codespace, fmt.Sprintf("token denom %s does not exist", token.GetDenom()))
	}

	if msg.MaxSupply > 0 {
		maxSupply := sdk.NewIntWithDecimal(int64(msg.MaxSupply), int(token.Decimal))
		if maxSupply.LT(hasIssuedAmt.Amount) {
			return nil, types.ErrInvalidAssetMaxSupply(k.codespace, fmt.Sprintf("max supply must not be less than %s", hasIssuedAmt.Amount.String()))
		}
		token.MaxSupply = maxSupply
	}

	if msg.Name != types.DoNotModify {
		token.Name = msg.Name
	}

	if msg.Mintable != types.Nil {
		token.Mintable = msg.Mintable.ToBool()
	}

	if err := k.setToken(ctx, token); err != nil {
		return nil, err
	}

	editTags := sdk.NewTags(
		types.TagId, []byte(msg.TokenId),
	)

	return editTags, nil
}

// TransferTokenOwner transfers the owner of the specified token to a new one
func (k Keeper) TransferTokenOwner(ctx sdk.Context, msg types.MsgTransferTokenOwner) (sdk.Tags, sdk.Error) {
	// get the destination token
	token, err := k.getToken(ctx, msg.TokenId)
	if err != nil {
		return nil, err
	}

	if !msg.SrcOwner.Equals(token.Owner) {
		return nil, types.ErrInvalidOwner(k.codespace, fmt.Sprintf("the address %s is not the owner of the token %s", msg.SrcOwner.String(), msg.TokenId))
	}

	token.Owner = msg.DstOwner
	// update token information
	if err := k.setToken(ctx, token); err != nil {
		return nil, err
	}

	// reset all index for query-token
	if err := k.resetStoreKeyForQueryToken(ctx, msg, token); err != nil {
		return nil, err
	}
	tags := sdk.NewTags(
		types.TagId, []byte(msg.TokenId),
	)

	return tags, nil
}

// MintToken mint specified amount token to a specified owner
func (k Keeper) MintToken(ctx sdk.Context, msg types.MsgMintToken) (sdk.Tags, sdk.Error) {
	token, err := k.getToken(ctx, msg.TokenId)
	if err != nil {
		return nil, err
	}

	if !msg.Owner.Equals(token.Owner) {
		return nil, types.ErrInvalidOwner(k.codespace, fmt.Sprintf("the address %s is not the owner of the token %s", msg.Owner.String(), msg.TokenId))
	}

	if !token.Mintable {
		return nil, types.ErrAssetNotMintable(k.codespace, fmt.Sprintf("the token %s is set to be non-mintable", msg.TokenId))
	}

	hasIssuedAmt, found := k.bk.GetTotalSupply(ctx, token.GetDenom())
	if !found {
		return nil, types.ErrAssetNotExists(k.codespace, fmt.Sprintf("token denom %s does not exist", token.GetDenom()))
	}

	//check the denom
	expDenom := token.GetDenom()
	if expDenom != hasIssuedAmt.Denom {
		return nil, types.ErrAssetNotExists(k.codespace, fmt.Sprintf("denom of mint token is not equal issued token,expected:%s,actual:%s", expDenom, hasIssuedAmt.Denom))
	}

	mintAmt := sdk.NewIntWithDecimal(int64(msg.Amount), int(token.Decimal))
	if mintAmt.Add(hasIssuedAmt.Amount).GT(token.MaxSupply) {
		exp := sdk.NewIntWithDecimal(1, int(token.Decimal))
		canAmt := token.MaxSupply.Sub(hasIssuedAmt.Amount).Div(exp)
		return nil, types.ErrInvalidAssetMaxSupply(k.codespace, fmt.Sprintf("The amount of mint tokens plus the total amount of issues has exceeded the maximum issue total,only accepts amount (0, %s]", canAmt.String()))
	}

	mintCoin := sdk.NewCoin(expDenom, mintAmt)
	//add TotalSupply
	if err := k.bk.IncreaseTotalSupply(ctx, mintCoin); err != nil {
		return nil, err
	}

	mintAcc := msg.To
	if mintAcc.Empty() {
		mintAcc = token.Owner
	}

	//add mintCoin to special account
	_, tags, err := k.bk.AddCoins(ctx, mintAcc, sdk.Coins{mintCoin})
	if err != nil {
		return nil, err
	}
	ctx.CoinFlowTags().AppendCoinFlowTag(ctx, msg.Owner.String(), mintAcc.String(), mintCoin.String(), sdk.MintTokenFlow, "")
	return tags, nil
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

// AddToken save a new token
func (k Keeper) AddToken(ctx sdk.Context, token types.FungibleToken) sdk.Error {
	if k.HasToken(ctx, types.GetTokenID(token.GetSymbol())) {
		return types.ErrAssetAlreadyExists(k.codespace, fmt.Sprintf("token already exists: %s", token.GetUniqueID()))
	}

	//
	if err := k.setToken(ctx, token); err != nil {
		return err
	}

	// Set token to be prefixed with owner and source
	if err := k.setTokens(ctx, token.GetOwner(), token); err != nil {
		return err
	}

	// Set token to be prefixed with source
	if err := k.setTokens(ctx, sdk.AccAddress{}, token); err != nil {
		return err
	}

	return nil
}

// HasToken asset a token exited
func (k Keeper) HasToken(ctx sdk.Context, tokenId string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(KeyToken(tokenId))
}

// GetParamSet return asset params from the global param store
func (k Keeper) GetParamSet(ctx sdk.Context) types.Params {
	var p types.Params
	k.paramSpace.GetParamSet(ctx, &p)
	return p
}

// SetParamSet set asset params from the global param store
func (k Keeper) SetParamSet(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

func (k Keeper) iterateTokensWithOwner(ctx sdk.Context, owner sdk.AccAddress, op func(token types.FungibleToken) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, KeyTokens(owner, ""))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var tokenId string
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &tokenId)
		token, err := k.getToken(ctx, tokenId)
		if err != nil {
			continue
		}
		if stop := op(token); stop {
			break
		}
	}
}

func (k Keeper) setTokens(ctx sdk.Context, owner sdk.AccAddress, token types.FungibleToken) sdk.Error {
	store := ctx.KVStore(k.storeKey)

	tokenId := types.GetTokenID(token.GetSymbol())

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(tokenId)

	store.Set(KeyTokens(owner, tokenId), bz)
	return nil
}

func (k Keeper) setToken(ctx sdk.Context, token types.FungibleToken) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(token)

	tokenId := types.GetTokenID(token.GetSymbol())
	store.Set(KeyToken(tokenId), bz)
	return nil
}

// reset all index by DstOwner of token for query-token command
func (k Keeper) resetStoreKeyForQueryToken(ctx sdk.Context, msg types.MsgTransferTokenOwner, token types.FungibleToken) sdk.Error {
	store := ctx.KVStore(k.storeKey)

	tokenId := types.GetTokenID(token.GetSymbol())
	// delete the old key
	store.Delete(KeyTokens(msg.SrcOwner, tokenId))

	// add the new key
	return k.setTokens(ctx, msg.DstOwner, token)
}

func (k Keeper) getToken(ctx sdk.Context, tokenId string) (token types.FungibleToken, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyToken(tokenId))
	if bz == nil {
		return token, types.ErrAssetNotExists(k.codespace, fmt.Sprintf("token %s does not exist", tokenId))
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &token)
	return token, nil
}
