package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	iristypes "github.com/irisnet/irishub/types"
)

type BaseToken struct {
	ID              string         `json:"id" yaml:"id"`
	Family          AssetFamily    `json:"family" yaml:"family"`
	Source          AssetSource    `json:"source" yaml:"source"`
	Symbol          string         `json:"symbol" yaml:"symbol"`
	Name            string         `json:"name" yaml:"name"`
	Decimal         uint8          `json:"decimal" yaml:"decimal"`
	CanonicalSymbol string         `json:"canonical_symbol" yaml:"canonical_symbol"`
	MinUnitAlias    string         `json:"min_unit_alias" yaml:"min_unit_alias"`
	InitialSupply   sdk.Int        `json:"initial_supply" yaml:"initial_supply"`
	MaxSupply       sdk.Int        `json:"max_supply" yaml:"max_supply"`
	Mintable        bool           `json:"mintable" yaml:"mintable"`
	Owner           sdk.AccAddress `json:"owner" yaml:"owner"`
}

func NewBaseToken(family AssetFamily, source AssetSource, symbol string, name string,
	decimal uint8, canonicalSymbol string, minUnitAlias string, initialSupply sdk.Int, maxSupply sdk.Int,
	mintable bool, owner sdk.AccAddress,
) BaseToken {
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
	BaseToken `json:"base_token" yaml:"base_token"`
}

// NewFungibleToken - construct fungible token
func NewFungibleToken(source AssetSource, symbol string, name string, decimal uint8,
	canonicalSymbol string, minUnitAlias string, initialSupply sdk.Int, maxSupply sdk.Int,
	mintable bool, owner sdk.AccAddress,
) FungibleToken {
	token := FungibleToken{
		BaseToken: NewBaseToken(
			FUNGIBLE, source, symbol, name, decimal, canonicalSymbol, minUnitAlias, initialSupply, maxSupply, mintable, owner,
		),
	}

	token.ID = token.GetUniqueID()
	return token
}

// GetDecimal returns Decimal
func (ft FungibleToken) GetDecimal() uint8 {
	return ft.Decimal
}

// IsMintable returns Mintable
func (ft FungibleToken) IsMintable() bool {
	return ft.Mintable
}

// GetOwner returns Owner
func (ft FungibleToken) GetOwner() sdk.AccAddress {
	return ft.Owner
}

// GetSource returns Source
func (ft FungibleToken) GetSource() AssetSource {
	return ft.Source
}

// GetSymbol returns Symbol
func (ft FungibleToken) GetSymbol() string {
	return ft.Symbol
}

// GetUniqueID returns UniqueID
func (ft FungibleToken) GetUniqueID() string {
	switch ft.Source {
	case NATIVE:
		return strings.ToLower(ft.Symbol)
	case EXTERNAL:
		return strings.ToLower(fmt.Sprintf("x.%s", ft.Symbol))
	default:
		return ""
	}
}

// GetDenom returns denom
func (ft FungibleToken) GetDenom() string {
	denom, _ := iristypes.GetCoinMinDenom(ft.GetUniqueID())
	return denom
}

// GetInitSupply returns InitialSupply
func (ft FungibleToken) GetInitSupply() sdk.Int {
	return ft.InitialSupply
}

// GetCoinType returns CoinType
func (ft FungibleToken) GetCoinType() iristypes.CoinType {
	units := make(iristypes.Units, 2)
	units[0] = iristypes.NewUnit(ft.GetUniqueID(), 0)
	units[1] = iristypes.NewUnit(ft.GetDenom(), ft.Decimal)
	return iristypes.CoinType{
		Name:    ft.GetUniqueID(), // UniqueID == Coin Name
		MinUnit: units[1],
		Units:   units,
		Desc:    ft.Name,
	}
}

// Sanitize - sanitize strings type
func (ft FungibleToken) Sanitize() FungibleToken {
	ft.Symbol = strings.ToLower(strings.TrimSpace(ft.Symbol))
	ft.CanonicalSymbol = strings.ToLower(strings.TrimSpace(ft.CanonicalSymbol))
	ft.MinUnitAlias = strings.ToLower(strings.TrimSpace(ft.MinUnitAlias))
	ft.Name = strings.TrimSpace(ft.Name)
	return ft
}

// Tokens - construct FungibleToken array
type Tokens []FungibleToken

// Validate - validate Tokens
func (tokens Tokens) Validate() sdk.Error {
	if len(tokens) == 0 {
		return nil
	}

	for _, token := range tokens {
		exp := sdk.NewIntWithDecimal(1, int(token.Decimal))
		initialSupply := uint64(token.InitialSupply.Quo(exp).Int64())
		maxSupply := uint64(token.MaxSupply.Quo(exp).Int64())
		msg := NewMsgIssueToken(
			token.Family, token.GetSource(), token.Symbol, token.CanonicalSymbol, token.Name,
			token.Decimal, token.MinUnitAlias, initialSupply, maxSupply, token.Mintable, token.Owner,
		)
		if err := ValidateMsgIssueToken(&msg); err != nil {
			return err
		}
	}
	return nil
}

// GetTokenID returns tokenId by source and symbol
func GetTokenID(source AssetSource, symbol string) (string, sdk.Error) {
	switch source {
	case NATIVE:
		return strings.ToLower(fmt.Sprintf("i.%s", symbol)), nil
	case EXTERNAL:
		return strings.ToLower(fmt.Sprintf("x.%s", symbol)), nil
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
	if len(symbol) < MinimumAssetSymbolSize || len(symbol) > MaximumAssetSymbolSize ||
		!IsBeginWithAlpha(symbol) || !IsAlphaNumeric(symbol) || strings.Contains(symbol, iristypes.Iris) {
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
