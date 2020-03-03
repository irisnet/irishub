package types

import (
	"fmt"
	"strings"

	"github.com/irisnet/irishub/types"
	sdk "github.com/irisnet/irishub/types"
)

// FungibleToken defines a struct for the fungible token
type FungibleToken struct {
	Symbol          string           `json:"symbol"`
	Name            string           `json:"name"`
	Decimal         uint8            `json:"decimal"`
	CanonicalSymbol string           `json:"canonical_symbol"`
	MinUnitAlias    string           `json:"min_unit_alias"`
	InitialSupply   types.Int        `json:"initial_supply"`
	MaxSupply       types.Int        `json:"max_supply"`
	Mintable        bool             `json:"mintable"`
	Owner           types.AccAddress `json:"owner"`
}

// NewFungibleToken constructs a new FungibleToken instance
func NewFungibleToken(
	symbol,
	name,
	minUnit string,
	decimal uint8,
	initialSupply,
	maxSupply types.Int,
	mintable bool,
	owner types.AccAddress,
) FungibleToken {
	symbol = strings.ToLower(strings.TrimSpace(symbol))
	name = strings.TrimSpace(name)
	minUnit = strings.ToLower(strings.TrimSpace(minUnit))

	return FungibleToken{
		Symbol:        symbol,
		Name:          name,
		MinUnitAlias:  minUnit,
		Decimal:       decimal,
		InitialSupply: initialSupply,
		MaxSupply:     maxSupply,
		Mintable:      mintable,
		Owner:         owner,
	}
}

func (ft FungibleToken) GetSymbol() string {
	return ft.Symbol
}

func (ft FungibleToken) GetDecimal() uint8 {
	return ft.Decimal
}

func (ft FungibleToken) IsMintable() bool {
	return ft.Mintable
}

func (ft FungibleToken) GetOwner() types.AccAddress {
	return ft.Owner
}

func (ft FungibleToken) GetDenom() string {
	denom, _ := sdk.GetCoinMinDenom(ft.GetSymbol())
	return denom
}

func (ft FungibleToken) GetInitSupply() types.Int {
	return ft.InitialSupply
}

func (ft FungibleToken) GetCoinType() types.CoinType {
	units := make(types.Units, 2)
	units[0] = types.NewUnit(ft.GetSymbol(), 0)
	units[1] = types.NewUnit(ft.GetDenom(), ft.Decimal)

	return types.CoinType{
		Name:    ft.GetSymbol(),
		MinUnit: units[1],
		Units:   units,
		Desc:    ft.Name,
	}
}

// String implements fmt.Stringer
func (ft FungibleToken) String() string {
	ct := ft.GetCoinType()
	initSupply, _ := ct.Convert(types.NewCoin(ft.GetDenom(), ft.InitialSupply).String(), ft.GetSymbol())
	maxSupply, _ := ct.Convert(types.NewCoin(ft.GetDenom(), ft.MaxSupply).String(), ft.GetSymbol())

	return fmt.Sprintf(`FungibleToken:
  Name:              %s
  Symbol:            %s
  Scale:             %d
  MinUnit:           %s
  Initial Supply:    %s
  Max Supply:        %s
  Mintable:          %v
  Owner:             %s`,
		ft.Name, ft.Symbol,
		ft.Decimal, ft.MinUnitAlias, initSupply, maxSupply, ft.Mintable, ft.Owner)
}

// Tokens is a set of tokens
type Tokens []FungibleToken

// String implements Stringer
func (tokens Tokens) String() string {
	if len(tokens) == 0 {
		return "[]"
	}

	out := ""
	for _, token := range tokens {
		out += fmt.Sprintf("%s \n", token.String())
	}

	return out[:len(out)-1]
}

func (tokens Tokens) Validate() sdk.Error {
	if len(tokens) == 0 {
		return nil
	}

	for _, token := range tokens {
		exp := sdk.NewIntWithDecimal(1, int(token.Decimal))
		initialSupply := uint64(token.InitialSupply.Div(exp).Int64())
		maxSupply := uint64(token.MaxSupply.Div(exp).Int64())

		msg := NewMsgIssueToken(token.Symbol, token.MinUnitAlias, token.Name, token.Decimal, initialSupply, maxSupply, token.Mintable, token.Owner)
		if err := ValidateMsgIssueToken(&msg); err != nil {
			return err
		}
	}

	return nil
}

// CheckSymbol checks if the given symbol is valid
func CheckSymbol(symbol string) sdk.Error {
	if strings.Contains(strings.ToLower(symbol), sdk.Iris) {
		return ErrInvalidAssetSymbol(DefaultCodespace, fmt.Sprintf("symbol can not contains : %s", sdk.Iris))
	}

	if len(symbol) < MinimumAssetSymbolLen || len(symbol) > MaximumAssetSymbolLen {
		return ErrInvalidAssetSymbol(DefaultCodespace, fmt.Sprintf("invalid asset symbol: %s", symbol))
	}

	if !IsBeginWithAlpha(symbol) || !IsAlphaNumeric(symbol) {
		return ErrInvalidAssetSymbol(DefaultCodespace, fmt.Sprintf("invalid asset symbol: %s", symbol))
	}

	return nil
}
