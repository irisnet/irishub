package types

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	iristypes "github.com/irisnet/irishub/types"
)

var (
	MaximumAssetMaxSupply   = uint64(1000000000000) // maximal limitation for asset max supply，1000 billion
	MaximumAssetInitSupply  = uint64(100000000000)  // maximal limitation for asset initial supply，100 billion
	MaximumAssetDecimal     = uint8(18)             // maximal limitation for asset decimal
	MinimumAssetSymbolSize  = 3                     // minimal limitation for the length of the asset's symbol
	MaximumAssetSymbolSize  = 8                     // maximal limitation for the length of the asset's symbol
	MinimumAssetMinUnitSize = 3                     // minimal limitation for the length of the asset's min_unit
	MaximumAssetMinUnitSize = 10                    // maximal limitation for the length of the asset's min_unit
	MaximumAssetNameSize    = 32                    // maximal limitation for the length of the asset's name

	IsAlphaNumeric     = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString   // only accepts alphanumeric characters
	IsAlphaNumericDash = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString // only accepts alphanumeric characters, _ and -
	IsBeginWithAlpha   = regexp.MustCompile(`^[a-zA-Z].*`).MatchString      // only accepts alpha characters
)

// BaseToken
type BaseToken struct {
	Symbol        string         `json:"symbol" yaml:"symbol"`
	Name          string         `json:"name" yaml:"name"`
	Scale         uint8          `json:"scale" yaml:"scale"`
	MinUnit       string         `json:"min_unit" yaml:"min_unit"`
	InitialSupply sdk.Int        `json:"initial_supply" yaml:"initial_supply"`
	MaxSupply     sdk.Int        `json:"max_supply" yaml:"max_supply"`
	Mintable      bool           `json:"mintable" yaml:"mintable"`
	Owner         sdk.AccAddress `json:"owner" yaml:"owner"`
}

// NewBaseToken - construct fungible token
func NewBaseToken(symbol string, name string,
	scale uint8, minUnit string, initialSupply sdk.Int, maxSupply sdk.Int,
	mintable bool, owner sdk.AccAddress,
) BaseToken {
	symbol = strings.ToLower(strings.TrimSpace(symbol))
	minUnit = strings.ToLower(strings.TrimSpace(minUnit))
	name = strings.TrimSpace(name)

	if maxSupply.IsZero() {
		if mintable {
			maxSupply = sdk.NewInt(int64(MaximumAssetMaxSupply))
		} else {
			maxSupply = initialSupply
		}
	}

	return BaseToken{
		Symbol:        symbol,
		Name:          name,
		Scale:         scale,
		MinUnit:       minUnit,
		InitialSupply: initialSupply,
		MaxSupply:     maxSupply,
		Mintable:      mintable,
		Owner:         owner,
	}
}

// FungibleToken
type FungibleToken struct {
	BaseToken `json:"base_token" yaml:"base_token"`
}

// NewFungibleToken - construct fungible token
func NewFungibleToken(symbol string, name string, scale uint8, minUnit string, initialSupply sdk.Int, maxSupply sdk.Int,
	mintable bool, owner sdk.AccAddress,
) FungibleToken {
	token := FungibleToken{
		BaseToken: NewBaseToken(symbol, name, scale, minUnit, initialSupply, maxSupply, mintable, owner),
	}
	return token
}

// GetDecimal returns scale
func (ft FungibleToken) GetScale() uint8 {
	return ft.Scale
}

// IsMintable returns Mintable
func (ft FungibleToken) IsMintable() bool {
	return ft.Mintable
}

// GetOwner returns Owner
func (ft FungibleToken) GetOwner() sdk.AccAddress {
	return ft.Owner
}

// GetSymbol returns Symbol
func (ft FungibleToken) GetSymbol() string {
	return ft.Symbol
}

// GetMinUnit returns MinUnit
func (ft FungibleToken) GetMinUnit() string {
	return ft.MinUnit
}

// GetInitSupply returns InitialSupply
func (ft FungibleToken) GetInitSupply() sdk.Int {
	return ft.InitialSupply
}

// GetCoinType returns CoinType
func (ft FungibleToken) GetCoinType() iristypes.CoinType {
	units := make(iristypes.Units, 2)
	units[0] = iristypes.NewUnit(ft.Symbol, 0)
	units[1] = iristypes.NewUnit(ft.GetMinUnit(), ft.Scale)
	return iristypes.CoinType{
		Name:    ft.Name,
		MinUnit: units[1],
		Units:   units,
		Desc:    ft.Name,
	}
}

