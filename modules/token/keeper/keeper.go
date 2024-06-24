package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"mods.irisnet.org/token/types"
	v1 "mods.irisnet.org/token/types/v1"
)

// Keeper of the token store
type Keeper struct {
	storeKey         storetypes.StoreKey
	cdc              codec.Codec
	bankKeeper       types.BankKeeper
	accountKeeper    types.AccountKeeper
	evmKeeper        types.EVMKeeper
	ics20Keeper      types.ICS20Keeper
	blockedAddrs     map[string]bool
	feeCollectorName string
	authority        string
	registry         v1.SwapRegistry
}

// NewKeeper creates a new instance of Keeper.
//
// Parameters:
//
//	cdc: codec to marshal/unmarshal binary encoding/decoding.
//	key: store key for the module's store.
//	bankKeeper: bank Keeper module for interacting with accounts.
//	accountKeeper: Account Keeper for interacting with accounts.
//	evmKeeper: EVM Keeper module for interacting with Ethereum Virtual Machine transactions.
//	ics20Keeper: ICS20 Keeper module for interacting with ICS20 transactions.
//	feeCollectorName: name of the fee collector.
//	authority: authority string.
//
// Return type: Keeper.
func NewKeeper(
	cdc codec.Codec,
	key storetypes.StoreKey,
	bankKeeper types.BankKeeper,
	accountKeeper types.AccountKeeper,
	evmKeeper types.EVMKeeper,
	ics20Keeper types.ICS20Keeper,
	feeCollectorName string,
	authority string,
) Keeper {
	return Keeper{
		storeKey:         key,
		cdc:              cdc,
		bankKeeper:       bankKeeper,
		accountKeeper:    accountKeeper,
		evmKeeper:        evmKeeper,
		ics20Keeper:      ics20Keeper,
		feeCollectorName: feeCollectorName,
		blockedAddrs:     bankKeeper.GetBlockedAddresses(),
		registry:         make(v1.SwapRegistry),
		authority:        authority,
	}
}

// Codec returns a k.cdc.
func (k Keeper) Codec() codec.Codec {
	return k.cdc
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("irismod/%s", types.ModuleName))
}

// Hooks returns the keeper's Hooks struct.
func (k Keeper) Hooks() types.Hook {
	return erc20Hook{k}
}

// IssueToken issues a new token
func (k Keeper) IssueToken(
	ctx sdk.Context,
	symbol string,
	name string,
	minUnit string,
	scale uint32,
	initialSupply uint64,
	maxSupply uint64,
	mintable bool,
	owner sdk.AccAddress,
) error {
	token := v1.NewToken(
		symbol, name, minUnit, scale, initialSupply,
		maxSupply, mintable, owner,
	)

	if err := k.AddToken(ctx, token, true); err != nil {
		return err
	}

	precision := sdkmath.NewIntWithDecimal(1, int(token.Scale))
	initialCoin := sdk.NewCoin(
		token.MinUnit,
		sdk.NewIntFromUint64(token.InitialSupply).Mul(precision),
	)

	mintCoins := sdk.NewCoins(initialCoin)

	// mint coins into module account
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, mintCoins); err != nil {
		return err
	}

	// sent coins to owner's account
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, owner, mintCoins)
}

// EditToken edits the specified token
func (k Keeper) EditToken(
	ctx sdk.Context,
	symbol string,
	name string,
	maxSupply uint64,
	mintable types.Bool,
	owner sdk.AccAddress,
) error {
	// get the destination token
	token, err := k.getTokenBySymbol(ctx, symbol)
	if err != nil {
		return err
	}

	if owner.String() != token.Owner {
		return errorsmod.Wrapf(
			types.ErrInvalidOwner,
			"the address %s is not the owner of the token %s",
			owner,
			symbol,
		)
	}

	if maxSupply > 0 {
		issuedAmt := k.getTokenSupply(ctx, token.MinUnit)
		issuedMainUnitAmt := issuedAmt.Quo(sdkmath.NewIntWithDecimal(1, int(token.Scale)))

		if sdk.NewIntFromUint64(maxSupply).LT(issuedMainUnitAmt) {
			return errorsmod.Wrapf(
				types.ErrInvalidMaxSupply,
				"max supply must not be less than %s",
				issuedMainUnitAmt,
			)
		}

		token.MaxSupply = maxSupply
	}

	if name != v1.DoNotModify {
		token.Name = name
		metadata, exist := k.bankKeeper.GetDenomMetaData(ctx, token.MinUnit)
		if exist {
			metadata.Description = name
			k.bankKeeper.SetDenomMetaData(ctx, metadata)
		}
	}

	if mintable != types.Nil {
		token.Mintable = mintable.ToBool()
	}

	k.setToken(ctx, token)

	return nil
}

// TransferTokenOwner transfers the owner of the specified token to a new one
func (k Keeper) TransferTokenOwner(
	ctx sdk.Context,
	symbol string,
	srcOwner sdk.AccAddress,
	dstOwner sdk.AccAddress,
) error {
	token, err := k.getTokenBySymbol(ctx, symbol)
	if err != nil {
		return err
	}

	if srcOwner.String() != token.Owner {
		return errorsmod.Wrapf(
			types.ErrInvalidOwner,
			"the address %s is not the owner of the token %s",
			srcOwner,
			symbol,
		)
	}
	return k.changeTokenOwner(ctx, token, dstOwner)
}

