package types

import (
	time "time"

	"github.com/irisnet/irismod/types/exported"
)

// Parameter store keys
var (
	KeyAssetParams = []byte("AssetParams") // asset params key

	DefaultPreviousBlockTime = time.Now()
)

// ParamKeyTable returns the TypeTable for coinswap module
func ParamKeyTable() exported.KeyTable {
	return exported.NewKeyTable().RegisterParamSet(&Params{})
}

func (p *Params) ParamSetPairs() exported.ParamSetPairs {
	return exported.ParamSetPairs{
		exported.NewParamSetPair(KeyAssetParams, &p.AssetParams, validateAssetParams),
	}
}
