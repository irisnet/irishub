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

// Minter difines current inflation state
type Minter struct {
	LastUpdate    time.Time `json:"last_update"`        // time which the last update was made to the minter
	InflationBase sdk.Int   `json:"inflation_basement"` // base inflation
}

// Create a new minter object
func NewMinter(lastUpdate time.Time, inflationBase sdk.Int) Minter {
	return Minter{
		LastUpdate:    lastUpdate,
		InflationBase: inflationBase,
	}
}

// DefaultMinter returns minter object for a new chain
func DefaultMinter() Minter {
	return NewMinter(
		time.Unix(0, 0).UTC(),
		initialIssue.Mul(sdk.NewIntWithDecimal(1, 18)), // 20*(10^8)iris, 20*(10^8)*(10^18)iris-atto
	)
}

// ValidateMinter returns err if the Minter is invalid
func ValidateMinter(m Minter) error {
	if m.LastUpdate.Before(time.Unix(0, 0)) {
		return fmt.Errorf("minter last update time(%s) should not be a time before January 1, 1970 UTC", m.LastUpdate.String())
	}
	if !m.InflationBase.GT(sdk.ZeroInt()) {
		return fmt.Errorf("minter inflation basement (%s) should be positive", m.InflationBase.String())
	}
	return nil
}

// NextAnnualProvisions gets the provisions for a block based on the annual provisions rate
func (m Minter) NextAnnualProvisions(params Params) (provisions sdk.Dec) {
	return params.Inflation.MulInt(m.InflationBase)
}

// BlockProvision gets the provisions for a block based on the annual provisions rate
func (m Minter) BlockProvision(params Params) sdk.Coin {
	provisions := m.NextAnnualProvisions(params)
	blockInflationAmount := provisions.QuoInt(sdk.NewInt(blocksPerYear))
	return sdk.NewCoin(params.MintDenom, blockInflationAmount.TruncateInt())
}
