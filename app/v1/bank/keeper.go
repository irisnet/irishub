package bank

import (
	"fmt"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/codec"
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

//var BurnedCoinsAccAddr = sdk.AccAddress(crypto.AddressHash([]byte("burnedCoins")))

// Keeper defines a module interface that facilitates the transfer of coins
// between accounts.
type Keeper interface {
	SendKeeper
	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error)
	AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error)
	BurnCoins(ctx sdk.Context, fromAddr sdk.AccAddress, amt sdk.Coins) (sdk.Tags, sdk.Error)
	SetMemoRegexp(ctx sdk.Context, fromAddr sdk.AccAddress, regexp string) (sdk.Tags, sdk.Error)
	IncreaseLoosenToken(ctx sdk.Context, amt sdk.Coins)
	DecreaseLoosenToken(ctx sdk.Context, amt sdk.Coins)
	IncreaseTotalSupply(ctx sdk.Context, amt sdk.Coin) sdk.Error
	DecreaseTotalSupply(ctx sdk.Context, amt sdk.Coin) sdk.Error
	SetTotalSupply(ctx sdk.Context, totalSupply sdk.Coin)
}

var _ Keeper = (*BaseKeeper)(nil)

// BaseKeeper manages transfers between accounts. It implements the Keeper
// interface.
type BaseKeeper struct {
	am  auth.AccountKeeper
	cdc *codec.Codec
}

func (keeper BaseKeeper) GetTotalSupply(ctx sdk.Context, denom string) (coin sdk.Coin, found bool) {
	return keeper.am.GetTotalSupply(ctx, denom)
}

func (keeper BaseKeeper) GetTotalSupplies(ctx sdk.Context) sdk.Iterator {
	return keeper.am.GetTotalSupplies(ctx)
}

// NewBaseKeeper returns a new BaseKeeper
func NewBaseKeeper(cdc *codec.Codec, am auth.AccountKeeper) BaseKeeper {
	return BaseKeeper{am: am, cdc: cdc}
}

// GetCoins returns the coins at the addr.
func (keeper BaseKeeper) GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return getCoins(ctx, keeper.am, addr)
}

// GetLoosenCoins returns the total loosen coins
func (keeper BaseKeeper) GetLoosenCoins(ctx sdk.Context) sdk.Coins {
	return getLoosenCoins(ctx, keeper.am)
}

// GetBurnedCoins returns the burned coins
//func (keeper BaseKeeper) GetBurnedCoins(ctx sdk.Context) sdk.Coins {
//	return getBurnedCoins(ctx, keeper.am)
//}

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

// SetMemoRegexp sets memo regexp for sender account
func (keeper BaseKeeper) SetMemoRegexp(ctx sdk.Context, fromAddr sdk.AccAddress, regexp string) (sdk.Tags, sdk.Error) {
	acc := keeper.am.GetAccount(ctx, fromAddr)
	if acc == nil {
		return nil, ErrInvalidAccount(DefaultCodespace, "account does not exist")
	}

	acc.SetMemoRegexp(regexp)
	keeper.am.SetAccount(ctx, acc)

	tags := sdk.NewTags(
		"memoRegexp", []byte(regexp),
	)

	return tags, nil
}

// BurnCoins burns coins from the given account
func (keeper BaseKeeper) BurnCoins(ctx sdk.Context, fromAddr sdk.AccAddress, amt sdk.Coins) (sdk.Tags, sdk.Error) {
	ctx.GasMeter().ConsumeGas(costBurnCoins, "burnCoins")

	if _, err := keeper.SendCoins(ctx, fromAddr, auth.BurnedCoinsAccAddr, amt); err != nil {
		return nil, err
	}

	burnTags := sdk.NewTags(
		"burnAmount", []byte(amt.String()),
	)
	ctx.Logger().Info("Tokens Burned Successfully", "burnFrom", fromAddr, "burnAmount", amt.String())

	return burnTags, nil
}

func (keeper BaseKeeper) IncreaseTotalSupply(ctx sdk.Context, coin sdk.Coin) sdk.Error {
	return keeper.am.IncreaseTotalSupply(ctx, coin)
}

func (keeper BaseKeeper) DecreaseTotalSupply(ctx sdk.Context, coin sdk.Coin) sdk.Error {
	return keeper.am.DecreaseTotalSupply(ctx, coin)
}

func (keeper BaseKeeper) SetTotalSupply(ctx sdk.Context, totalSupply sdk.Coin) {
	keeper.am.SetTotalSupply(ctx, totalSupply)
}

// InputOutputCoins handles a list of inputs and outputs
func (keeper BaseKeeper) InputOutputCoins(ctx sdk.Context, inputs []Input, outputs []Output) (sdk.Tags, sdk.Error) {
	return inputOutputCoins(ctx, keeper.am, inputs, outputs)
}

