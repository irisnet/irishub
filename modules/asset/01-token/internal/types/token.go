package types

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	iristypes "github.com/irisnet/irishub/types"
)

var (
	MaximumTokenMaxSupply   = uint64(1000000000000) // Maximal limitation for token max supply，1000 billion
	MaximumTokenInitSupply  = uint64(100000000000)  // Maximal limitation for token initial supply，100 billion
	MaximumTokenDecimal     = uint8(18)             // Maximal limitation for token decimal
	MinimumTokenSymbolSize  = 3                     // Minimal limitation for the length of the token's symbol
	MaximumTokenSymbolSize  = 8                     // Maximal limitation for the length of the token's symbol
	MinimumTokenMinUnitSize = 3                     // Minimal limitation for the length of the token's min_unit
	MaximumTokenMinUnitSize = 10                    // Maximal limitation for the length of the token's min_unit
	MaximumTokenNameSize    = 32                    // Maximal limitation for the length of the token's name

	IsAlphaNumeric     = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString   // Only accepts alphanumeric characters
	IsAlphaNumericDash = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString // Only accepts alphanumeric characters, _ and -
	IsBeginWithAlpha   = regexp.MustCompile(`^[a-zA-Z].*`).MatchString      // Only accepts alpha characters
)

// FungibleToken
type FungibleToken struct {
	Symbol        string         `json:"symbol" yaml:"symbol"`                 // Globally unique asset identifier
	Name          string         `json:"name" yaml:"name"`                     // The name of the token, for example: Irisnet Network
	Scale         uint8          `json:"scale" yaml:"scale"`                   // Maximum number of decimals supported by this token
	MinUnit       string         `json:"min_unit" yaml:"min_unit"`             // The smallest unit name of the token
	InitialSupply sdk.Int        `json:"initial_supply" yaml:"initial_supply"` // Initial Token Issuance
	MaxSupply     sdk.Int        `json:"max_supply" yaml:"max_supply"`         // Maximum Token Issuance
	Mintable      bool           `json:"mintable" yaml:"mintable"`             // Is it possible to issue additional shares after the token is issued?
	Owner         sdk.AccAddress `json:"owner" yaml:"owner"`                   // The actual controller of the token
}

// NewFungibleToken - construct fungible token
func NewFungibleToken(symbol string, name string,
	scale uint8, minUnit string, initialSupply sdk.Int, maxSupply sdk.Int,
	mintable bool, owner sdk.AccAddress,
) FungibleToken {
	symbol = strings.ToLower(strings.TrimSpace(symbol))
	minUnit = strings.ToLower(strings.TrimSpace(minUnit))
	name = strings.TrimSpace(name)

	if maxSupply.IsZero() {
		if mintable {
			maxSupply = sdk.NewInt(int64(MaximumTokenMaxSupply))
		} else {
			maxSupply = initialSupply
		}
	}

	return FungibleToken{
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

// GetScale returns scale
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
	return iristypes.CoinType{
		Name:     ft.Name,
		MinUnit:  iristypes.NewUnit(ft.Symbol, 0),
		MainUnit: iristypes.NewUnit(ft.GetMinUnit(), ft.Scale),
		Desc:     ft.Name,
	}
}

// Sanitize - sanitize strings type
func (ft *FungibleToken) Sanitize() {
	ft.Symbol = strings.ToLower(strings.TrimSpace(ft.Symbol))
	ft.MinUnit = strings.ToLower(strings.TrimSpace(ft.MinUnit))
	ft.Name = strings.TrimSpace(ft.Name)
}

// Tokens - construct FungibleToken array
type Tokens []FungibleToken

// Validate - validate Tokens
func (tokens Tokens) Validate() error {
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

// ValidateName - check the validity of the name
func ValidateName(name string) error {
	nameLen := len(name)
	if nameLen == 0 || nameLen > MaximumTokenNameSize {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("invalid token name %s, only accepts length (0, %d]", name, MaximumTokenNameSize))
	}
	return nil
}

// ValidateSymbol - check the validity of the symbol
func ValidateSymbol(symbol string) error {
	symbolLen := len(symbol)
	if symbolLen < MinimumTokenSymbolSize || symbolLen > MaximumTokenSymbolSize ||
		!IsBeginWithAlpha(symbol) || !IsAlphaNumeric(symbol) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("invalid token symbol %s, only accepts alphanumeric characters, and begin with an english letter, length [%d, %d]", symbol, MinimumTokenSymbolSize, MaximumTokenSymbolSize))
	}

	if strings.Contains(strings.ToLower(symbol), iristypes.Iris) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("invalid token symbol %s, can not contain native token symbol %s", symbol, iristypes.Iris))
	}
	return nil
}

// ValidateScale - check the validity of the scale
func ValidateScale(scale uint8) error {
	if scale > MaximumTokenDecimal {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("invalid token decimal %d, only accepts value [0, %d]", scale, MaximumTokenDecimal))
	}
	return nil
}

// ValidateMinUnit - check the validity of the minUnit
func ValidateMinUnit(minUnit string) error {
	minUnitsLen := len(minUnit)
	if minUnitsLen < MinimumTokenMinUnitSize ||
		minUnitsLen > MaximumTokenMinUnitSize ||
		!IsAlphaNumericDash(minUnit) ||
		!IsBeginWithAlpha(minUnit) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("invalid token min_unit %s, only accepts alphanumeric characters, and begin with an english letter, length [%d, %d]", minUnit, MinimumTokenMinUnitSize, MaximumTokenMinUnitSize))
	}
	if strings.Contains(strings.ToLower(minUnit), iristypes.Iris) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("invalid token minUnit %s, can not contain native token minUnit %s", minUnit, iristypes.Iris))
	}
	return nil
}

// ValidateSupply - check the validity of the initialSupply,maxSupply
func ValidateSupply(initialSupply, maxSupply uint64) error {
	if initialSupply > MaximumTokenInitSupply {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("invalid token initial supply %d, only accepts value [0, %d]", initialSupply, MaximumTokenInitSupply))
	}

	if maxSupply < initialSupply || maxSupply > MaximumTokenMaxSupply {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("invalid token max supply %d, only accepts value [%d, %d]", maxSupply, initialSupply, MaximumTokenMaxSupply))
	}
	return nil
}

// ValidateMaxSupply - check the validity of the maxSupply
func ValidateMaxSupply(maxSupply uint64) error {
	if maxSupply > MaximumTokenMaxSupply {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("invalid token max supply %d, only accepts value (0, %d]", maxSupply, MaximumTokenMaxSupply))
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
