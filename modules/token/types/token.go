package types

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"

	"github.com/gogo/protobuf/proto"
	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type TokenI interface {
	GetSymbol() string
	GetName() string
	GetScale() uint32
	GetMinUnit() string
	GetInitialSupply() uint64
	GetMaxSupply() uint64
	GetMintable() bool
	GetOwner() sdk.AccAddress

	ToMainCoin(coin sdk.Coin) (sdk.DecCoin, error)
	ToMinCoin(coin sdk.DecCoin) (sdk.Coin, error)
}

var _ proto.Message = &Token{}

// NewToken constructs a new Token instance
func NewToken(
	symbol string,
	name string,
	minUnit string,
	scale uint32,
	initialSupply,
	maxSupply uint64,
	mintable bool,
	owner sdk.AccAddress,
) Token {
	symbol = strings.ToLower(strings.TrimSpace(symbol))
	minUnit = strings.ToLower(strings.TrimSpace(minUnit))
	name = strings.TrimSpace(name)

	if maxSupply == 0 {
		if mintable {
			maxSupply = MaximumMaxSupply
		} else {
			maxSupply = initialSupply
		}
	}

	return Token{
		Symbol:        symbol,
		Name:          name,
		MinUnit:       minUnit,
		Scale:         scale,
		InitialSupply: initialSupply,
		MaxSupply:     maxSupply,
		Mintable:      mintable,
		Owner:         owner.String(),
	}
}

// GetSymbol implements exported.TokenI
func (t Token) GetSymbol() string {
	return t.Symbol
}

// GetName implements exported.TokenI
func (t Token) GetName() string {
	return t.Name
}

// GetScale implements exported.TokenI
func (t Token) GetScale() uint32 {
	return t.Scale
}

// GetMinUnit implements exported.TokenI
func (t Token) GetMinUnit() string {
	return t.MinUnit
}

// GetInitialSupply implements exported.TokenI
func (t Token) GetInitialSupply() uint64 {
	return t.InitialSupply
}

// GetMaxSupply implements exported.TokenI
func (t Token) GetMaxSupply() uint64 {
	return t.MaxSupply
}

// GetMintable implements exported.TokenI
func (t Token) GetMintable() bool {
	return t.Mintable
}

// GetOwner implements exported.TokenI
func (t Token) GetOwner() sdk.AccAddress {
	owner, _ := sdk.AccAddressFromBech32(t.Owner)
	return owner
}

func (t Token) String() string {
	bz, _ := yaml.Marshal(t)
	return string(bz)
}

//ToMainCoin return the main denom coin from args
func (t Token) ToMainCoin(coin sdk.Coin) (sdk.DecCoin, error) {
	if t.Symbol != coin.Denom && t.MinUnit != coin.Denom {
		return sdk.NewDecCoinFromDec(coin.Denom, sdk.ZeroDec()), sdkerrors.Wrapf(ErrTokenNotExists, "token not match")
	}

	if t.Symbol == coin.Denom {
		return sdk.NewDecCoin(coin.Denom, coin.Amount), nil
	}

	precision := math.Pow10(int(t.Scale))
	precisionStr := strconv.FormatFloat(precision, 'f', 0, 64)
	precisionDec, err := sdk.NewDecFromStr(precisionStr)
	if err != nil {
		return sdk.DecCoin{}, err
	}

	// dest amount = src amount / 10^(scale)
	amount := sdk.NewDecFromInt(coin.Amount).Quo(precisionDec)
	return sdk.NewDecCoinFromDec(t.Symbol, amount), nil
}

//ToMinCoin return the min denom coin from args
func (t Token) ToMinCoin(coin sdk.DecCoin) (newCoin sdk.Coin, err error) {
	if t.Symbol != coin.Denom && t.MinUnit != coin.Denom {
		return sdk.NewCoin(coin.Denom, sdk.ZeroInt()), sdkerrors.Wrapf(ErrTokenNotExists, "token not match")
	}

	if t.MinUnit == coin.Denom {
		return sdk.NewCoin(coin.Denom, coin.Amount.TruncateInt()), nil
	}

	precision := math.Pow10(int(t.Scale))
	precisionStr := strconv.FormatFloat(precision, 'f', 0, 64)
	precisionDec, err := sdk.NewDecFromStr(precisionStr)
	if err != nil {
		return sdk.Coin{}, err
	}

	// dest amount = src amount * 10^(dest scale)
	amount := coin.Amount.Mul(precisionDec)
	return sdk.NewCoin(t.MinUnit, amount.TruncateInt()), nil
}

