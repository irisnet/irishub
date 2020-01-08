package coinswap

import "math/big"

type GenesisState struct {
	Params Params `json:"params"`
}

type Params struct {
	Fee Rat `json:"fee"`
}

type Rat struct {
	*big.Rat `json:"rat"`
}
