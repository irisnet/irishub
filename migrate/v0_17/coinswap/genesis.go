package coinswap

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	Params Params `json:"params"`
	Pools  []Pool `json:"pool"`
}

type Params struct {
	Fee Rat `json:"fee"`
}

type Rat struct {
	*big.Rat `json:"rat"`
}

type Pool struct {
	Coins sdk.Coins `json:"coins"`
	Name  string    `json:"name"`
}
