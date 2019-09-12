package types

import (
	"errors"
	"fmt"
	"strings"
)

const (
	//1 iris = 10^3 iris-milli
	Milli      = "milli"
	MilliScale = 3

	//1 iris = 10^6 iris-micro
	Micro      = "micro"
	MicroScale = 6

	//1 iris = 10^9 iris-nano
	Nano      = "nano"
	NanoScale = 9

	//1 iris = 10^12 iris-pico
	Pico      = "pico"
	PicoScale = 12

	//1 iris = 10^15 iris-femto
	Femto      = "femto"
	FemtoScale = 15

	//1 iris = 10^18 iris-atto
	Atto      = "atto"
	AttoScale = 18

	MinDenomSuffix = "-min"
)

var (
	IrisCoinType    = NewIrisCoinType()
	AttoScaleFactor = IrisCoinType.MinUnit.GetScaleFactor()
)

type Unit struct {
	Denom   string `json:"denom"`
	Decimal uint8  `json:"decimal"`
}

func (u Unit) String() string {
	return fmt.Sprintf("%s: %d",
		u.Denom, u.Decimal,
	)
}

func NewUnit(denom string, decimal uint8) Unit {
	return Unit{
		Denom:   denom,
		Decimal: decimal,
	}
}

func (u Unit) GetScaleFactor() Int {
	return NewIntWithDecimal(1, int(u.Decimal))
}

type Units []Unit

func (u Units) String() (out string) {
	for _, val := range u {
		out += val.String() + ",  "
	}
	if len(out) > 3 {
		out = out[:len(out)-3]
	}
	return
}

type CoinType struct {
	Name    string `json:"name"`
	MinUnit Unit   `json:"min_unit"`
	Units   Units  `json:"units"`
	Desc    string `json:"desc"`
}

func (ct CoinType) Convert(srcCoinStr string, destDenom string) (destCoinStr string, err error) {
	srcDenom, srcAmt, err := ParseCoinParts(srcCoinStr)
	if err != nil {
		return destCoinStr, err
	}
	var destUnit Unit
	if destUnit, err = ct.GetUnit(destDenom); err != nil {
		return destCoinStr, errors.New("destination unit (%s) not defined" + destDenom)
	}

	if srcUnit, err := ct.GetUnit(srcDenom); err == nil {
		if srcUnit.Denom == destDenom {
			return srcCoinStr, nil
		}
		// dest amount = src amount * (10^(dest scale) / 10^(src scale))
		rat := NewRatFromInt(destUnit.GetScaleFactor(), srcUnit.GetScaleFactor())
		amount, err := NewRatFromDecimal(srcAmt, int(ct.MinUnit.Decimal)) // convert src amount to dest unit
		if err != nil {
			return destCoinStr, err
		}
		amt := amount.Mul(rat).DecimalString(int(ct.MinUnit.Decimal))
		destCoinStr = fmt.Sprintf("%s%s", amt, destUnit.Denom)
		return destCoinStr, nil
	}
	return destCoinStr, errors.New("source unit (%s) not defined" + srcDenom)
}

func (ct CoinType) ConvertToMinDenomCoin(srcCoinStr string) (coin Coin, err error) {
	if destCoinStr, err := ct.Convert(srcCoinStr, ct.MinUnit.Denom); err == nil {
		coin, err = ParseCoin(destCoinStr)
		return coin, err
	}
	return coin, errors.New("convert error")
}

func (ct CoinType) GetUnit(denom string) (u Unit, err error) {
	for _, unit := range ct.Units {
		if strings.ToLower(denom) == strings.ToLower(unit.Denom) {
			return unit, nil
		}
	}
	return u, errors.New("unit (%s) not found" + denom)
}

func (ct CoinType) GetMainUnit() (unit Unit) {
	unit, _ = ct.GetUnit(ct.Name)
	return unit
}

func (ct CoinType) String() string {
	return fmt.Sprintf(`CoinType:
  Name:     %s
  MinUnit:  %s
  Units:    %s
  Desc:     %s`,
		ct.Name, ct.MinUnit, ct.Units, ct.Desc,
	)
}

func NewIrisCoinType() CoinType {
	units := make(Units, 7)

	units[0] = NewUnit(Iris, 0)
	units[1] = NewUnit(fmt.Sprintf("%s-%s", Iris, Milli), MilliScale)
	units[2] = NewUnit(fmt.Sprintf("%s-%s", Iris, Micro), MicroScale)
	units[3] = NewUnit(fmt.Sprintf("%s-%s", Iris, Nano), NanoScale)
	units[4] = NewUnit(fmt.Sprintf("%s-%s", Iris, Pico), PicoScale)
	units[5] = NewUnit(fmt.Sprintf("%s-%s", Iris, Femto), FemtoScale)
	units[6] = NewUnit(fmt.Sprintf("%s-%s", Iris, Atto), AttoScale)

	return CoinType{
		Name:    Iris,
		Units:   units,
		MinUnit: units[6],
		Desc:    "IRIS Network",
	}
}

func GetCoinName(coinStr string) (coinName string, err error) {
	denom, _, err := ParseCoinParts(coinStr)
	if err != nil {
		return coinName, err
	}

	if denom == Iris || denom == IrisAtto {
		return Iris, nil
	}

	if !strings.HasPrefix(denom, Iris+"-") && !strings.HasSuffix(denom, MinDenomSuffix) {
		return denom, nil
	}

	return GetCoinNameByDenom(denom)
}

func GetCoinNameByDenom(denom string) (coinName string, err error) {
	denom = strings.ToLower(denom)
	if strings.HasPrefix(denom, Iris+"-") {
		_, err := IrisCoinType.GetUnit(denom)
		if err != nil {
			return "", fmt.Errorf("invalid denom for getting coin name: %s", denom)
		}
		return Iris, nil
	}
	if !IsCoinMinDenomValid(denom) {
		return "", fmt.Errorf("invalid denom for getting coin name: %s", denom)
	}
	coinName = strings.TrimSuffix(denom, MinDenomSuffix)
	if coinName == "" {
		return coinName, fmt.Errorf("coin name is empty")
	}
	return coinName, nil
}

func GetCoinMinDenom(coinName string) (denom string, err error) {
	coinName = strings.ToLower(strings.TrimSpace(coinName))

	if coinName == Iris {
		return IrisAtto, nil
	}

	return fmt.Sprintf("%s%s", coinName, MinDenomSuffix), nil
}
