package mint

import (
	"fmt"
	"time"

	sdk "github.com/irisnet/irishub/types"
)

var (
	nanoToMiliSecond = 1000000
	miliSecondPerYear = 60 * 60 * 8766 * 1000
)

// current inflation state
type Minter struct {
	LastUpdate       time.Time `json:"last_update"`       // time which the last update was made to the minter
	AnnualProvisions sdk.Dec   `json:"annual_provisions"` // current annual expected provisions
}

// Create a new minter object
func NewMinter(lastUpdate time.Time, annualProvisions sdk.Dec) Minter {
	return Minter{
		LastUpdate:       lastUpdate,
		AnnualProvisions: annualProvisions,
	}
}

// minter object for a new chain
func InitialMinter() Minter {
	return NewMinter(
		time.Unix(0, 0),
		sdk.NewDec(0),
	)
}

func validateMinter(minter Minter) error {
	if minter.LastUpdate.Nanosecond() < 0 {
		return fmt.Errorf("mint parameter Inflation should be positive, is %s ", minter.LastUpdate.String())
	}
	if minter.AnnualProvisions.LT(sdk.ZeroDec()) {
		return fmt.Errorf("mint annual provisions should be positive, is %s ", minter.AnnualProvisions.String())
	}
	return nil
}

// get the provisions for a block based on the annual provisions rate
func (m Minter) NextAnnualProvisions(params Params) (provisions sdk.Dec) {
	return params.Inflation.MulInt(params.InflationBasement)
}

// get the provisions for a block based on the annual provisions rate
func (m Minter) BlockProvision(params Params, inflationTime time.Time) sdk.Coin {
	inflationPeriod := inflationTime.Sub(m.LastUpdate)
	millisecond := inflationPeriod.Nanoseconds()/int64(nanoToMiliSecond)
	blockInflationAmount := m.AnnualProvisions.Mul(sdk.NewDec(millisecond)).Quo(sdk.NewDec(int64(miliSecondPerYear)))
	return sdk.NewCoin(params.MintDenom, blockInflationAmount.TruncateInt())
}