// ValidateToken checks if the given token is valid
func ValidateToken(token Token) error {
	_, err := sdk.AccAddressFromBech32(token.Owner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	nameLen := len(strings.TrimSpace(token.Name))
	if nameLen == 0 || nameLen > MaximumNameLen {
		return sdkerrors.Wrapf(ErrInvalidName, "invalid token name %s, only accepts length (0, %d]", token.Name, MaximumNameLen)
	}

	if err := CheckSymbol(token.Symbol); err != nil {
		return err
	}

	if err := CheckMinUnit(token.MinUnit); err != nil {
		return err
	}

	if token.InitialSupply > MaximumInitSupply {
		return sdkerrors.Wrapf(ErrInvalidInitSupply, "invalid token initial supply %d, only accepts value [0, %d]", token.InitialSupply, MaximumInitSupply)
	}

	if token.MaxSupply < token.InitialSupply || token.MaxSupply > MaximumMaxSupply {
		return sdkerrors.Wrapf(ErrInvalidMaxSupply, "invalid token max supply %d, only accepts value [%d, %d]", token.MaxSupply, token.InitialSupply, MaximumMaxSupply)
	}

	if token.Scale > MaximumScale {
		return sdkerrors.Wrapf(ErrInvalidScale, "invalid token scale %d, only accepts value [0, %d]", token.Scale, MaximumScale)
	}

	return nil
}

// CheckMinUnit checks if the given minUnit is valid
func CheckMinUnit(minUnit string) error {
	minUnitLen := len(strings.TrimSpace(minUnit))
	if minUnitLen < MinimumMinUnitLen || minUnitLen > MaximumMinUnitLen || !IsAlphaNumeric(minUnit) || !IsBeginWithAlpha(minUnit) {
		return sdkerrors.Wrapf(ErrInvalidMinUnit, "invalid token min_unit %s, only accepts alphanumeric characters, and begin with an english letter, length [%d, %d]", minUnit, MinimumMinUnitLen, MaximumMinUnitLen)
	}
	return CheckKeywords(minUnit)
}

// CheckSymbol checks if the given symbol is valid
func CheckSymbol(symbol string) error {
	if len(symbol) < MinimumSymbolLen || len(symbol) > MaximumSymbolLen {
		return sdkerrors.Wrapf(ErrInvalidSymbol, "invalid symbol: %s,  only accepts length [%d, %d]", symbol, MinimumSymbolLen, MaximumSymbolLen)
	}

	if !IsBeginWithAlpha(symbol) || !IsAlphaNumeric(symbol) {
		return sdkerrors.Wrapf(ErrInvalidSymbol, "invalid symbol: %s, only accepts alphanumeric characters, and begin with an english letter", symbol)
	}
	return CheckKeywords(symbol)
}

// CheckKeywords checks if the given denom begin with `TokenKeywords`
func CheckKeywords(denom string) error {
	if IsBeginWithKeyword(denom) {
		return sdkerrors.Wrapf(ErrInvalidSymbol, "invalid token: %s, can not begin with keyword: (%s)", denom, keywords)
	}
	return nil
}

type Bool string

const (
	False Bool = "false"
	True  Bool = "true"
	Nil   Bool = ""
)

func (b Bool) ToBool() bool {
	v := string(b)
	if len(v) == 0 {
		return false
	}
	result, _ := strconv.ParseBool(v)
	return result
}

func (b Bool) String() string {
	return string(b)
}

// Marshal needed for protobuf compatibility
func (b Bool) Marshal() ([]byte, error) {
	return []byte(b), nil
}

// Unmarshal needed for protobuf compatibility
func (b *Bool) Unmarshal(data []byte) error {
	*b = Bool(data[:])
	return nil
}

// Marshals to JSON using string
func (b Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

// UnmarshalJSON from using string
func (b *Bool) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	*b = Bool(s)
	return nil
}
func ParseBool(v string) (Bool, error) {
	if len(v) == 0 {
		return Nil, nil
	}
	result, err := strconv.ParseBool(v)
	if err != nil {
		return Nil, err
	}
	if result {
		return True, nil
	}
	return False, nil
}