// UnsafeTransferTokenOwner transfer the token owner without authorization
// NOTE: this method should be used with caution
func (k Keeper) UnsafeTransferTokenOwner(ctx sdk.Context, symbol string, to sdk.AccAddress) error {
	token, err := k.getTokenBySymbol(ctx, symbol)
	if err != nil {
		return err
	}
	return k.changeTokenOwner(ctx, token, to)
}

// MintToken mints the specified amount of token to the specified recipient
// NOTE: empty owner means that the external caller is responsible to manage the token authority
func (k Keeper) MintToken(
	ctx sdk.Context,
	coinMinted sdk.Coin,
	recipient sdk.AccAddress,
	owner sdk.AccAddress,
) error {
	token, err := k.getTokenByMinUnit(ctx, coinMinted.Denom)
	if err != nil {
		return err
	}

	if owner.String() != token.Owner {
		return errorsmod.Wrapf(
			types.ErrInvalidOwner,
			"the address %s is not the owner of the token %s",
			owner,
			token.Symbol,
		)
	}

	if !token.Mintable {
		return errorsmod.Wrapf(types.ErrNotMintable, "%s", token.Symbol)
	}

	supply := k.getTokenSupply(ctx, token.MinUnit)
	precision := sdkmath.NewIntWithDecimal(1, int(token.Scale))
	mintableAmt := sdk.NewIntFromUint64(token.MaxSupply).Mul(precision).Sub(supply)

	if coinMinted.Amount.GT(mintableAmt) {
		return errorsmod.Wrapf(
			types.ErrInvalidAmount,
			"the amount exceeds the mintable token amount; expected (0, %d], got %d",
			mintableAmt, coinMinted.Amount,
		)
	}

	mintCoins := sdk.NewCoins(coinMinted)

	// mint coins
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, mintCoins); err != nil {
		return err
	}

	if recipient.Empty() {
		recipient = owner
	}

	// sent coins to the recipient account
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, mintCoins)
}

// BurnToken burns the specified amount of token
func (k Keeper) BurnToken(
	ctx sdk.Context,
	coinBurnt sdk.Coin,
	owner sdk.AccAddress,
) error {
	_, err := k.getTokenByMinUnit(ctx, coinBurnt.Denom)
	if err != nil {
		return err
	}

	burnCoins := sdk.NewCoins(coinBurnt)
	// burn coins
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, owner, types.ModuleName, burnCoins); err != nil {
		return err
	}

	k.AddBurnCoin(ctx, coinBurnt)

	return k.bankKeeper.BurnCoins(ctx, types.ModuleName, burnCoins)
}

// SwapFeeToken swap the fee token
func (k Keeper) SwapFeeToken(
	ctx sdk.Context,
	feePaid sdk.Coin,
	sender sdk.AccAddress,
	recipient sdk.AccAddress,
) (sdk.Coin, sdk.Coin, error) {
	burnedCoin, mintedCoin, err := k.calcFeeTokenMinted(ctx, feePaid)
	if err != nil {
		return sdk.Coin{}, sdk.Coin{}, err
	}

	burnedCoins := sdk.NewCoins(burnedCoin)
	// burn coins
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, burnedCoins); err != nil {
		return sdk.Coin{}, sdk.Coin{}, err
	}
	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, burnedCoins); err != nil {
		return sdk.Coin{}, sdk.Coin{}, err
	}

	// mint coins
	mintedCoins := sdk.NewCoins(mintedCoin)
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, mintedCoins); err != nil {
		return sdk.Coin{}, sdk.Coin{}, err
	}

	if recipient == nil {
		recipient = sender
	}
	return burnedCoin, mintedCoin, k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		recipient,
		mintedCoins,
	)
}

// WithSwapRegistry sets the swap registry in the Keeper and returns the updated Keeper instance.
//
// registry: The swap registry to set.
// Returns the updated Keeper instance.
func (k Keeper) WithSwapRegistry(registry v1.SwapRegistry) Keeper {
	k.registry = registry
	return k
}

func (k Keeper) calcFeeTokenMinted(
	ctx sdk.Context,
	feePaid sdk.Coin,
) (burnt, minted sdk.Coin, err error) {
	tokenBurned, err := k.getTokenByMinUnit(ctx, feePaid.Denom)
	if err != nil {
		return burnt, minted, err
	}

	swapParams, ok := k.registry[tokenBurned.GetMinUnit()]
	if !ok {
		return burnt, minted, types.ErrInvalidSwap
	}

	tokenMinted, err := k.GetToken(ctx, swapParams.MinUnit)
	if err != nil {
		return burnt, minted, err
	}

	burntAmt, mintAmt := types.LossLessSwap(
		feePaid.Amount,
		swapParams.Ratio,
		tokenBurned.GetScale(),
		tokenMinted.GetScale(),
	)
	return sdk.NewCoin(tokenBurned.MinUnit, burntAmt), sdk.NewCoin(swapParams.MinUnit, mintAmt), nil
}

func (k Keeper) changeTokenOwner(
	ctx sdk.Context,
	srcToken v1.Token,
	dstOwner sdk.AccAddress,
) error {
	srcOwner, err := sdk.AccAddressFromBech32(srcToken.Owner)
	if err != nil {
		return err
	}

	srcToken.Owner = dstOwner.String()
	// update token
	k.setToken(ctx, srcToken)

	// reset all indices
	k.resetStoreKeyForQueryToken(ctx, srcToken.Symbol, srcOwner, dstOwner)

	return nil
}
