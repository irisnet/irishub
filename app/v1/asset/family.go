package asset

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

type TokenFamily byte

const (
	FUNGIBLE     TokenFamily = 0x00
	NON_FUNGIBLE TokenFamily = 0x01
)

var (
	TokenFamilyToStringMap = map[TokenFamily]string{
		FUNGIBLE:     "fungible",
		NON_FUNGIBLE: "non-fungible",
	}
	StringToTokenFamilyMap = map[string]TokenFamily{
		"fungible":     FUNGIBLE,
		"non-fungible": NON_FUNGIBLE,
	}
)

func TokenFamilyFromString(str string) (TokenFamily, error) {
	if family, ok := StringToTokenFamilyMap[strings.ToLower(str)]; ok {
		return family, nil
	}
	return TokenFamily(0xff), errors.Errorf("'%s' is not a valid token family", str)
}

func IsValidTokenFamily(family TokenFamily) bool {
	_, ok := TokenFamilyToStringMap[family]
	return ok
}

func (family TokenFamily) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(fmt.Sprintf("%s", family.String())))
	default:
		s.Write([]byte(fmt.Sprintf("%v", byte(family))))
	}
}

func (family TokenFamily) String() string {
	return TokenFamilyToStringMap[family]
}

// Marshal needed for protobuf compatibility
func (family TokenFamily) Marshal() ([]byte, error) {
	return []byte{byte(family)}, nil
}

// Unmarshal needed for protobuf compatibility
func (family *TokenFamily) Unmarshal(data []byte) error {
	*family = TokenFamily(data[0])
	return nil
}

// Marshals to JSON using string
func (family TokenFamily) MarshalJSON() ([]byte, error) {
	return json.Marshal(family.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (family *TokenFamily) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return nil
	}

	bz2, err := TokenFamilyFromString(s)
	if err != nil {
		return err
	}
	*family = bz2
	return nil
}
