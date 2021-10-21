package types

import (
	"fmt"
	"regexp"
	"strings"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/token/types"
)

const (
	DoNotModify = "[do-not-modify]"
	MinDenomLen = 3
	MaxDenomLen = 128

	MaxTokenURILen = 256

	ReservedPeg  = "peg"
	ReservedIBC  = "ibc"
	ReservedHTLT = "htlt"
	ReservedTIBC = "tibc"
)

var (
	// IsAlphaNumeric only accepts [a-z0-9]
	IsAlphaNumeric = regexp.MustCompile(`^[a-z0-9]+$`).MatchString
	// IsBeginWithAlpha only begin with [a-z]
	IsBeginWithAlpha = regexp.MustCompile(`^[a-z].*`).MatchString

	keywords          = strings.Join([]string{ReservedPeg, ReservedIBC, ReservedHTLT, ReservedTIBC}, "|")
	regexpKeywordsFmt = fmt.Sprintf("^(%s).*", keywords)
	regexpKeyword     = regexp.MustCompile(regexpKeywordsFmt).MatchString
)

// ValidateDenomID verifies whether the  parameters are legal
func ValidateDenomID(denomID string) error {
	if len(denomID) < MinDenomLen || len(denomID) > MaxDenomLen {
		return sdkerrors.Wrapf(ErrInvalidDenom, "the length of denom(%s) only accepts value [%d, %d]", denomID, MinDenomLen, MaxDenomLen)
	}
	if !IsBeginWithAlpha(denomID) || !IsAlphaNumeric(denomID) {
		return sdkerrors.Wrapf(ErrInvalidDenom, "the denom(%s) only accepts alphanumeric characters, and begin with an english letter", denomID)
	}
	return ValidateKeywords(denomID)
}

// ValidateTokenID verify that the tokenID is legal
func ValidateTokenID(tokenID string) error {
	if len(tokenID) < MinDenomLen || len(tokenID) > MaxDenomLen {
		return sdkerrors.Wrapf(ErrInvalidTokenID, "the length of nft id(%s) only accepts value [%d, %d]", tokenID, MinDenomLen, MaxDenomLen)
	}
	if !IsBeginWithAlpha(tokenID) || !IsAlphaNumeric(tokenID) {
		return sdkerrors.Wrapf(ErrInvalidTokenID, "nft id(%s) only accepts alphanumeric characters, and begin with an english letter", tokenID)
	}
	return nil
}

// ValidateTokenURI verify that the tokenURI is legal
func ValidateTokenURI(tokenURI string) error {
	if len(tokenURI) > MaxTokenURILen {
		return sdkerrors.Wrapf(ErrInvalidTokenURI, "the length of nft uri(%s) only accepts value [0, %d]", tokenURI, MaxTokenURILen)
	}
	return nil
}

// Modified returns whether the field is modified
func Modified(target string) bool {
	return target != types.DoNotModify
}

// ValidateKeywords checks if the given denomId begins with `DenomKeywords`
func ValidateKeywords(denomId string) error {
	if regexpKeyword(denomId) {
		return sdkerrors.Wrapf(ErrInvalidDenom, "invalid denomId: %s, can not begin with keyword: (%s)", denomId, keywords)
	}
	return nil
}
