package bank

import (
	"fmt"

	"github.com/irisnet/irishub/modules/auth"
	sdk "github.com/irisnet/irishub/types"
)

const (
	costGetCoins       sdk.Gas = 10
	costGetLoosenCoins sdk.Gas = 10
	costGetBurnedCoins sdk.Gas = 10
	costBurnCoins      sdk.Gas = 100
	costHasCoins       sdk.Gas = 10
	costSetCoins       sdk.Gas = 100
	costSubtractCoins  sdk.Gas = 10
	costAddCoins       sdk.Gas = 10
)

// Keeper defines a module interface that facilitates the transfer of coins
// between accounts.
type Keeper interface {
	SendKeeper
	IncreaseLoosenToken(ctx sdk.Context, amt sdk.Coins)
	DecreaseLoosenToken(ctx sdk.Context, amt sdk.Coins)
	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error)
	AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error)
	BurnCoinsFromAddr(ctx sdk.Context, fromAddr sdk.AccAddress, amt sdk.Coins) (sdk.Tags, sdk.Error)
	BurnCoinsFromPool(ctx sdk.Context, pool string, amt sdk.Coins) (sdk.Tags, sdk.Error)
}

var _ Keeper = (*BaseKeeper)(nil)

// BaseKeeper manages transfers between accounts. It implements the Keeper
// interface.
type BaseKeeper struct {
	am auth.AccountKeeper
}

// NewBaseKeeper returns a new BaseKeeper
func NewBaseKeeper(am auth.AccountKeeper) BaseKeeper {
	return BaseKeeper{am: am}
}

// GetCoins returns the coins at the addr.
func (keeper BaseKeeper) GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return getCoins(ctx, keeper.am, addr)
}

// GetLoosenCoins returns the total loosen coins
func (keeper BaseKeeper) GetLoosenCoins(ctx sdk.Context) sdk.Coins {
	return getLoosenCoins(ctx, keeper.am)
}

// GetLoosenCoins returns the burned coins
func (keeper BaseKeeper) GetBurnedCoins(ctx sdk.Context) sdk.Coins {
	return getBurnedCoins(ctx, keeper.am)
}

// HasCoins returns whether or not an account has at least amt coins.
func (keeper BaseKeeper) HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool {
	return hasCoins(ctx, keeper.am, addr, amt)
}

// SubtractCoins subtracts amt from the coins at the addr.
func (keeper BaseKeeper) SubtractCoins(
	ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins,
) (sdk.Coins, sdk.Tags, sdk.Error) {

	return subtractCoins(ctx, keeper.am, addr, amt)
}

// AddCoins adds amt to the coins at the addr.
func (keeper BaseKeeper) AddCoins(
	ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins,
) (sdk.Coins, sdk.Tags, sdk.Error) {

	return addCoins(ctx, keeper.am, addr, amt)
}

// SendCoins moves coins from one account to another
func (keeper BaseKeeper) SendCoins(
	ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins,
) (sdk.Tags, sdk.Error) {

	return sendCoins(ctx, keeper.am, fromAddr, toAddr, amt)
}

func (keeper BaseKeeper) IncreaseLoosenToken(
	ctx sdk.Context, amt sdk.Coins) {
	keeper.am.IncreaseTotalLoosenToken(ctx, amt)
}

// SendCoins moves coins from one account to another
func (keeper BaseKeeper) DecreaseLoosenToken(
	ctx sdk.Context, amt sdk.Coins) {
	keeper.am.DecreaseTotalLoosenToken(ctx, amt)
}

// BurnCoins burns coins from one account
func (keeper BaseKeeper) BurnCoinsFromAddr(
	ctx sdk.Context, fromAddr sdk.AccAddress, amt sdk.Coins,
) (sdk.Tags, sdk.Error) {
	_, _, err := subtractCoins(ctx, keeper.am, fromAddr, amt)
	if err != nil {
		return nil, err
	}
	return burnCoins(ctx, keeper.am, fromAddr.String(), amt)
}

// BurnCoins burns coins from one account
func (keeper BaseKeeper) BurnCoinsFromPool(
	ctx sdk.Context, pool string, amt sdk.Coins,
) (sdk.Tags, sdk.Error) {
	return burnCoins(ctx, keeper.am, pool, amt)
}

