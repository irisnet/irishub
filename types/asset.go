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

// ConvertIDToTokenKeyID return the store key suffix of a token
func ConvertIDToTokenKeyID(tokenID string) (key string, err error) {
	if !reToken.MatchString(tokenID) {
		return "", fmt.Errorf("token id convert error: invalid denom")
	}

	if strings.Contains(tokenID, ".") {
		return strings.ToLower(tokenID), nil
	} else {
		return strings.ToLower(fmt.Sprintf("i.%s", tokenID)), nil
	}
}

// ConvertDenomToTokenKeyID return the store key suffix of a token
func ConvertDenomToTokenKeyID(denom string) (key string, err error) {
	tokenID, err := ConvertDenomToTokenID(denom)
	if err != nil {
		return "", err
	}

	key, err = ConvertIDToTokenKeyID(tokenID)
	if err != nil {
		return "", err
	}

	return key, nil
}

// ConvertDenomToTokenID return the token id of the given denom
func ConvertDenomToTokenID(denom string) (tokenID string, err error) {
	if !reDnm.MatchString(denom) {
		return "", fmt.Errorf("token id convert error: invalid denom")
	}

	tokenID = strings.ToLower(strings.Split(denom, "-")[0])
	return tokenID, nil
}
