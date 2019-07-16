package human

import (
	"github.com/irisnet/irishub/types"
)

type Stringer interface {
	HumanString(assetConvert AssetConvert) string
}

type AssetConvert interface {
	ToMainUnit(coins types.Coins) string
}