// InputOutputCoins handles a list of inputs and outputs
func (keeper BaseKeeper) InputOutputCoins(ctx sdk.Context, inputs []Input, outputs []Output) (sdk.Tags, sdk.Error) {
	return inputOutputCoins(ctx, keeper.am, inputs, outputs)
}

//______________________________________________________________________________________________

// SendKeeper defines a module interface that facilitates the transfer of coins
// between accounts without the possibility of creating coins.
type SendKeeper interface {
	ViewKeeper
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) (sdk.Tags, sdk.Error)
	InputOutputCoins(ctx sdk.Context, inputs []Input, outputs []Output) (sdk.Tags, sdk.Error)
}

var _ SendKeeper = (*BaseSendKeeper)(nil)

// SendKeeper only allows transfers between accounts without the possibility of
// creating coins. It implements the SendKeeper interface.
type BaseSendKeeper struct {
	am auth.AccountKeeper
}

// NewBaseSendKeeper returns a new BaseSendKeeper.
func NewBaseSendKeeper(am auth.AccountKeeper) BaseSendKeeper {
	return BaseSendKeeper{am: am}
}

// GetCoins returns the coins at the addr.
func (keeper BaseSendKeeper) GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return getCoins(ctx, keeper.am, addr)
}

// GetLoosenCoins returns the total loosen coins
func (keeper BaseSendKeeper) GetLoosenCoins(ctx sdk.Context) sdk.Coins {
	return getLoosenCoins(ctx, keeper.am)
}

// GetLoosenCoins returns the burned coins
func (keeper BaseSendKeeper) GetBurnedCoins(ctx sdk.Context) sdk.Coins {
	return getBurnedCoins(ctx, keeper.am)
}

// HasCoins returns whether or not an account has at least amt coins.
func (keeper BaseSendKeeper) HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool {
	return hasCoins(ctx, keeper.am, addr, amt)
}

// SendCoins moves coins from one account to another
func (keeper BaseSendKeeper) SendCoins(
	ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins,
) (sdk.Tags, sdk.Error) {

	return sendCoins(ctx, keeper.am, fromAddr, toAddr, amt)
}

// InputOutputCoins handles a list of inputs and outputs
func (keeper BaseSendKeeper) InputOutputCoins(
	ctx sdk.Context, inputs []Input, outputs []Output,
) (sdk.Tags, sdk.Error) {

	return inputOutputCoins(ctx, keeper.am, inputs, outputs)
}

//______________________________________________________________________________________________

// ViewKeeper defines a module interface that facilitates read only access to
// account balances.
type ViewKeeper interface {
	GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetLoosenCoins(ctx sdk.Context) sdk.Coins
	GetBurnedCoins(ctx sdk.Context) sdk.Coins
	HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool
}

var _ ViewKeeper = (*BaseViewKeeper)(nil)

// BaseViewKeeper implements a read only keeper implementation of ViewKeeper.
type BaseViewKeeper struct {
	am auth.AccountKeeper
}

// NewBaseViewKeeper returns a new BaseViewKeeper.
func NewBaseViewKeeper(am auth.AccountKeeper) BaseViewKeeper {
	return BaseViewKeeper{am: am}
}

// GetCoins returns the coins at the addr.
func (keeper BaseViewKeeper) GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return getCoins(ctx, keeper.am, addr)
}

// GetLoosenCoins returns the total loosen coins
func (keeper BaseViewKeeper) GetLoosenCoins(ctx sdk.Context) sdk.Coins {
	return getLoosenCoins(ctx, keeper.am)
}

// GetLoosenCoins returns the burned coins
func (keeper BaseViewKeeper) GetBurnedCoins(ctx sdk.Context) sdk.Coins {
	return getBurnedCoins(ctx, keeper.am)
}

// HasCoins returns whether or not an account has at least amt coins.
func (keeper BaseViewKeeper) HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool {
	return hasCoins(ctx, keeper.am, addr, amt)
}

//______________________________________________________________________________________________

func getCoins(ctx sdk.Context, am auth.AccountKeeper, addr sdk.AccAddress) sdk.Coins {
	ctx.GasMeter().ConsumeGas(costGetCoins, "getCoins")
	acc := am.GetAccount(ctx, addr)
	if acc == nil {
		return sdk.Coins{}
	}
	return acc.GetCoins()
}

