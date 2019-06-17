package asset

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

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

func IsValidAssetSource(source AssetSource) bool {
	_, ok := AssetSourceToStringMap[source]
	return ok
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

// Marshal needed for protobuf compatibility
func (source AssetSource) Marshal() ([]byte, error) {
	return []byte{byte(source)}, nil
}

// Unmarshal needed for protobuf compatibility
func (source *AssetSource) Unmarshal(data []byte) error {
	*source = AssetSource(data[0])
	return nil
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
