package bank

import (
	"fmt"

	"github.com/irisnet/irishub/app/v1/auth"
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
	costGetFrozenCoin  sdk.Gas = 10
	costFreezeCoin     sdk.Gas = 10
	costUnfreezeCoin   sdk.Gas = 10
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
	FreezeCoinFromAddr(ctx sdk.Context, Addr sdk.AccAddress, amt sdk.Coin) (sdk.Tags, sdk.Error)
	UnfreezeCoinFromAddr(ctx sdk.Context, Addr sdk.AccAddress,  amt sdk.Coin) (sdk.Tags, sdk.Error)
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

//provide interface to search the total frozen token for specific coin
func (keeper BaseKeeper) GetFrozenCoin(ctx sdk.Context, am auth.AccountKeeper, addr sdk.AccAddress, denom string) sdk.Coin {
	ctx.GasMeter().ConsumeGas(costGetFrozenCoin, "getfrozenCoin")
	acc := am.GetAccount(ctx, addr)
	if acc == nil {
		return sdk.Coin{}
	}
	return acc.GetFrozenCoinByDenom(denom)
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

// FreezeCoinFromAddr freezes coins from one account
func (keeper BaseKeeper) FreezeCoinFromAddr(
	ctx sdk.Context, Addr sdk.AccAddress, amt sdk.Coin,
) (sdk.Tags, sdk.Error) {
	_, _, err := subtractCoins(ctx, keeper.am, Addr, sdk.Coins{amt})
	if err != nil {
		return nil, err
	}

	_, err = increaseFrozenCoin(ctx, keeper.am, Addr, amt)
	if err != nil {
		return nil, err
	}

	//send the frozen token to specific account
	//if msg.owner=msg.holder, then send the token to account.frozen coins
	//if msg.owner=msg.holder, then move the token to msg.owner account.frozen coins


	return freezeCoin(ctx, keeper.am, Addr.String(), amt)
}

// UnfreezeCoinFromAddr unfreezes coins from one account
func (keeper BaseKeeper) UnfreezeCoinFromAddr(
	ctx sdk.Context, Addr sdk.AccAddress, amt sdk.Coin,
) (sdk.Tags, sdk.Error) {

	_, err := decreaseFrozenCoin(ctx, keeper.am, Addr, amt)
	if err != nil {
		return nil, err
	}
	_, _, err = addCoins(ctx, keeper.am, Addr, sdk.Coins{amt})
	if err != nil {
		return nil, err
	}
	return unfreezeCoin(ctx, keeper.am, Addr.String(), amt)
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

//set the frozen coins for holders.the frozen coin will be stored in othFrozenTokens or FrozenTokens
func setFrozenCoin(ctx sdk.Context, am auth.AccountKeeper, Addr sdk.AccAddress, amt sdk.Coin) sdk.Error {
	ctx.GasMeter().ConsumeGas(costFreezeCoin, "setFrozenCoins")
	acc := am.GetAccount(ctx, Addr)
	if acc == nil {
		acc = am.NewAccountWithAddress(ctx, Addr)
	}

	err := acc.SetFrozenCoin(amt)
	if err != nil {
		// Handle w/ #870
		panic(err)
	}

	am.SetAccount(ctx, acc)
	return nil
}

func deductFrozenCoin(ctx sdk.Context, am auth.AccountKeeper, addr sdk.AccAddress, amt sdk.Coin) sdk.Error {
	ctx.GasMeter().ConsumeGas(costFreezeCoin, "setFrozenCoins")
	acc := am.GetAccount(ctx, addr)
	if acc == nil {
		acc = am.NewAccountWithAddress(ctx, addr)
	}
	var err error

	err = acc.DeductFrozenCoin(amt)

	if err != nil {
		// Handle w/ #870
		panic(err)
		//return err
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

//increaseFrozenCoins increases the frozen the coins at addr.
func increaseFrozenCoin(ctx sdk.Context, am auth.AccountKeeper, Addr sdk.AccAddress, amt sdk.Coin) (sdk.Tags, sdk.Error) {
	ctx.GasMeter().ConsumeGas(costFreezeCoin, "frozenCoins")
	var err sdk.Error
	var tags sdk.Tags
	err = setFrozenCoin(ctx, am, Addr, amt)

	tags = sdk.NewTags("sender", []byte(Addr.String()))
	return tags, err
}

//decreaseFrozenCoins increases the frozen the coins at addr.
func decreaseFrozenCoin(ctx sdk.Context, am auth.AccountKeeper, addr sdk.AccAddress, amt sdk.Coin) (sdk.Tags, sdk.Error) {
	ctx.GasMeter().ConsumeGas(costUnfreezeCoin, "unfrozenCoins")
	var err sdk.Error
	err = deductFrozenCoin(ctx, am, addr, amt)
	tags := sdk.NewTags("sender", []byte(addr.String()))
	return tags, err
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

	ctx.Logger().Info("Execute Burntoken Successed", "burnFrom", from, "burnAmount", amt.String())

	return burnTags, nil
}

//getTotalFrozenToken get the total frozen token for specific coin
func getTotalFrozenToken(ctx sdk.Context, am auth.AccountKeeper, denom string) sdk.Coin {
	return am.GetFrozenToken(ctx, []byte(denom))
}

// freezeCoin moves coins to frozen token from
// NOTE: Make sure to revert state changes from tx on error
func freezeCoin(ctx sdk.Context, am auth.AccountKeeper, from string, amt sdk.Coin) (sdk.Tags, sdk.Error) {
	ctx.GasMeter().ConsumeGas(costFreezeCoin, "freezeCoins")

	am.IncreaseFrozenToken(ctx, amt)
	freezeTags := sdk.NewTags(
		"freezeFrom", []byte(from),
		"freezeAmount", []byte(amt.String()),
	)

	ctx.Logger().Info("Execute Frozentoken Successed", "freezeFrom", from, "freezeAmount", amt.String())

	return freezeTags, nil
}

// unfreezeCoins add moves frozen token to coins
// NOTE: Make sure to revert state changes from tx on error
func unfreezeCoin(ctx sdk.Context, am auth.AccountKeeper, from string, amt sdk.Coin) (sdk.Tags, sdk.Error) {
	ctx.GasMeter().ConsumeGas(costFreezeCoin, "unfreezeCoins")

	am.DecreaseFrozenToken(ctx, amt)
	unfreezeTags := sdk.NewTags(
		"unfreezeFrom", []byte(from),
		"unfreezeAmount", []byte(amt.String()),
	)

	ctx.Logger().Info("Execute Unfrozentoken Successed", "unfreezeFrom", from, "unfreezeAmount", amt.String())

	return unfreezeTags, nil
}

// InputOutputCoins handles a list of inputs and outputs
// NOTE: Make sure to revert state changes from tx on error
func inputOutputCoins(ctx sdk.Context, am auth.AccountKeeper, inputs []Input, outputs []Output) (sdk.Tags, sdk.Error) {
	allTags := sdk.EmptyTags()

	multiInMultiOut := true
	if len(inputs) == 1 && len(outputs) == 1 {
		multiInMultiOut = false
		ctx.CoinFlowTags().AppendCoinFlowTag(ctx, inputs[0].Address.String(), outputs[0].Address.String(), inputs[0].Coins.String(), sdk.TransferFlow, "")
	}

	for _, in := range inputs {
		_, tags, err := subtractCoins(ctx, am, in.Address, in.Coins)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(tags)
		if multiInMultiOut {
			ctx.CoinFlowTags().AppendCoinFlowTag(ctx, in.Address.String(), ctx.CoinFlowTrigger(), in.Coins.String(), sdk.TransferFlow, "")
		}
	}

	for _, out := range outputs {
		_, tags, err := addCoins(ctx, am, out.Address, out.Coins)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(tags)
		if multiInMultiOut {
			ctx.CoinFlowTags().AppendCoinFlowTag(ctx, ctx.CoinFlowTrigger(), out.Address.String(), out.Coins.String(), sdk.TransferFlow, "")
		}
	}

	return allTags, nil
}
