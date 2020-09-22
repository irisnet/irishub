package keeper

import (
	"fmt"
	"strings"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/irisnet/irismod/modules/token/types"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      codec.Marshaler

	// The bankKeeper to reduce the supply of the network
	bankKeeper types.BankKeeper

	feeCollectorName string

	// params subspace
	paramSpace paramstypes.Subspace
}

func NewKeeper(
	cdc codec.Marshaler, key sdk.StoreKey, paramSpace paramstypes.Subspace,
	bankKeeper types.BankKeeper, feeCollectorName string,
) Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:         key,
		cdc:              cdc,
		paramSpace:       paramSpace,
		bankKeeper:       bankKeeper,
		feeCollectorName: feeCollectorName,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("irismod/%s", types.ModuleName))
}

// IssueToken issues a new token
func (k Keeper) IssueToken(ctx sdk.Context, msg types.MsgIssueToken) error {
	token := types.NewToken(
		msg.Symbol, msg.Name, msg.MinUnit, msg.Scale, msg.InitialSupply,
		msg.MaxSupply, msg.Mintable, msg.Owner,
	)

	if err := k.AddToken(ctx, token); err != nil {
		return err
	}

	initialSupply := sdk.NewCoin(
		token.MinUnit,
		sdk.NewIntWithDecimal(int64(msg.InitialSupply), int(msg.Scale)),
	)

	mintCoins := sdk.NewCoins(initialSupply)

	// mint coins into module account
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, mintCoins); err != nil {
		return err
	}

	// sent coins to owner's account
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, token.Owner, mintCoins)
}

// EditToken edits the specified token
func (k Keeper) EditToken(ctx sdk.Context, msg types.MsgEditToken) error {
	// get the destination token
	tokenI, err := k.GetToken(ctx, msg.Symbol)
	if err != nil {
		return err
	}

	token := tokenI.(*types.Token)

	if !msg.Owner.Equals(token.Owner) {
		return sdkerrors.Wrapf(types.ErrInvalidOwner, "the address %d is not the owner of the token %s", msg.Owner, msg.Symbol)
	}

	if msg.MaxSupply > 0 {
		issuedAmt := k.getTokenSupply(ctx, token.MinUnit)
		issuedMainUnitAmt := uint64(issuedAmt.Quo(sdk.NewIntWithDecimal(1, int(token.Scale))).Int64())
		if msg.MaxSupply < issuedMainUnitAmt {
			return sdkerrors.Wrapf(types.ErrInvalidMaxSupply, "max supply must not be less than %d", issuedMainUnitAmt)
		}

		token.MaxSupply = msg.MaxSupply
	}

	if msg.Name != types.DoNotModify {
		token.Name = strings.TrimSpace(msg.Name)
	}

	if msg.Mintable != types.Nil {
		token.Mintable = msg.Mintable.ToBool()
	}

	return k.setToken(ctx, *token)
}

// TransferTokenOwner transfers the owner of the specified token to a new one
func (k Keeper) TransferTokenOwner(ctx sdk.Context, msg types.MsgTransferTokenOwner) error {
	tokenI, err := k.GetToken(ctx, msg.Symbol)
	if err != nil {
		return err
	}

	token := tokenI.(*types.Token)

	if !msg.SrcOwner.Equals(token.Owner) {
		return sdkerrors.Wrapf(types.ErrInvalidOwner, "the address %s is not the owner of the token %s", msg.SrcOwner.String(), msg.Symbol)
	}

	token.Owner = msg.DstOwner
	// update token information
	if err := k.setToken(ctx, *token); err != nil {
		return err
	}

	// reset all index for query-token
	return k.resetStoreKeyForQueryToken(ctx, msg, *token)
}

// MintToken mints specified amount token to a specified owner
func (k Keeper) MintToken(ctx sdk.Context, msg types.MsgMintToken) error {
	tokenI, err := k.GetToken(ctx, msg.Symbol)
	if err != nil {
		return err
	}

	token := tokenI.(*types.Token)

	if !msg.Owner.Equals(token.Owner) {
		return sdkerrors.Wrapf(types.ErrInvalidOwner, "the address %s is not the owner of the token %s", msg.Owner.String(), msg.Symbol)
	}

	if !token.Mintable {
		return sdkerrors.Wrapf(types.ErrNotMintable, "the token %s is set to be non-mintable", msg.Symbol)
	}

	issuedAmt := k.getTokenSupply(ctx, token.MinUnit)
	mintableMaxAmt := sdk.NewIntWithDecimal(int64(token.MaxSupply), int(token.Scale)).Sub(issuedAmt)
	mintableMaxMainUnitAmt := uint64(mintableMaxAmt.Quo(sdk.NewIntWithDecimal(1, int(token.Scale))).Int64())

	if msg.Amount > mintableMaxMainUnitAmt {
		return sdkerrors.Wrapf(
			types.ErrInvalidMaxSupply,
			"The amount of minting tokens plus the total amount of issued tokens has exceeded the maximum supply, only accepts amount (0, %d]",
			mintableMaxMainUnitAmt,
		)
	}

	mintCoin := sdk.NewCoin(token.MinUnit, sdk.NewIntWithDecimal(int64(msg.Amount), int(token.Scale)))
	mintCoins := sdk.NewCoins(mintCoin)

	// mint coins
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, mintCoins); err != nil {
		return err
	}

	mintAcc := msg.To
	if mintAcc.Empty() {
		mintAcc = token.Owner
	}

	// sent coins to owner's account
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, mintAcc, mintCoins)
}
