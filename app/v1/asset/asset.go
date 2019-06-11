package asset

import (
	"github.com/irisnet/irishub/types"
	"strings"
)

type Asset interface {
	GetDecimal() uint8
	IsMintable() bool
	GetUniqueID() string
	GetDenom() string
}

type BaseAsset struct {
	Family         AssetFamily      `json:"family"`
	Source         AssetSource      `json:"source"`
	Gateway        string           `json:"gateway"`
	Symbol         string           `json:"symbol"`
	Name           string           `json:"name"`
	Decimal        uint8            `json:"decimal"`
	SymbolMinAlias string           `json:"symbol_min_alias"`
	InitialSupply  uint64           `json:"initial_supply"`
	MaxSupply      uint64           `json:"max_supply"`
	Mintable       bool             `json:"mintable"`
	Owner          types.AccAddress `json:"owner"`
}

func NewBaseAsset(family AssetFamily, source AssetSource, gateway string, symbol string, name string, decimal uint8, alias string, initialSupply uint64, maxSupply uint64, mintable bool, owner types.AccAddress) BaseAsset {
	return BaseAsset{
		Family:         family,
		Source:         source,
		Gateway:        gateway,
		Symbol:         symbol,
		Name:           name,
		Decimal:        decimal,
		SymbolMinAlias: alias,
		InitialSupply:  initialSupply,
		MaxSupply:      maxSupply,
		Mintable:       mintable,
		Owner:          owner,
	}
}

func (BaseAsset) GetFamily() AssetFamily {
	panic("implement me")
}

func (BaseAsset) GetDecimal() uint8 {
	panic("implement me")
}

func (BaseAsset) IsMintable() bool {
	panic("implement me")
}

func (ba BaseAsset) GetUniqueID() string {
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

	return "invalid_asset_id"
}

func (ba BaseAsset) GetDenom() string {
	var sb strings.Builder
	sb.WriteString(ba.GetUniqueID())
	sb.WriteString("-min")
	return strings.ToLower(sb.String())
}

// Fungible Token
type FungibleToken struct {
	BaseAsset
}

func NewFungibleToken(source AssetSource, gateway string, symbol string, name string, decimal uint8, alias string, initialSupply uint64, maxSupply uint64, mintable bool, owner types.AccAddress) FungibleToken {
	return FungibleToken{
		BaseAsset: NewBaseAsset(
			FUNGIBLE, source, gateway, symbol, name, decimal, alias, initialSupply, maxSupply, mintable, owner,
		),
	}
}

func (ft FungibleToken) GetDecimal() uint8 {
	return ft.Decimal
}

func (ft FungibleToken) IsMintable() bool {
	return ft.Mintable
}

// Non-fungible Token
type NonFungibleToken struct {
	BaseAsset
}

func (nft NonFungibleToken) GetDecimal() uint8 {
	return 0
}

func (nft NonFungibleToken) IsMintable() bool {
	return true
}