// Sanitize - sanitize strings type
func (ft FungibleToken) Sanitize() FungibleToken {
	ft.Symbol = strings.ToLower(strings.TrimSpace(ft.Symbol))
	ft.MinUnit = strings.ToLower(strings.TrimSpace(ft.MinUnit))
	ft.Name = strings.TrimSpace(ft.Name)
	return ft
}

// Tokens - construct FungibleToken array
type Tokens []FungibleToken

// Validate - validate Tokens
func (tokens Tokens) Validate() sdk.Error {
	for _, token := range tokens {
		exp := sdk.NewIntWithDecimal(1, int(token.Scale))
		initialSupply := uint64(token.InitialSupply.Quo(exp).Int64())
		maxSupply := uint64(token.MaxSupply.Quo(exp).Int64())
		msg := NewMsgIssueToken(
			token.Symbol, token.Name, token.Scale, token.MinUnit,
			initialSupply, maxSupply, token.Mintable, token.Owner,
		)
		if err := ValidateMsgIssueToken(&msg); err != nil {
			return err
		}
	}
	return nil
}

func ValidateName(name string) sdk.Error {
	nameLen := len(name)
	if nameLen == 0 || nameLen > MaximumAssetNameSize {
		return ErrInvalidAssetName(DefaultCodespace, fmt.Sprintf("invalid token name %s, only accepts length (0, %d]", name, MaximumAssetNameSize))
	}
	return nil
}

func ValidateSymbol(symbol string) sdk.Error {
	symbolLen := len(symbol)
	if symbolLen < MinimumAssetSymbolSize || symbolLen > MaximumAssetSymbolSize ||
		!IsBeginWithAlpha(symbol) || !IsAlphaNumeric(symbol) {
		return ErrInvalidAssetSymbol(DefaultCodespace, fmt.Sprintf("invalid token symbol %s, only accepts alphanumeric characters, and begin with an english letter, length [%d, %d]", symbol, MinimumAssetSymbolSize, MaximumAssetSymbolSize))
	}

	// TODO
	if strings.Contains(strings.ToLower(symbol), "iris") {
		return ErrInvalidAssetSymbol(DefaultCodespace, fmt.Sprintf("invalid token symbol %s, can not contain native token symbol %s", symbol, "iris"))
	}
	return nil
}

func ValidateScale(scale uint8) sdk.Error {
	if scale > MaximumAssetDecimal {
		return ErrInvalidAssetScale(DefaultCodespace, fmt.Sprintf("invalid token decimal %d, only accepts value [0, %d]", scale, MaximumAssetDecimal))
	}
	return nil
}

func ValidateMinUnit(minUnit string) sdk.Error {
	minUnitsLen := len(minUnit)
	if minUnitsLen > 0 && (minUnitsLen < MinimumAssetMinUnitSize ||
		minUnitsLen > MaximumAssetMinUnitSize ||
		!IsAlphaNumeric(minUnit) ||
		!IsBeginWithAlpha(minUnit)) {
		return ErrInvalidAssetMinUnit(DefaultCodespace, fmt.Sprintf("invalid token min_unit_alias %s, only accepts alphanumeric characters, and begin with an english letter, length [%d, %d]", minUnit, MinimumAssetMinUnitSize, MaximumAssetMinUnitSize))
	}
	return nil
}

func ValidateSupply(initialSupply, maxSupply uint64) sdk.Error {
	if initialSupply > MaximumAssetInitSupply {
		return ErrInvalidAssetInitSupply(DefaultCodespace, fmt.Sprintf("invalid token initial supply %d, only accepts value [0, %d]", initialSupply, MaximumAssetInitSupply))
	}

	if maxSupply < initialSupply || maxSupply > MaximumAssetMaxSupply {
		return ErrInvalidAssetMaxSupply(DefaultCodespace, fmt.Sprintf("invalid token max supply %d, only accepts value [%d, %d]", maxSupply, initialSupply, MaximumAssetMaxSupply))
	}
	return nil
}

func ValidateMaxSupply(maxSupply uint64) sdk.Error {
	if maxSupply > MaximumAssetMaxSupply {
		return ErrInvalidAssetMaxSupply(DefaultCodespace, fmt.Sprintf("invalid token max supply %d, only accepts value (0, %d]", maxSupply, MaximumAssetMaxSupply))

	}
	return nil
}

type Bool string

const (
	False Bool = "false"
	True  Bool = "true"
	Nil   Bool = ""
)

// ToBool
func (b Bool) ToBool() bool {
	v := string(b)
	if len(v) == 0 {
		return false
	}
	result, _ := strconv.ParseBool(v)
	return result
}

// ToBool
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

// Unmarshals from JSON assuming Bech32 encoding
func (b *Bool) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return nil
	}
	*b = Bool(s)
	return nil
}

// ParseBool
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
