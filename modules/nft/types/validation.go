package types

import (
	"fmt"
	"regexp"
	"strings"

	errormod "cosmossdk.io/errors"
)

const (
	DoNotModify = "[do-not-modify]"

	MaxTokenURILen = 256

	ReservedPeg  = "peg"
	ReservedIBC  = "ibc"
	ReservedHTLT = "htlt"
	ReservedTIBC = "tibc"
)

var (
	// DenomID or TokenID can be 3 ~ 101 characters long and support letters, followed by either
	// a letter, a number or a separator ('/', ':', '.', '_' or '-').
	idString = `[a-z][a-zA-Z0-9/]{2,100}`
	regexpID = regexp.MustCompile(fmt.Sprintf(`^%s$`, idString)).MatchString

	keywords          = strings.Join([]string{ReservedPeg, ReservedIBC, ReservedHTLT, ReservedTIBC}, "|")
	regexpKeywordsFmt = fmt.Sprintf("^(%s).*", keywords)
	regexpKeyword     = regexp.MustCompile(regexpKeywordsFmt).MatchString
)

// ValidateDenomID verifies whether the  parameters are legal
func ValidateDenomID(denomID string) error {
	boolPrifix := strings.HasPrefix(denomID, "tibc-")
	if !regexpID(denomID) && !boolPrifix {
		return errormod.Wrapf(ErrInvalidDenom, "denomID can only accept characters that match the regular expression: (%s),but got (%s)", idString, denomID)
	}
	return nil
}

// ValidateTokenID verify that the tokenID is legal
func ValidateTokenID(tokenID string) error {
	if !regexpID(tokenID) {
		return errormod.Wrapf(ErrInvalidDenom, "tokenID can only accept characters that match the regular expression: (%s),but got (%s)", idString, tokenID)
	}
	return nil
}

// ValidateTokenURI verify that the tokenURI is legal
func ValidateTokenURI(tokenURI string) error {
	if len(tokenURI) > MaxTokenURILen {
		return errormod.Wrapf(ErrInvalidTokenURI, "the length of nft uri(%s) only accepts value [0, %d]", tokenURI, MaxTokenURILen)
	}
	return nil
}

// Modified returns whether the field is modified
func Modified(target string) bool {
	return target != DoNotModify
}

func Modify(origin, target string) string {
	if target == DoNotModify {
		return origin
	}
	return target
}

// ValidateKeywords checks if the given denomId begins with `DenomKeywords`
func ValidateKeywords(denomId string) error {
	if regexpKeyword(denomId) {
		return errormod.Wrapf(ErrInvalidDenom, "invalid denomId: %s, can not begin with keyword: (%s)", denomId, keywords)
	}
	return nil
}

func IsIBCDenom(denomID string) bool {
	return strings.HasPrefix(denomID, "ibc/")
}
