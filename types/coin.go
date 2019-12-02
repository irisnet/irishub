package types

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	Iris     = "iris"
	IrisAtto = "iris-atto"
)

var (
	reABS           = `([a-z][0-9a-z]{2}[:])?`
	reCoinName      = reABS + `(([a-z][a-z0-9]{2,7}|x)\.)?([a-z][a-z0-9]{2,7})`
	reDenom         = reCoinName + `(-[a-z]{3,5})?`
	reAmount        = `[0-9]+(\.[0-9]+)?`
	reSpace         = `[[:space:]]*`
	reDenomCompiled = regexp.MustCompile(fmt.Sprintf(`^%s$`, reDenom))
	reCoinCompiled  = regexp.MustCompile(fmt.Sprintf(`^(%s)%s(%s)$`, reAmount, reSpace, reDenom))
)

func IsCoinMinDenomValid(denom string) bool {
	if denom != IrisAtto && (!strings.HasSuffix(denom, MinDenomSuffix) || strings.HasPrefix(denom, Iris+"-")) {
		return false
	}
	return reDenomCompiled.MatchString(denom)
}