func getLoosenCoins(ctx sdk.Context, am auth.AccountKeeper) sdk.Coins {
	ctx.GasMeter().ConsumeGas(costGetLoosenCoins, "getLoosenCoins")
	return am.GetTotalLoosenToken(ctx)
}

func getBurnedCoins(ctx sdk.Context, am auth.AccountKeeper) sdk.Coins {
	ctx.GasMeter().ConsumeGas(costGetBurnedCoins, "getBurnedCoins")
	return am.GetBurnedToken(ctx)
}

func setCoins(ctx sdk.Context, am auth.AccountKeeper, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	ctx.GasMeter().ConsumeGas(costSetCoins, "setCoins")
	acc := am.GetAccount(ctx, addr)
	if acc == nil {
		acc = am.NewAccountWithAddress(ctx, addr)
	}
	err := acc.SetCoins(amt)
	if err != nil {
		// Handle w/ #870
		panic(err)
	}
	am.SetAccount(ctx, acc)
	return nil
}

// HasCoins returns whether or not an account has at least amt coins.
func hasCoins(ctx sdk.Context, am auth.AccountKeeper, addr sdk.AccAddress, amt sdk.Coins) bool {
	ctx.GasMeter().ConsumeGas(costHasCoins, "hasCoins")
	return getCoins(ctx, am, addr).IsAllGTE(amt)
}

// SubtractCoins subtracts amt from the coins at the addr.
func subtractCoins(ctx sdk.Context, am auth.AccountKeeper, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error) {
	ctx.GasMeter().ConsumeGas(costSubtractCoins, "subtractCoins")
	oldCoins := getCoins(ctx, am, addr)
	newCoins, hasNeg := oldCoins.SafeMinus(amt)
	if hasNeg {
		return amt, nil, sdk.ErrInsufficientCoins(fmt.Sprintf("%s is less than %s", oldCoins, amt))
	}
	err := setCoins(ctx, am, addr, newCoins)
	tags := sdk.NewTags("sender", []byte(addr.String()))
	return newCoins, tags, err
}

// AddCoins adds amt to the coins at the addr.
func addCoins(ctx sdk.Context, am auth.AccountKeeper, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error) {
	ctx.GasMeter().ConsumeGas(costAddCoins, "addCoins")
	oldCoins := getCoins(ctx, am, addr)
	newCoins := oldCoins.Plus(amt)
	if !newCoins.IsNotNegative() {
		return amt, nil, sdk.ErrInsufficientCoins(fmt.Sprintf("%s is less than %s", oldCoins, amt))
	}
	err := setCoins(ctx, am, addr, newCoins)
	tags := sdk.NewTags("recipient", []byte(addr.String()))
	return newCoins, tags, err
}

// SendCoins moves coins from one account to another
// NOTE: Make sure to revert state changes from tx on error
func sendCoins(ctx sdk.Context, am auth.AccountKeeper, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) (sdk.Tags, sdk.Error) {
	_, subTags, err := subtractCoins(ctx, am, fromAddr, amt)
	if err != nil {
		return nil, err
	}

	_, addTags, err := addCoins(ctx, am, toAddr, amt)
	if err != nil {
		return nil, err
	}

	return subTags.AppendTags(addTags), nil
}

// burnCoins moves coins from burn address
// NOTE: Make sure to revert state changes from tx on error
func burnCoins(ctx sdk.Context, am auth.AccountKeeper, from string, amt sdk.Coins) (sdk.Tags, sdk.Error) {
	ctx.GasMeter().ConsumeGas(costBurnCoins, "burnCoins")
	am.DecreaseTotalLoosenToken(ctx, amt)
	am.IncreaseBurnedToken(ctx, amt)
	burnTags := sdk.NewTags(
		"burnFrom", []byte(from),
		"burnAmount", []byte(amt.String()),
	)

	return burnTags, nil
}

// InputOutputCoins handles a list of inputs and outputs
// NOTE: Make sure to revert state changes from tx on error
func inputOutputCoins(ctx sdk.Context, am auth.AccountKeeper, inputs []Input, outputs []Output) (sdk.Tags, sdk.Error) {
	allTags := sdk.EmptyTags()

	for _, in := range inputs {
		_, tags, err := subtractCoins(ctx, am, in.Address, in.Coins)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(tags)
	}

	for _, out := range outputs {
		_, tags, err := addCoins(ctx, am, out.Address, out.Coins)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(tags)
	}

	return allTags, nil
}
