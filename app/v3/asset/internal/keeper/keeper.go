package keeper

import (
	"fmt"
	"strings"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/app/v3/asset/exported"
	"github.com/irisnet/irishub/app/v3/asset/internal/types"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
	bk       types.BankKeeper
	gk       types.GuardianKeeper

	// codespace
	codespace sdk.CodespaceType
	// params subspace
	paramSpace params.Subspace
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, bk types.BankKeeper, gk types.GuardianKeeper, codespace sdk.CodespaceType, paramSpace params.Subspace) Keeper {
	keeper := Keeper{
		storeKey:   key,
		cdc:        cdc,
		bk:         bk,
		gk:         gk,
		codespace:  codespace,
		paramSpace: paramSpace.WithTypeTable(types.ParamTypeTable()),
	}

	auth.RegisterTotalSupplyKeyGen(types.TotalSupplyKeyGen)

	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("%s/%s", "iris", types.MsgRoute))
}

// IssueToken issues a new token
func (k Keeper) IssueToken(ctx sdk.Context, msg types.MsgIssueToken) (sdk.Tags, sdk.Error) {
	symbol := strings.ToLower(msg.Symbol)
	name := strings.TrimSpace(msg.Name)
	minUnitAlias := strings.ToLower(strings.TrimSpace(msg.MinUnitAlias))

	token := types.NewFungibleToken(
		symbol, name, minUnitAlias, msg.Decimal, msg.InitialSupply,
		msg.MaxSupply, msg.Mintable, msg.Owner,
	)

	if err := k.AddToken(ctx, token); err != nil {
		return nil, err
	}

	initialSupply := sdk.NewCoin(
		token.GetDenom(),
		sdk.NewIntWithDecimal(int64(msg.InitialSupply), int(msg.Decimal)),
	)

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
		types.TagSymbol, []byte(token.GetSymbol()),
		types.TagDenom, []byte(token.GetDenom()),
		types.TagOwner, []byte(token.GetOwner().String()),
	)

	return createTags, nil
}

// EditToken edits the specified token
func (k Keeper) EditToken(ctx sdk.Context, msg types.MsgEditToken) (sdk.Tags, sdk.Error) {
	// get the destination token
	token, err := k.getToken(ctx, msg.Symbol)
	if err != nil {
		return nil, err
	}

	if !msg.Owner.Equals(token.Owner) {
		return nil, types.ErrInvalidOwner(k.codespace, fmt.Sprintf("the address %d is not the owner of the token %s", msg.Owner, msg.Symbol))
	}

	issuedAmt, found := k.bk.GetTotalSupply(ctx, token.GetDenom())
	if !found {
		return nil, types.ErrAssetNotExists(k.codespace, fmt.Sprintf("token denom %s does not exist", token.GetDenom()))
	}

	if msg.MaxSupply > 0 {
		issuedMainUnitAmt := uint64(issuedAmt.Amount.Div(sdk.NewIntWithDecimal(1, int(token.Decimal))).Int64())
		if msg.MaxSupply < issuedMainUnitAmt {
			return nil, types.ErrInvalidAssetMaxSupply(k.codespace, fmt.Sprintf("max supply must not be less than %d", issuedMainUnitAmt))
		}

		token.MaxSupply = msg.MaxSupply
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
		types.TagSymbol, []byte(msg.Symbol),
	)

	return editTags, nil
}

// TransferTokenOwner transfers the owner of the specified token to a new one
func (k Keeper) TransferTokenOwner(ctx sdk.Context, msg types.MsgTransferTokenOwner) (sdk.Tags, sdk.Error) {
	// get the destination token
	token, err := k.getToken(ctx, msg.Symbol)
	if err != nil {
		return nil, err
	}

	if !msg.SrcOwner.Equals(token.Owner) {
		return nil, types.ErrInvalidOwner(k.codespace, fmt.Sprintf("the address %s is not the owner of the token %s", msg.SrcOwner.String(), msg.Symbol))
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
		types.TagSymbol, []byte(msg.Symbol),
	)

	return tags, nil
}

