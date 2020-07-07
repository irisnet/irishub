package types

import (
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	//1 iris = 10^18 iris-atto
	Atto      = "atto"
	AttoScale = 18
)

var IrisCoinType = CoinType{
	Name:     Iris,
	MainUnit: NewUnit(Iris, 0),
	MinUnit:  NewUnit(fmt.Sprintf("%s-%s", Iris, Atto), AttoScale),
	Desc:     "IRIS Network",
}

type Unit struct {
	Denom string `json:"denom"` //denom of unit
	Scale uint8  `json:"scale"` //scale of unit
}

func NewUnit(denom string, scale uint8) Unit {
	return Unit{
		Denom: denom,
		Scale: scale,
	}
}

//GetScaleFactor return 1 * 10^scale
func (u Unit) GetScaleFactor() sdk.Int {
	return sdk.NewIntWithDecimal(1, int(u.Scale))
}

type CoinType struct {
	Name     string `json:"name"`      //description name of CoinType
	MinUnit  Unit   `json:"min_unit"`  //the min unit of CoinType
	MainUnit Unit   `json:"main_unit"` //the max unit of CoinType
	Desc     string `json:"desc"`      //the description of CoinType
}

//TODO

//ConvertToMainCoin return the main denom coin from args
func (ct CoinType) ConvertToMainCoin(coin sdk.Coin) (sdk.Coin, error) {
	if !ct.hasUnit(coin.Denom) {
		return coin, errors.New("coinType unit (%s) not defined" + coin.Denom)
	}

	if ct.isMainUnit(coin.Denom) {
		return coin, nil
	}

	// dest amount = src amount * (10^(dest scale) / 10^(src scale))
	dstScale := sdk.NewDecFromInt(ct.MainUnit.GetScaleFactor())
	srcScale := sdk.NewDecFromInt(ct.MinUnit.GetScaleFactor())
	amount := sdk.NewDecFromInt(coin.Amount)

	amt := amount.Mul(dstScale).Quo(srcScale)
	return sdk.NewCoin(ct.MainUnit.Denom, amt.RoundInt()), nil
}

//ConvertToMinCoin return the min denom coin from args
func (ct CoinType) ConvertToMinCoin(coin sdk.Coin) (newCoins sdk.Coin, err error) {
	if !ct.hasUnit(coin.Denom) {
		return coin, errors.New("coinType unit (%s) not defined" + coin.Denom)
	}

	if ct.isMinUnit(coin.Denom) {
		return coin, nil
	}

	// dest amount = src amount * (10^(dest scale) / 10^(src scale))
	srcScale := sdk.NewDecFromInt(ct.MainUnit.GetScaleFactor())
	dstScale := sdk.NewDecFromInt(ct.MinUnit.GetScaleFactor())
	amount := sdk.NewDecFromInt(coin.Amount)

	amt := amount.Mul(dstScale).Quo(srcScale)
	return sdk.NewCoin(ct.MainUnit.Denom, amt.RoundInt()), nil
}

func (ct CoinType) isMainUnit(name string) bool {
	return ct.MainUnit.Denom == strings.TrimSpace(name)
}

func (ct CoinType) isMinUnit(name string) bool {
	return ct.MinUnit.Denom == strings.TrimSpace(name)
}

func (ct CoinType) hasUnit(name string) bool {
	if ct.isMainUnit(name) || ct.isMinUnit(name) {
		return true
	}
	return false
}
