package asset

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

type TokenSource byte

const (
	NATIVE   TokenSource = 0x00
	EXTERNAL TokenSource = 0x01
	GATEWAY  TokenSource = 0x02
)

var (
	TokenSourceToStringMap = map[TokenSource]string{
		NATIVE:   "native",
		EXTERNAL: "external",
		GATEWAY:  "gateway",
	}
	StringToTokenSourceMap = map[string]TokenSource{
		"native":   NATIVE,
		"external": EXTERNAL,
		"gateway":  GATEWAY,
	}
)

func TokenSourceFromString(str string) (TokenSource, error) {
	if source, ok := StringToTokenSourceMap[strings.ToLower(str)]; ok {
		return source, nil
	}
	return TokenSource(0xff), errors.Errorf("'%s' is not a valid token source", str)
}

func IsValidTokenSource(source TokenSource) bool {
	_, ok := TokenSourceToStringMap[source]
	return ok
}

func (source TokenSource) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(fmt.Sprintf("%s", source.String())))
	default:
		s.Write([]byte(fmt.Sprintf("%v", byte(source))))
	}
}

func (source TokenSource) String() string {
	return TokenSourceToStringMap[source]
}

// Marshal needed for protobuf compatibility
func (source TokenSource) Marshal() ([]byte, error) {
	return []byte{byte(source)}, nil
}

// Unmarshal needed for protobuf compatibility
func (source *TokenSource) Unmarshal(data []byte) error {
	*source = TokenSource(data[0])
	return nil
}

// Marshals to JSON using string
func (source TokenSource) MarshalJSON() ([]byte, error) {
	return json.Marshal(source.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (source *TokenSource) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return nil
	}

	bz2, err := TokenSourceFromString(s)
	if err != nil {
		return err
	}
	*source = bz2
	return nil
}
