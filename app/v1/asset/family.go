package asset

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

type AssetFamily byte

const (
	FUNGIBLE     AssetFamily = 0x00
	NON_FUNGIBLE AssetFamily = 0x01
)

var (
	AssetFamilyToStringMap = map[AssetFamily]string{
		FUNGIBLE:     "fungible",
		NON_FUNGIBLE: "non-fungible",
	}
	StringToAssetFamilyMap = map[string]AssetFamily{
		"fungible":     FUNGIBLE,
		"non-fungible": NON_FUNGIBLE,
	}
)

func AssetFamilyFromString(str string) (AssetFamily, error) {
	if family, ok := StringToAssetFamilyMap[strings.ToLower(str)]; ok {
		return family, nil
	}
	return AssetFamily(0xff), errors.Errorf("'%s' is not a valid asset family", str)
}

func IsValidAssetFamily(family AssetFamily) bool {
	_, ok := AssetFamilyToStringMap[family]
	return ok
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

// Marshal needed for protobuf compatibility
func (family AssetFamily) Marshal() ([]byte, error) {
	return []byte{byte(family)}, nil
}

// Unmarshal needed for protobuf compatibility
func (family *AssetFamily) Unmarshal(data []byte) error {
	*family = AssetFamily(data[0])
	return nil
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
