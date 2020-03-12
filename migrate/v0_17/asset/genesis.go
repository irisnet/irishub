package asset

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	Params Params `json:"params"` // asset params
	Tokens Tokens `json:"tokens"` // issued tokens
}

type Params struct {
	AssetTaxRate      sdk.Dec  `json:"asset_tax_rate"`       // e.g., 40%
	IssueTokenBaseFee sdk.Coin `json:"issue_token_base_fee"` // e.g., 300000*10^18iris-atto
	MintTokenFeeRatio sdk.Dec  `json:"mint_token_fee_ratio"` // e.g., 10%
}

type Tokens []FungibleToken

type FungibleToken struct {
	BaseToken `json:"base_token"`
}

type BaseToken struct {
	Id              string         `json:"id"`
	Family          AssetFamily    `json:"family"`
	Source          AssetSource    `json:"source"`
	Gateway         string         `json:"gateway"`
	Symbol          string         `json:"symbol"`
	Name            string         `json:"name"`
	Decimal         uint8          `json:"decimal"`
	CanonicalSymbol string         `json:"canonical_symbol"`
	MinUnitAlias    string         `json:"min_unit_alias"`
	InitialSupply   sdk.Int        `json:"initial_supply"`
	MaxSupply       sdk.Int        `json:"max_supply"`
	Mintable        bool           `json:"mintable"`
	Owner           sdk.AccAddress `json:"owner"`
}

type AssetFamily byte

const (
	FUNGIBLE AssetFamily = 0x00
	//NON_FUNGIBLE AssetFamily = 0x01
)

var (
	AssetFamilyToStringMap = map[AssetFamily]string{
		FUNGIBLE: "fungible",
		//NON_FUNGIBLE: "non-fungible",
	}
	StringToAssetFamilyMap = map[string]AssetFamily{
		"fungible": FUNGIBLE,
		//"non-fungible": NON_FUNGIBLE,
	}
)

func AssetFamilyFromString(str string) (AssetFamily, error) {
	if family, ok := StringToAssetFamilyMap[strings.ToLower(str)]; ok {
		return family, nil
	}
	return AssetFamily(0xff), errors.Errorf("'%s' is not a valid asset family", str)
}

func (family AssetFamily) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(fmt.Sprintf("%s", family.String())))
	default:
		s.Write([]byte(fmt.Sprintf("%v", byte(family))))
	}
}

func (family AssetFamily) String() string {
	return AssetFamilyToStringMap[family]
}

// Marshals to JSON using string
func (family AssetFamily) MarshalJSON() ([]byte, error) {
	return json.Marshal(family.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (family *AssetFamily) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return nil
	}

	bz2, err := AssetFamilyFromString(s)
	if err != nil {
		return err
	}
	*family = bz2
	return nil
}

type AssetSource byte

const (
	NATIVE   AssetSource = 0x00
	EXTERNAL AssetSource = 0x01
	GATEWAY  AssetSource = 0x02
)

var (
	AssetSourceToStringMap = map[AssetSource]string{
		NATIVE:   "native",
		EXTERNAL: "external",
		GATEWAY:  "gateway",
	}
	StringToAssetSourceMap = map[string]AssetSource{
		"native":   NATIVE,
		"external": EXTERNAL,
		"gateway":  GATEWAY,
	}
)

func AssetSourceFromString(str string) (AssetSource, error) {
	if source, ok := StringToAssetSourceMap[strings.ToLower(str)]; ok {
		return source, nil
	}
	return AssetSource(0xff), errors.Errorf("'%s' is not a valid token source", str)
}

func (source AssetSource) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(fmt.Sprintf("%s", source.String())))
	default:
		s.Write([]byte(fmt.Sprintf("%v", byte(source))))
	}
}

func (source AssetSource) String() string {
	return AssetSourceToStringMap[source]
}

// Marshals to JSON using string
func (source AssetSource) MarshalJSON() ([]byte, error) {
	return json.Marshal(source.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (source *AssetSource) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return nil
	}

	bz2, err := AssetSourceFromString(s)
	if err != nil {
		return err
	}
	*source = bz2
	return nil
}
