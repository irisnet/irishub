package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	//1 iris = 10^3 iris-milli
	Milli = "milli"

	//1 iris = 10^6 iris-micro
	Micro = "micro"

	//1 iris = 10^9 iris-nano
	Nano = "nano"

	//1 iris = 10^12 iris-pico
	Pico = "pico"

	//1 iris = 10^15 iris-femto
	Femto = "femto"

	//1 iris = 10^18 iris-atto
	Atto = "atto"

	Native     Origin = 0x01
	External   Origin = 0x02
	UserIssued Origin = 0x03
)

var (
	MainUnit = func(coinName string) Unit {
		return NewUnit(coinName, 0)
	}

	MilliUnit = func(coinName string) Unit {
		denom := fmt.Sprintf("%s-%s", coinName, Milli)
		return NewUnit(denom, 3)
	}

	MicroUnit = func(coinName string) Unit {
		denom := fmt.Sprintf("%s-%s", coinName, Micro)
		return NewUnit(denom, 6)
	}

	NanoUnit = func(coinName string) Unit {
		denom := fmt.Sprintf("%s-%s", coinName, Nano)
		return NewUnit(denom, 9)
	}

	PicoUnit = func(coinName string) Unit {
		denom := fmt.Sprintf("%s-%s", coinName, Pico)
		return NewUnit(denom, 12)
	}

	FemtoUnit = func(coinName string) Unit {
		denom := fmt.Sprintf("%s-%s", coinName, Femto)
		return NewUnit(denom, 15)
	}

	AttoUnit = func(coinName string) Unit {
		denom := fmt.Sprintf("%s-%s", coinName, Atto)
		return NewUnit(denom, 18)
	}
)

type Origin = byte

func ToOrigin(origin string) (og Origin, err error) {
	switch strings.ToLower(origin) {
	case "native":
		return Native, nil
	case "external":
		return External, nil
	case "userissued":
		return UserIssued, nil
	}
	return og, errors.New("not support type:" + origin)
}

type Unit struct {
	Denom   string `json:"denom"`
	Decimal int    `json:"decimal"`
}

func NewUnit(denom string, decimal int) Unit {
	return Unit{
		Denom:   denom,
		Decimal: decimal,
	}
}
func (u Unit) GetPrecision() Int {
	return NewIntWithDecimal(1, u.Decimal)
}

type Units = []Unit

type CoinType struct {
	Name    string `json:"name"`
	MinUnit Unit   `json:"min_unit"`
	Units   Units  `json:"units"`
	Origin  Origin `json:"origin"`
	Desc    string `json:"desc"`
}

func (ct CoinType) Convert(orgCoinStr string, denom string) (destCoinStr string, err error) {
	orgDenom, orgAmt, err := GetCoin(orgCoinStr)
	if err != nil {
		return destCoinStr, err
	}
	var destUint Unit
	if destUint, err = ct.GetUnit(denom); err != nil {
		return destCoinStr, errors.New("not exist unit " + orgDenom)
	}
	// target Coin = original amount * (10^(target decimal) / 10^(original decimal))
	if orgUnit, err := ct.GetUnit(orgDenom); err == nil {
		rat := NewRatFromInt(destUint.GetPrecision(), orgUnit.GetPrecision())
		amount, err := NewRatFromDecimal(orgAmt, ct.MinUnit.Decimal) //Convert the original amount to the target accuracy
		if err != nil {
			return destCoinStr, err
		}
		amt := amount.Mul(rat).DecimalString(ct.MinUnit.Decimal)
		destCoinStr = fmt.Sprintf("%s%s", amt, destUint.Denom)
		return destCoinStr, nil
	}
	return destCoinStr, errors.New("not exist unit " + orgDenom)
}

func (ct CoinType) ConvertToMinCoin(coinStr string) (coin Coin, err error) {
	minUint := ct.GetMinUnit()

	if destCoinStr, err := ct.Convert(coinStr, minUint.Denom); err == nil {
		coin, err = parseCoin(destCoinStr)
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
	return u, errors.New("not find unit " + denom)
}

func (ct CoinType) GetMinUnit() (unit Unit) {
	return ct.MinUnit
}

func (ct CoinType) GetMainUnit() (unit Unit) {
	unit, _ = ct.GetUnit(ct.Name)
	return unit
}

func (ct CoinType) String() string {
	bz, _ := json.Marshal(ct)
	return string(bz)
}

func NewDefaultCoinType(name string) CoinType {
	units := GetDefaultUnits(name)
	return CoinType{
		Name:    name,
		Units:   units,
		MinUnit: units[6],
		Origin:  Native,
		Desc:    "IRIS Network",
	}
}

//TODO currently only iris token is supported
func CoinTypeKey(coinName string) string {
	return fmt.Sprintf("%s/%s/%s", "global", "coin_types", coinName)
}

func GetDefaultUnits(coin string) Units {
	units := make(Units, 7)
	units[0] = MainUnit(coin)
	units[1] = MilliUnit(coin)
	units[2] = MicroUnit(coin)
	units[3] = NanoUnit(coin)
	units[4] = PicoUnit(coin)
	units[5] = FemtoUnit(coin)
	units[6] = AttoUnit(coin)
	return units
}

func GetCoin(coinStr string) (denom, amount string, err error) {
	var (
		reDnm  = `[A-Za-z\-]{2,15}`
		reAmt  = `[0-9]+[.]?[0-9]*`
		reSpc  = `[[:space:]]*`
		reCoin = regexp.MustCompile(fmt.Sprintf(`^(%s)%s(%s)$`, reAmt, reSpc, reDnm))
	)

	coinStr = strings.TrimSpace(coinStr)

	matches := reCoin.FindStringSubmatch(coinStr)
	if matches == nil {
		err = fmt.Errorf("invalid coin expression: %s", coinStr)
		return
	}
	denom, amount = matches[2], matches[1]
	return
}

func parseCoin(coinStr string) (coin Coin, err error) {
	denom, amount, err := GetCoin(coinStr)
	if err != nil {
		return Coin{}, err
	}

	amt, ok := NewIntFromString(amount)
	if !ok {
		return Coin{}, fmt.Errorf("invalid coin amount: %s", amount)
	}
	denom = strings.ToLower(denom)
	return NewCoin(denom, amt), nil
}

func GetCoinName(coinStr string) (coinName string, err error) {
	denom, _, err := GetCoin(coinStr)
	if err != nil {
		return coinName, err
	}
	coinName = strings.Split(denom, "-")[0]
	coinName = strings.ToLower(coinName)
	if coinName == "" {
		return coinName, fmt.Errorf("coin name is empty")
	}
	return coinName, nil
}


