package human

import (
	"fmt"

	"github.com/irisnet/irishub/types"
)

type Stringer interface {
	fmt.Stringer
	HumanString(assetConvert AssetConvert) string
}

type AssetConvert interface {
	ToMainUnit(coins types.Coins) string
}
