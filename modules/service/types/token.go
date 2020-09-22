package types

import (
	"math"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// TokenI defines the interface for Token
type TokenI interface {
	GetMinUnit() string
	GetScale() uint32
	ToMainCoin(coin sdk.Coin) (sdk.DecCoin, error)
	ToMinCoin(coin sdk.DecCoin) (sdk.Coin, error)
}

// MockToken represents a mock implementation for TokenI
type MockToken struct {
	Symbol  string
	MinUnit string
	Scale   uint32
}

func (token MockToken) ToMainCoin(coin sdk.Coin) (sdk.DecCoin, error) {
	if token.Symbol != coin.Denom && token.MinUnit != coin.Denom {
		return sdk.NewDecCoinFromDec(coin.Denom, sdk.ZeroDec()), sdkerrors.Wrapf(ErrInvalidPricing, "token not match")
	}

	if token.Symbol == coin.Denom {
		return sdk.NewDecCoin(coin.Denom, coin.Amount), nil
	}

	precision := math.Pow10(int(token.Scale))
	precisionStr := strconv.FormatFloat(precision, 'f', 0, 64)
	precisionDec, err := sdk.NewDecFromStr(precisionStr)
	if err != nil {
		return sdk.DecCoin{}, err
	}

	// dest amount = src amount / 10^(scale)
	amount := sdk.NewDecFromInt(coin.Amount).Quo(precisionDec)

	return sdk.NewDecCoinFromDec(token.Symbol, amount), nil
}

func (token MockToken) ToMinCoin(coin sdk.DecCoin) (sdk.Coin, error) {
	if token.Symbol != coin.Denom && token.MinUnit != coin.Denom {
		return sdk.NewCoin(coin.Denom, sdk.ZeroInt()), sdkerrors.Wrapf(ErrInvalidPricing, "token not match")
	}

	if token.MinUnit == coin.Denom {
		return sdk.NewCoin(coin.Denom, coin.Amount.TruncateInt()), nil
	}

	precision := math.Pow10(int(token.Scale))
	precisionStr := strconv.FormatFloat(precision, 'f', 0, 64)
	precisionDec, err := sdk.NewDecFromStr(precisionStr)
	if err != nil {
		return sdk.Coin{}, err
	}

	// dest amount = src amount * 10^(dest scale)
	amount := coin.Amount.Mul(precisionDec)

	return sdk.NewCoin(token.MinUnit, amount.TruncateInt()), nil
}

// GetSymbol gets the symbol
func (token MockToken) GetSymbol() string {
	return token.Symbol
}

// GetMinUnit gets the min unit
func (token MockToken) GetMinUnit() string {
	return token.MinUnit
}

// GetScale gets the scale
func (token MockToken) GetScale() uint32 {
	return token.Scale
}
