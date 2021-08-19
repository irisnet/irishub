package types

import (
	fmt "fmt"
	"math"
	"regexp"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// MaximumMaxSupply is the maximum limitation for the token max supply
	MaximumMaxSupply = math.MaxUint64
	// MaximumInitSupply is maximum limitation for the token initial supply，100 billion
	MaximumInitSupply = uint64(100000000000)
	// MaximumScale is the maximum limitation for token decimals
	MaximumScale = uint32(9)
	// MinimumSymbolLen is the minimum limitation for the length of the token's symbol
	MinimumSymbolLen = 3
	// MaximumSymbolLen is the maximum limitation for the length of the token's symbol
	MaximumSymbolLen = 64
	// MaximumNameLen is the maximum limitation for the length of the token's name
	MaximumNameLen = 32
	// MinimumMinUnitLen is the minimum limitation for the length of the token's min unit
	MinimumMinUnitLen = 3
	// MaximumMinUnitLen is the maximum limitation for the length of the token's min unit
	MaximumMinUnitLen = 64
)

const (
	ReservedPeg  = "peg"
	ReservedIBC  = "ibc"
	ReservedTIBC = "tibc"
	ReservedLpt  = "lpt"
	ReservedHTLT = "htlt"
)

var (
	keywords = strings.Join([]string{
		ReservedPeg, ReservedIBC, ReservedLpt, ReservedHTLT, ReservedTIBC,
	}, "|")
	regexpKeywordsFmt = fmt.Sprintf("^(%s).*", keywords)
	regexpKeyword     = regexp.MustCompile(regexpKeywordsFmt).MatchString

	regexpSymbolFmt = fmt.Sprintf("^[a-z][a-z0-9]{%d,%d}$", MinimumSymbolLen-1, MaximumSymbolLen-1)
	regexpSymbol    = regexp.MustCompile(regexpSymbolFmt).MatchString

	regexpMinUintFmt = fmt.Sprintf("^[a-z][a-z0-9]{%d,%d}$", MinimumMinUnitLen-1, MaximumMinUnitLen-1)
	regexpMinUint    = regexp.MustCompile(regexpMinUintFmt).MatchString
)

// ValidateToken checks if the given token is valid
func ValidateToken(token Token) error {
	if len(token.Owner) > 0 {
		if _, err := sdk.AccAddressFromBech32(token.Owner); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
		}
	}
	if err := ValidateName(token.Name); err != nil {
		return err
	}
	if err := ValidateSymbol(token.Symbol); err != nil {
		return err
	}
	if err := ValidateMinUnit(token.MinUnit); err != nil {
		return err
	}
	if err := ValidateInitialSupply(token.InitialSupply); err != nil {
		return err
	}
	if token.MaxSupply < token.InitialSupply {
		return sdkerrors.Wrapf(ErrInvalidMaxSupply, "invalid token max supply %d, only accepts value [%d, %d]", token.MaxSupply, token.InitialSupply, uint64(MaximumMaxSupply))
	}
	return ValidateScale(token.Scale)
}

// ValidateInitialSupply verifies whether the initial supply is legal
func ValidateInitialSupply(initialSupply uint64) error {
	if initialSupply > MaximumInitSupply {
		return sdkerrors.Wrapf(ErrInvalidInitSupply, "invalid token initial supply %d, only accepts value [0, %d]", initialSupply, MaximumInitSupply)
	}
	return nil
}

// ValidateName verifies whether the given name is legal
func ValidateName(name string) error {
	if len(name) == 0 || len(name) > MaximumNameLen {
		return sdkerrors.Wrapf(ErrInvalidName, "invalid token name %s, only accepts length (0, %d]", name, MaximumNameLen)
	}
	return nil
}

// ValidateScale verifies whether the given scale is legal
func ValidateScale(scale uint32) error {
	if scale > MaximumScale {
		return sdkerrors.Wrapf(ErrInvalidScale, "invalid token scale %d, only accepts value [0, %d]", scale, MaximumScale)
	}
	return nil
}

// ValidateMinUnit checks if the given min unit is valid
func ValidateMinUnit(minUnit string) error {
	if !regexpMinUint(minUnit) {
		return sdkerrors.Wrapf(ErrInvalidMinUnit, "invalid minUnit: %s, only accepts english lowercase letters and numbers, length [%d, %d], and begin with an english letter, regexp: %s", minUnit, MinimumMinUnitLen, MaximumMinUnitLen, regexpMinUintFmt)
	}
	return ValidateKeywords(minUnit)
}

// ValidateSymbol checks if the given symbol is valid
func ValidateSymbol(symbol string) error {
	if !regexpSymbol(symbol) {
		return sdkerrors.Wrapf(ErrInvalidSymbol, "invalid symbol: %s, only accepts english lowercase letters and numbers, length [%d, %d], and begin with an english letter, regexp: %s", symbol, MinimumSymbolLen, MaximumSymbolLen, regexpSymbolFmt)
	}
	return ValidateKeywords(symbol)
}

// ValidateKeywords checks if the given denom begins with `TokenKeywords`
func ValidateKeywords(denom string) error {
	if regexpKeyword(denom) {
		return sdkerrors.Wrapf(ErrInvalidSymbol, "invalid token: %s, can not begin with keyword: (%s)", denom, keywords)
	}
	return nil
}

// ValidateAmount checks if the given denom begins with `TokenKeywords`
func ValidateAmount(amount uint64) error {
	if amount == 0 {
		return sdkerrors.Wrapf(ErrInvalidMaxSupply, "invalid token amount %d, only accepts value (0, %d]", amount, uint64(MaximumMaxSupply))
	}
	return nil
}
