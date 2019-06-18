package asset

import (
	"fmt"
	"strings"

	"github.com/irisnet/irishub/types"
)

type BaseToken struct {
	Id             string           `json:"id"`
	Family         AssetFamily      `json:"family"`
	Source         AssetSource      `json:"source"`
	Gateway        string           `json:"gateway"`
	Symbol         string           `json:"symbol"`
	Name           string           `json:"name"`
	Decimal        uint8            `json:"decimal"`
	SymbolAtSource string           `json:"symbol_at_source"`
	SymbolMinAlias string           `json:"symbol_min_alias"`
	InitialSupply  types.Int        `json:"initial_supply"`
	MaxSupply      types.Int        `json:"max_supply"`
	Mintable       bool             `json:"mintable"`
	Owner          types.AccAddress `json:"owner"`
}

func NewBaseToken(family AssetFamily, source AssetSource, gateway string, symbol string, name string, decimal uint8, symbolAtSource string, symbolMinAlias string, initialSupply types.Int, totalSupply types.Int, maxSupply types.Int, mintable bool, owner types.AccAddress) BaseToken {
	return BaseToken{
		Family:         family,
		Source:         source,
		Gateway:        strings.ToLower(gateway),
		Symbol:         strings.ToLower(symbol),
		Name:           name,
		Decimal:        decimal,
		SymbolAtSource: strings.ToLower(symbolAtSource),
		SymbolMinAlias: strings.ToLower(symbolMinAlias),
		InitialSupply:  initialSupply,
		MaxSupply:      maxSupply,
		Mintable:       mintable,
		Owner:          owner,
	}
}

// FungibleToken
type FungibleToken struct {
	BaseToken `json:"base_token"`
}

func NewFungibleToken(source AssetSource, gateway string, symbol string, name string, decimal uint8, symbolAtSource string, symbolMinAlias string, initialSupply types.Int, totalSupply types.Int, maxSupply types.Int, mintable bool, owner types.AccAddress) FungibleToken {
	token := FungibleToken{
		BaseToken: NewBaseToken(
			FUNGIBLE, source, gateway, symbol, name, decimal, symbolAtSource, symbolMinAlias, initialSupply, totalSupply, maxSupply, mintable, owner,
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
	return strings.ToLower(fmt.Sprintf("%s-min", ft.GetUniqueID()))
}

func (ft FungibleToken) GetInitSupply() types.Int {
	return ft.InitialSupply
}

func (ft FungibleToken) GetCoinType() types.CoinType {

	units := make(types.Units, 2)
	units[0] = types.NewUnit(ft.GetUniqueID(), 0)
	units[1] = types.NewUnit(ft.GetDenom(), int(ft.Decimal))
	return types.CoinType{
		Name:    ft.GetUniqueID(),
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

	return fmt.Sprintf(`FungibleToken %s:
  Family:            %s
  Source:            %s
  Gateway:           %s
  Name:              %s
  Symbol:            %s
  Symbol At Source:  %s
  Symbol Min Alias:  %s
  Decimal:           %d
  Initial Supply:    %s
  Max Supply:        %s
  Mintable:          %v
  Owner:             %s`,
		ft.GetUniqueID(), ft.Family, ft.Source, ft.Gateway, ft.Name, ft.Symbol, ft.SymbolAtSource, ft.SymbolMinAlias,
		ft.Decimal, initSupply, maxSupply, ft.Mintable, ft.Owner.String())
}

// -----------------------------

func GetKeyID(source AssetSource, symbol string, gateway string) (string, types.Error) {
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

func GetKeyIDFromUniqueID(uniqueID string) string {
	if strings.Contains(uniqueID, ".") {
		return strings.ToLower(uniqueID)
	} else {
		return strings.ToLower(fmt.Sprintf("i.%s", uniqueID))
	}
}