// MintToken mints specified amount token to a specified owner
func (k Keeper) MintToken(ctx sdk.Context, msg types.MsgMintToken) (sdk.Tags, sdk.Error) {
	token, err := k.getToken(ctx, msg.Symbol)
	if err != nil {
		return nil, err
	}

	if !msg.Owner.Equals(token.Owner) {
		return nil, types.ErrInvalidOwner(k.codespace, fmt.Sprintf("the address %s is not the owner of the token %s", msg.Owner.String(), msg.Symbol))
	}

	if !token.Mintable {
		return nil, types.ErrAssetNotMintable(k.codespace, fmt.Sprintf("the token %s is set to be non-mintable", msg.Symbol))
	}

	issuedAmt, found := k.bk.GetTotalSupply(ctx, token.GetDenom())
	if !found {
		return nil, types.ErrAssetNotExists(k.codespace, fmt.Sprintf("token denom %s does not exist", token.GetDenom()))
	}

	//check the denom
	expDenom := token.GetDenom()
	if expDenom != issuedAmt.Denom {
		return nil, types.ErrAssetNotExists(k.codespace, fmt.Sprintf("denom of minting token is not equal to the issued token, expected:%s, actual:%s", expDenom, issuedAmt.Denom))
	}

	mintableMaxAmt := sdk.NewIntWithDecimal(int64(token.MaxSupply), int(token.Decimal)).Sub(issuedAmt.Amount)
	mintableMaxMainUnitAmt := uint64(mintableMaxAmt.Div(sdk.NewIntWithDecimal(1, int(token.Decimal))).Int64())

	if msg.Amount > mintableMaxMainUnitAmt {
		return nil, types.ErrInvalidAssetMaxSupply(k.codespace, fmt.Sprintf("The amount of minting tokens plus the total amount of issued tokens has exceeded the maximum supply, only accepts amount (0, %d]", mintableMaxMainUnitAmt))
	}

	mintCoin := sdk.NewCoin(expDenom, sdk.NewIntWithDecimal(int64(msg.Amount), int(token.Decimal)))

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

// GetAllTokens returns all existing tokens
func (k Keeper) GetAllTokens(ctx sdk.Context) (tokens []exported.TokenI) {
	k.IterateTokens(ctx, func(token types.FungibleToken) (stop bool) {
		tokens = append(tokens, token)
		return false
	})

	return
}

// AddToken saves a new token
func (k Keeper) AddToken(ctx sdk.Context, token types.FungibleToken) sdk.Error {
	if k.HasToken(ctx, token.GetSymbol()) {
		return types.ErrAssetAlreadyExists(k.codespace, fmt.Sprintf("token already exists: %s", token.GetSymbol()))
	}

	// set token
	if err := k.setToken(ctx, token); err != nil {
		return err
	}

	// Set token to be prefixed with owner
	if err := k.setOwnerToken(ctx, token.GetOwner(), token); err != nil {
		return err
	}

	return nil
}

// HasToken asserts a token exists
func (k Keeper) HasToken(ctx sdk.Context, symbol string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(KeyToken(symbol))
}

// GetParamSet returns asset params from the global param store
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
		var symbol string
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &symbol)

		token, err := k.getToken(ctx, symbol)
		if err != nil {
			continue
		}

		if stop := op(token); stop {
			break
		}
	}
}

func (k Keeper) setOwnerToken(ctx sdk.Context, owner sdk.AccAddress, token types.FungibleToken) sdk.Error {
	store := ctx.KVStore(k.storeKey)

	symbol := token.GetSymbol()
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(symbol)

	store.Set(KeyTokens(owner, symbol), bz)
	return nil
}

func (k Keeper) setToken(ctx sdk.Context, token types.FungibleToken) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(token)

	store.Set(KeyToken(token.GetSymbol()), bz)
	return nil
}

// reset all index by DstOwner of token for query-token command
func (k Keeper) resetStoreKeyForQueryToken(ctx sdk.Context, msg types.MsgTransferTokenOwner, token types.FungibleToken) sdk.Error {
	store := ctx.KVStore(k.storeKey)

	// delete the old key
	store.Delete(KeyTokens(msg.SrcOwner, token.GetSymbol()))

	// add the new key
	return k.setOwnerToken(ctx, msg.DstOwner, token)
}

func (k Keeper) getToken(ctx sdk.Context, symbol string) (token types.FungibleToken, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(KeyToken(symbol))
	if bz == nil {
		return token, types.ErrAssetNotExists(k.codespace, fmt.Sprintf("token %s does not exist", symbol))
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &token)
	return token, nil
}

// GetToken wraps getToken for export
func (k Keeper) GetToken(ctx sdk.Context, symbol string) (token exported.TokenI, err sdk.Error) {
	return k.getToken(ctx, symbol)
}
