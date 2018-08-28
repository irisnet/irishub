package types

import (
	"fmt"
	"strings"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//   "that's one big rat!"
//          ______
//         / / /\ \____oo
//     __ /___...._____ _\o
//  __|     |_    |_

// NOTE: never use new(Rat) or else
// we will panic unmarshalling into the
// nil embedded big.Rat
type Rat struct {
	sdk.Rat
}

func NewRat(rat sdk.Rat) Rat{
	return Rat{
		rat,
	}
}

func (r Rat) DecimalString(prec int) string {
	floatStr := r.Rat.Rat.FloatString(prec)
	str := strings.Split(floatStr,".")
	if len(str) == 1 {
		return str[0]
	}
	dot := strings.TrimRightFunc(str[1],func(rune rune) bool{
		return rune == '0'
	})
	if len(dot) == 0 {
		return str[0]
	}
	return fmt.Sprintf("%s.%s",str[0],dot)
}