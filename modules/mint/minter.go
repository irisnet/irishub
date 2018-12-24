package mint

import (
	"fmt"
	"time"

	stakeTypes "github.com/irisnet/irishub/modules/stake/types"
	sdk "github.com/irisnet/irishub/types"
)

const (
	nanoToMiliSecond  = 1000000
	miliSecondPerYear = 60 * 60 * 8766 * 1000
)

// current inflation state
type Minter struct {
	LastUpdate        time.Time `json:"last_update"`       // time which the last update was made to the minter
	MintDenom         string    `json:"mint_denom"`        // type of coin to mint
	InflationBasement sdk.Int   `json:"inflation_basement"`
}

// Create a new minter object
func NewMinter(lastUpdate time.Time, mintDenom string, inflationBasement sdk.Int) Minter {
	return Minter{
		LastUpdate:        lastUpdate,
		MintDenom:         mintDenom,
		InflationBasement: inflationBasement,
	}
}

// minter object for a new chain
func InitialMinter() Minter {
	return NewMinter(
		time.Unix(0, 0),
		stakeTypes.StakeDenom,
		sdk.NewIntWithDecimal(2, 9).Mul(sdk.NewIntWithDecimal(1, 18)), // 2*(10^9)iris, 2*(10^9)*(10^18)iris-atto
	)
}

func validateMinter(minter Minter) error {
	if minter.LastUpdate.Nanosecond() < 0 {
		return fmt.Errorf("mint parameter Inflation should be positive, is %s ", minter.LastUpdate.String())
	}
	return nil
}

// get the provisions for a block based on the annual provisions rate
func (m Minter) NextAnnualProvisions(params Params) (provisions sdk.Dec) {
	return params.Inflation.MulInt(m.InflationBasement)
}

// get the provisions for a block based on the annual provisions rate
func (m Minter) BlockProvision(params Params, annualProvisions sdk.Dec, inflationTime time.Time) sdk.Coin {
	inflationPeriod := inflationTime.Sub(m.LastUpdate)
	millisecond := inflationPeriod.Nanoseconds() / int64(nanoToMiliSecond)
	blockInflationAmount := annualProvisions.Mul(sdk.NewDec(millisecond)).Quo(sdk.NewDec(int64(miliSecondPerYear)))
	return sdk.NewCoin(m.MintDenom, blockInflationAmount.TruncateInt())
}
