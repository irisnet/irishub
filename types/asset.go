package types

import (
	"fmt"
	"regexp"
	"strings"
)

const FormatUniABSPrefix = "uni:"

var (
	reToken = regexp.MustCompile(`[A-Za-z0-9\.]{3,17}`)
	reDnm   = regexp.MustCompile(`[A-Za-z0-9\.\-]{3,21}`)
)

// ConvertIdToTokenKeyId return the store key suffix of a token
func ConvertIdToTokenKeyId(tokenId string) (key string, err error) {
	if !reToken.MatchString(tokenId) {
		return "", fmt.Errorf("token id convert error: invalid denom")
	}

	if strings.Contains(tokenId, ".") {
		return strings.ToLower(tokenId), nil
	} else {
		return strings.ToLower(fmt.Sprintf("i.%s", tokenId)), nil
	}
}

// ConvertDenomToTokenKeyId return the store key suffix of a token
func ConvertDenomToTokenKeyId(denom string) (key string, err error) {
	tokenId, err := ConvertDenomToTokenId(denom)
	if err != nil {
		return "", err
	}

	key, err = ConvertIdToTokenKeyId(tokenId)
	if err != nil {
		return "", err
	}

	return key, nil
}

// ConvertDenomToTokenId return the token id of the given denom
func ConvertDenomToTokenId(denom string) (tokenId string, err error) {
	if !reDnm.MatchString(denom) {
		return "", fmt.Errorf("token id convert error: invalid denom")
	}

	tokenId = strings.ToLower(strings.Split(denom, "-")[0])
	return tokenId, nil
}