func (keeper BaseKeeper) Init(ctx sdk.Context) {
	var burnedCoins sdk.Coins
	store := ctx.KVStore(protocol.KeyAccount)
	bz := store.Get([]byte("burnedToken"))
	if bz == nil {
		burnedCoins = sdk.Coins{}
	} else {
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &burnedCoins)
	}

	// If someone send coins to BurnedCoinsAccAddr in protocol v0,
	// It is handled as burned coins in protocol v1.
	// Subtract all BurnedCoinsAccAddr coins in v0 from looseToken.
	v0burnedAddrCoins := keeper.GetCoins(ctx, auth.BurnedCoinsAccAddr)
	if !v0burnedAddrCoins.Empty() {
		keeper.DecreaseLoosenToken(ctx, v0burnedAddrCoins)
	}

	keeper.IncreaseLoosenToken(ctx, burnedCoins)
	keeper.AddCoins(ctx, auth.BurnedCoinsAccAddr, burnedCoins)
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

// GetBurnedCoins returns the burned coins
//func (keeper BaseSendKeeper) GetBurnedCoins(ctx sdk.Context) sdk.Coins {
//	return getBurnedCoins(ctx, keeper.am)
//}

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

func (keeper BaseSendKeeper) GetTotalSupply(ctx sdk.Context, denom string) (coin sdk.Coin, found bool) {
	return keeper.am.GetTotalSupply(ctx, denom)
}

//______________________________________________________________________________________________

// ViewKeeper defines a module interface that facilitates read only access to
// account balances.
type ViewKeeper interface {
	GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetLoosenCoins(ctx sdk.Context) sdk.Coins
	//GetBurnedCoins(ctx sdk.Context) sdk.Coins
	HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool
	GetTotalSupply(ctx sdk.Context, denom string) (coin sdk.Coin, found bool)
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

// GetBurnedCoins returns the burned coins
//func (keeper BaseViewKeeper) GetBurnedCoins(ctx sdk.Context) sdk.Coins {
//	return getBurnedCoins(ctx, keeper.am)
//}

// HasCoins returns whether or not an account has at least amt coins.
func (keeper BaseViewKeeper) HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool {
	return hasCoins(ctx, keeper.am, addr, amt)
}

func (keeper BaseViewKeeper) GetTotalSupply(ctx sdk.Context, denom string) (coin sdk.Coin, found bool) {
	return keeper.am.GetTotalSupply(ctx, denom)
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

//func getBurnedCoins(ctx sdk.Context, am auth.AccountKeeper) sdk.Coins {
//	ctx.GasMeter().ConsumeGas(costGetBurnedCoins, "getBurnedCoins")
//	return am.GetBurnedToken(ctx)
//}

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
	if !amt.IsValid() {
		panic(fmt.Sprintf("invalid coins [%s]", amt))
	}
	ctx.GasMeter().ConsumeGas(costSubtractCoins, "subtractCoins")
	oldCoins := getCoins(ctx, am, addr)
	newCoins, hasNeg := oldCoins.SafeSub(amt)
	if hasNeg {
		return amt, nil, sdk.ErrInsufficientCoins(fmt.Sprintf("subtracting [%s] from [%s] yields negative coin(s)", amt, oldCoins))
	}
	err := setCoins(ctx, am, addr, newCoins)
	tags := sdk.NewTags("sender", []byte(addr.String()))
	return newCoins, tags, err
}

// AddCoins adds amt to the coins at the addr.
func addCoins(ctx sdk.Context, am auth.AccountKeeper, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error) {
	if !amt.IsValid() {
		panic(fmt.Sprintf("invalid coins [%s]", amt))
	}
	ctx.GasMeter().ConsumeGas(costAddCoins, "addCoins")
	oldCoins := getCoins(ctx, am, addr)
	newCoins := oldCoins.Add(amt)
	err := setCoins(ctx, am, addr, newCoins)

	// adding coins to BurnedCoinsAccAddr is equivalent to burning coins
	if addr.Equals(auth.BurnedCoinsAccAddr) {
		for _, coin := range amt {
			if coin.Denom == sdk.IrisAtto {
				// Decrease total loose token for iris
				am.DecreaseTotalLoosenToken(ctx, sdk.Coins{coin})
			} else {
				// Decrease total supply for tokens other than iris
				am.DecreaseTotalSupply(ctx, coin)
			}
		}
	}

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
//func burnCoins(ctx sdk.Context, am auth.AccountKeeper, from string, amt sdk.Coins) (sdk.Tags, sdk.Error) {
//	ctx.GasMeter().ConsumeGas(costBurnCoins, "burnCoins")
//	am.DecreaseTotalLoosenToken(ctx, amt)
//	am.IncreaseBurnedToken(ctx, amt)
//	burnTags := sdk.NewTags(
//		"burnFrom", []byte(from),
//		"burnAmount", []byte(amt.String()),
//	)
//
//	ctx.Logger().Info("Execute Burntoken Successed", "burnFrom", from, "burnAmount", amt.String())
//
//	return burnTags, nil
//}

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
