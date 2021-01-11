package mint

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	Minter Minter `json:"minter"` // minter object
	Params Params `json:"params"` // inflation params
}

type Minter struct {
	LastUpdate    time.Time `json:"last_update"` // time which the last update was made to the minter
	MintDenom     string    `json:"mint_denom"`  // type of coin to mint
	InflationBase sdk.Int   `json:"inflation_basement"`
}

type Params struct {
	Inflation sdk.Dec `json:"inflation"` // inflation rate
}
