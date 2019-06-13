package asset

import (
	"fmt"
	"strings"

	"github.com/irisnet/irishub/types"
)

type Asset interface {
	GetDecimal() uint8
	IsMintable() bool
	GetUniqueID() string
	GetDenom() string
	String() string

	GetOwner() types.AccAddress
	GetSource() AssetSource
	GetSymbol() string
	GetGateway() string
	GetInitSupply() uint64
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

func (BaseAsset) GetDecimal() uint8 {
	panic("implement me")
}

func (BaseAsset) IsMintable() bool {
	panic("implement me")
}

func (ba BaseAsset) GetOwner() types.AccAddress {
	return ba.Owner
}

func (ba BaseAsset) GetSource() AssetSource {
	return ba.Source
}

func (ba BaseAsset) GetSymbol() string {
	return ba.Symbol
}

func (ba BaseAsset) GetGateway() string {
	return ba.Gateway
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

func (ba BaseAsset) GetInitSupply() uint64 {
	return ba.InitialSupply
}

// String implements fmt.Stringer
func (ba BaseAsset) String() string {
	return fmt.Sprintf(`Asset %s:
  Family:            %s
  Source:            %s
  Symbol:            %s
  Symbol Min Alias:  %s
  Decimal:           %d
  Initial Supply:    %d
  Max Supply:        %d
  Mintable:          %v
  Owner:             %s`,
		ba.GetUniqueID(), ba.Family, ba.Source, ba.Symbol, ba.SymbolMinAlias,
		ba.Decimal, ba.InitialSupply, ba.MaxSupply, ba.Mintable, ba.Owner.String())
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

func GetKeyID(source AssetSource, symbol string, gateway string) (string, types.Error) {
	switch source {
	case NATIVE:
		return strings.ToLower(fmt.Sprintf("i.%s", symbol)), nil
	case EXTERNAL:
		return strings.ToLower(fmt.Sprintf("x.%s", symbol)), nil
	case GATEWAY:
		return strings.ToLower(fmt.Sprintf("%s.%s", gateway, symbol)), nil
	default:
		return "", ErrInvalidAssetSource(DefaultCodespace, source)
	}
}

func GetKeyIDFromUniqueID(uniqueID string) string {
	if strings.Contains(uniqueID, ".") {
		return strings.ToLower(uniqueID)
	} else {
		return strings.ToLower(fmt.Sprintf("i.%s", uniqueID))
	}
}
