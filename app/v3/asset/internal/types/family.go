package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type AssetFamily byte

const (
	FUNGIBLE AssetFamily = 0x00
)

var (
	AssetFamilyToStringMap = map[AssetFamily]string{
		FUNGIBLE: "fungible",
	}
	StringToAssetFamilyMap = map[string]AssetFamily{
		"fungible": FUNGIBLE,
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
		return err
	}

	bz2, err := AssetFamilyFromString(s)
	if err != nil {
		return err
	}
	*family = bz2
	return nil
}

type Bool string

const (
	False Bool = "false"
	True  Bool = "true"
	Nil   Bool = ""
)

func (b Bool) ToBool() bool {
	v := string(b)
	if len(v) == 0 {
		return false
	}
	result, _ := strconv.ParseBool(v)
	return result
}

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
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	*b = Bool(s)
	return nil
}
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
