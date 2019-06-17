package asset

import (
	"fmt"
	"strings"

	"github.com/irisnet/irishub/types"
)

type Token interface {
	GetDecimal() uint8
	IsMintable() bool
	GetUniqueID() string
	GetDenom() string
	String() string

	GetOwner() types.AccAddress
	GetSource() TokenSource
	GetSymbol() string
	GetGateway() string
	GetInitSupply() types.Int
	GetTotalSupply() types.Int
	GetCoinType() types.CoinType
}

type BaseToken struct {
	Id             string           `json:"id"`
	Family         TokenFamily      `json:"family"`
	Source         TokenSource      `json:"source"`
	Gateway        string           `json:"gateway"`
	Symbol         string           `json:"symbol"`
	Name           string           `json:"name"`
	Decimal        uint8            `json:"decimal"`
	SymbolAtSource string           `json:"symbol_at_source"`
	SymbolMinAlias string           `json:"symbol_min_alias"`
	InitialSupply  types.Int        `json:"initial_supply"`
	TotalSupply    types.Int        `json:"total_supply"`
	MaxSupply      types.Int        `json:"max_supply"`
	Mintable       bool             `json:"mintable"`
	Owner          types.AccAddress `json:"owner"`
}

func NewBaseToken(family TokenFamily, source TokenSource, gateway string, symbol string, name string, decimal uint8, symbolAtSource string, symbolMinAlias string, initialSupply types.Int, totalSupply types.Int, maxSupply types.Int, mintable bool, owner types.AccAddress) BaseToken {
	baseToken := BaseToken{
		Family:         family,
		Source:         source,
		Gateway:        strings.ToLower(gateway),
		Symbol:         strings.ToLower(symbol),
		Name:           name,
		Decimal:        decimal,
		SymbolAtSource: strings.ToLower(symbolAtSource),
		SymbolMinAlias: strings.ToLower(symbolMinAlias),
		InitialSupply:  initialSupply,
		TotalSupply:    totalSupply,
		MaxSupply:      maxSupply,
		Mintable:       mintable,
		Owner:          owner,
	}

	baseToken.Id = baseToken.GetUniqueID()

	return baseToken
}

func (ba BaseToken) GetDecimal() uint8 {
	panic("implement me")
}

func (BaseToken) IsMintable() bool {
	panic("implement me")
}

func (BaseToken) GetCoinType() types.CoinType {
	panic("implement me")
}

// String implements fmt.Stringer
func (ba BaseToken) String() string {
	panic("implement me")
}

func (ba BaseToken) GetOwner() types.AccAddress {
	return ba.Owner
}

func (ba BaseToken) GetSource() TokenSource {
	return ba.Source
}

func (ba BaseToken) GetSymbol() string {
	return ba.Symbol
}

func (ba BaseToken) GetGateway() string {
	return ba.Gateway
}

func (ba BaseToken) GetUniqueID() string {
	if ba.Source == NATIVE {
		return strings.ToLower(ba.Symbol)
	}

	var sb strings.Builder

	if ba.Source == EXTERNAL {
		sb.WriteString("x.")
		sb.WriteString(ba.Symbol)
		return strings.ToLower(sb.String())
	}

	if ba.Source == GATEWAY {
		sb.WriteString(ba.Gateway)
		sb.WriteString(".")
		sb.WriteString(ba.Symbol)
		return strings.ToLower(sb.String())
	}

	return "invalid_token_id"
}

func (ba BaseToken) GetDenom() string {
	var sb strings.Builder
	sb.WriteString(ba.GetUniqueID())
	sb.WriteString("-min")
	return strings.ToLower(sb.String())
}

func (ba BaseToken) GetInitSupply() types.Int {
	return ba.InitialSupply
}

func (ba BaseToken) GetTotalSupply() types.Int {
	return ba.TotalSupply
}

// Fungible Token
type FungibleToken struct {
	BaseToken `json:"base_token"`
}

func NewFungibleToken(source TokenSource, gateway string, symbol string, name string, decimal uint8, symbolAtSource string, symbolMinAlias string, initialSupply types.Int, totalSupply types.Int, maxSupply types.Int, mintable bool, owner types.AccAddress) FungibleToken {
	return FungibleToken{
		BaseToken: NewBaseToken(
			FUNGIBLE, source, gateway, symbol, name, decimal, symbolAtSource, symbolMinAlias, initialSupply, totalSupply, maxSupply, mintable, owner,
		),
	}
}

func (ft FungibleToken) GetDecimal() uint8 {
	return ft.Decimal
}

func (ft FungibleToken) IsMintable() bool {
	return ft.Mintable
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
	totalSupply, _ := ct.Convert(types.NewCoin(ft.GetDenom(), ft.TotalSupply).String(), ft.GetUniqueID())

	return fmt.Sprintf(`Token %s:
  Family:            %s
  Source:            %s
  Gateway:           %s
  Name:              %s
  Symbol:            %s
  Symbol At Source:  %s
  Symbol Min Alias:  %s
  Decimal:           %d
  Initial Supply:    %s
  Total Supply:      %s
  Max Supply:        %s
  Mintable:          %v
  Owner:             %s`,
		ft.GetUniqueID(), ft.Family, ft.Source, ft.Gateway, ft.Name, ft.Symbol, ft.SymbolAtSource, ft.SymbolMinAlias,
		ft.Decimal, initSupply, totalSupply, maxSupply, ft.Mintable, ft.Owner.String())
}

// Non-fungible Token
type NonFungibleToken struct {
	BaseToken `json:"base_token"`
}

func (nft NonFungibleToken) GetDecimal() uint8 {
	return 0
}

func (nft NonFungibleToken) IsMintable() bool {
	return true
}

func GetKeyID(source TokenSource, symbol string, gateway string) (string, types.Error) {
	switch source {
	case NATIVE:
		return strings.ToLower(fmt.Sprintf("i.%s", symbol)), nil
	case EXTERNAL:
		return strings.ToLower(fmt.Sprintf("x.%s", symbol)), nil
	case GATEWAY:
		return strings.ToLower(fmt.Sprintf("%s.%s", gateway, symbol)), nil
	default:
		return "", ErrInvalidAssetSource(DefaultCodespace, fmt.Sprintf("invalid token source type %s", source))
	}
}

func GetKeyIDFromUniqueID(uniqueID string) string {
	if strings.Contains(uniqueID, ".") {
		return strings.ToLower(uniqueID)
	} else {
		return strings.ToLower(fmt.Sprintf("i.%s", uniqueID))
	}
}
