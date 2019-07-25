package types

import (
	"fmt"
)

type Stringer interface {
	fmt.Stringer
	HumanString(converter CoinsConverter) string
}

type CoinsConverter interface {
	ToMainUnit(coins Coins) string
}
