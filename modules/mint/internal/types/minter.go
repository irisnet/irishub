package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	blocksPerYear = 60 * 60 * 8766 / 5 // 5 second a block, 8766 = 365.25 * 24
)

var initialIssue = sdk.NewIntWithDecimal(20, 8)

// current inflation state
type Minter struct {
	LastUpdate    time.Time `json:"last_update"` // time which the last update was made to the minter
	InflationBase sdk.Int   `json:"inflation_basement"`
}

// Create a new minter object
func NewMinter(lastUpdate time.Time, inflationBase sdk.Int) Minter {
	return Minter{
		LastUpdate:    lastUpdate,
		InflationBase: inflationBase,
	}
}

// minter object for a new chain
func InitialMinter() Minter {
	return NewMinter(
		time.Unix(0, 0),
		initialIssue.Mul(sdk.NewIntWithDecimal(1, 18)), // 20*(10^8)iris, 20*(10^8)*(10^18)iris-atto
	)
}

func validateMinter(minter Minter) error {
	if minter.LastUpdate.Before(time.Unix(0, 0)) {
		return fmt.Errorf("minter last update time(%s) should not be a time before January 1, 1970 UTC", minter.LastUpdate.String())
	}
	if !minter.InflationBase.GT(sdk.ZeroInt()) {
		return fmt.Errorf("minter inflation basement (%s) should be positive", minter.InflationBase.String())
	}
	return nil
}

// get the provisions for a block based on the annual provisions rate
func (m Minter) NextAnnualProvisions(params Params) (provisions sdk.Dec) {
	return params.Inflation.MulInt(m.InflationBase)
}

// get the provisions for a block based on the annual provisions rate
func (m Minter) BlockProvision(params Params) sdk.Coin {
	provisions := m.NextAnnualProvisions(params)
	blockInflationAmount := provisions.QuoInt(sdk.NewInt(blocksPerYear))
	return sdk.NewCoin(params.MintDenom, blockInflationAmount.TruncateInt())
}
