package types

import (
	"fmt"
	"strings"

	"github.com/irisnet/irishub/types"
	sdk "github.com/irisnet/irishub/types"
)

type BaseToken struct {
	Id              string           `json:"id"`
	Family          AssetFamily      `json:"family"`
	Source          AssetSource      `json:"source"`
	Gateway         string           `json:"gateway"`
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

func NewBaseToken(family AssetFamily, source AssetSource, gateway string, symbol string, name string, decimal uint8, canonicalSymbol string, minUnitAlias string, initialSupply types.Int, maxSupply types.Int, mintable bool, owner types.AccAddress) BaseToken {
	gateway = strings.ToLower(strings.TrimSpace(gateway))
	symbol = strings.ToLower(strings.TrimSpace(symbol))
	canonicalSymbol = strings.ToLower(strings.TrimSpace(canonicalSymbol))
	minUnitAlias = strings.ToLower(strings.TrimSpace(minUnitAlias))
	name = strings.TrimSpace(name)

	if maxSupply.IsZero() {
		if mintable {
			maxSupply = sdk.NewInt(int64(MaximumAssetMaxSupply))
		} else {
			maxSupply = initialSupply
		}
	}

	return BaseToken{
		Family:          family,
		Source:          source,
		Gateway:         gateway,
		Symbol:          symbol,
		Name:            name,
		Decimal:         decimal,
		CanonicalSymbol: canonicalSymbol,
		MinUnitAlias:    minUnitAlias,
		InitialSupply:   initialSupply,
		MaxSupply:       maxSupply,
		Mintable:        mintable,
		Owner:           owner,
	}
}

// FungibleToken
type FungibleToken struct {
	BaseToken `json:"base_token"`
}

func NewFungibleToken(source AssetSource, gateway string, symbol string, name string, decimal uint8, canonicalSymbol string, minUnitAlias string, initialSupply types.Int, maxSupply types.Int, mintable bool, owner types.AccAddress) FungibleToken {
	token := FungibleToken{
		BaseToken: NewBaseToken(
			FUNGIBLE, source, gateway, symbol, name, decimal, canonicalSymbol, minUnitAlias, initialSupply, maxSupply, mintable, owner,
		),
	}

	token.Id = token.GetUniqueID()
	return token
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

func (ft FungibleToken) GetSource() AssetSource {
	return ft.Source
}

func (ft FungibleToken) GetSymbol() string {
	return ft.Symbol
}

func (ft FungibleToken) GetGateway() string {
	return ft.Gateway
}

func (ft FungibleToken) GetUniqueID() string {
	switch ft.Source {
	case NATIVE:
		return strings.ToLower(ft.Symbol)
	case EXTERNAL:
		return strings.ToLower(fmt.Sprintf("x.%s", ft.Symbol))
	case GATEWAY:
		return strings.ToLower(fmt.Sprintf("%s.%s", ft.Gateway, ft.Symbol))
	default:
		return ""
	}
}

func (ft FungibleToken) GetDenom() string {
	denom, _ := sdk.GetCoinMinDenom(ft.GetUniqueID())
	return denom
}

func (ft FungibleToken) GetInitSupply() types.Int {
	return ft.InitialSupply
}

func (ft FungibleToken) GetCoinType() types.CoinType {
	units := make(types.Units, 2)
	units[0] = types.NewUnit(ft.GetUniqueID(), 0)
	units[1] = types.NewUnit(ft.GetDenom(), ft.Decimal)
	return types.CoinType{
		Name:    ft.GetUniqueID(),	// UniqueID == Coin Name
		MinUnit: units[1],
		Units:   units,
		Desc:    ft.Name,
	}
}

// String implements fmt.Stringer
func (ft FungibleToken) String() string {

	ct := ft.GetCoinType()

	initSupply, _ := ct.Convert(types.NewCoin(ft.GetDenom(), ft.InitialSupply).String(), ft.GetUniqueID())
	maxSupply, _ := ct.Convert(types.NewCoin(ft.GetDenom(), ft.MaxSupply).String(), ft.GetUniqueID())
	owner := ""
	if !ft.Owner.Empty() {
		owner = ft.Owner.String()
	}

	return fmt.Sprintf(`FungibleToken %s:
  Family:            %s
  Source:            %s
  Gateway:           %s
  Name:              %s
  Symbol:            %s
  CanonicalSymbol:   %s
  MinUnitAlias:      %s
  Decimal:           %d
  Initial Supply:    %s
  Max Supply:        %s
  Mintable:          %v
  Owner:             %s`,
		ft.GetUniqueID(), ft.Family, ft.Source, ft.Gateway, ft.Name, ft.Symbol, ft.CanonicalSymbol, ft.MinUnitAlias,
		ft.Decimal, initSupply, maxSupply, ft.Mintable, owner)
}

func (ft FungibleToken) Sanitize() FungibleToken {
	ft.Gateway = strings.ToLower(strings.TrimSpace(ft.Gateway))
	ft.Symbol = strings.ToLower(strings.TrimSpace(ft.Symbol))
	ft.CanonicalSymbol = strings.ToLower(strings.TrimSpace(ft.CanonicalSymbol))
	ft.MinUnitAlias = strings.ToLower(strings.TrimSpace(ft.MinUnitAlias))
	ft.Name = strings.TrimSpace(ft.Name)
	return ft
}

type Tokens []FungibleToken

func (tokens Tokens) String() string {
	if len(tokens) == 0 {
		return "[]"
	}

	out := ""
	for _, token := range tokens {
		out += fmt.Sprintf("%v \n", token.String())
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
		msg := NewMsgIssueToken(token.Family, token.GetSource(), token.Gateway, token.Symbol, token.CanonicalSymbol, token.Name, token.Decimal, token.MinUnitAlias, initialSupply, maxSupply, token.Mintable, token.Owner)
		if err := msg.ValidateBasic(); err != nil {
			return err
		}
	}
	return nil
}

// -----------------------------

func GetTokenID(source AssetSource, symbol string, gateway string) (string, types.Error) {
	switch source {
	case NATIVE:
		return strings.ToLower(fmt.Sprintf("i.%s", symbol)), nil
	case EXTERNAL:
		return strings.ToLower(fmt.Sprintf("x.%s", symbol)), nil
	case GATEWAY:
		return strings.ToLower(fmt.Sprintf("%s.%s", gateway, symbol)), nil
	default:
		return "", ErrInvalidAssetSource(DefaultCodespace, fmt.Sprintf("invalid asset source type %s", source))
	}
}

// CheckTokenID checks if the given token id is valid
func CheckTokenID(id string) sdk.Error {
	prefix, symbol := GetTokenIDParts(id)

	// check gateway moniker
	if prefix != "" && prefix != "i" && prefix != "x" {
		if err := ValidateMoniker(prefix); err != nil {
			return err
		}
	}

	// check symbol
	if len(symbol) < MinimumAssetSymbolSize || len(symbol) > MaximumAssetSymbolSize || !IsBeginWithAlpha(symbol) || !IsAlphaNumeric(symbol) || strings.Contains(symbol, sdk.Iris) {
		return ErrInvalidAssetSymbol(DefaultCodespace, fmt.Sprintf("invalid asset symbol: %s", symbol))
	}

	return nil
}

// GetTokenIDParts returns the source prefix and symbol
func GetTokenIDParts(id string) (prefix string, symbol string) {
	parts := strings.Split(strings.ToLower(id), ".")

	if len(parts) > 1 {
		// external or gateway asset
		prefix = parts[0]
		symbol = strings.Join(parts[1:], ".")
	} else {
		symbol = parts[0]
	}

	return
}
