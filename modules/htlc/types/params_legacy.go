package types

import (
	time "time"
)

// Parameter store keys
var (
	KeyAssetParams = []byte("AssetParams") // asset params key

	DefaultPreviousBlockTime = time.Now()
)

// ParamKeyTable returns the TypeTable for coinswap module
func ParamKeyTable() KeyTable {
	return NewKeyTable().RegisterParamSet(&Params{})
}

func (p *Params) ParamSetPairs() ParamSetPairs {
	return ParamSetPairs{
		NewParamSetPair(KeyAssetParams, &p.AssetParams, validateAssetParams),
	}
}
