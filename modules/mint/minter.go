package mint

import (
	"fmt"
	stakeTypes "github.com/irisnet/irishub/modules/stake/types"
	sdk "github.com/irisnet/irishub/types"
	"time"
)

const (
	blocksPerYear     = 60 * 60 * 8766 / 5 // 5 second a block, 8766 = 365.25 * 24
)

// current inflation state
type Minter struct {
	LastUpdate    time.Time `json:"last_update"`       // time which the last update was made to the minter
	MintDenom     string    `json:"mint_denom"`        // type of coin to mint
	InflationBase sdk.Int   `json:"inflation_basement"`
}

// Create a new minter object
func NewMinter(lastUpdate time.Time, mintDenom string, inflationBase sdk.Int) Minter {
	return Minter{
		LastUpdate:    lastUpdate,
		MintDenom:     mintDenom,
		InflationBase: inflationBase,
	}
}

// minter object for a new chain
func InitialMinter() Minter {
	return NewMinter(
		time.Unix(0, 0),
		stakeTypes.StakeDenom,
		sdk.InitialIssue.Mul(sdk.NewIntWithDecimal(1, 18)), // 2*(10^9)iris, 2*(10^9)*(10^18)iris-atto
	)
}

func validateMinter(minter Minter) error {
	if minter.LastUpdate.Before(time.Unix(0, 0)) {
		return fmt.Errorf("mint last update time(%s) should not be a time before January 1, 1970 UTC", minter.LastUpdate.String())
	}
	if len(minter.MintDenom) == 0 {
		return fmt.Errorf("mint token denom (%s) should not be empty", minter.MintDenom)
	}
	if !minter.InflationBase.GT(sdk.ZeroInt()) {
		return fmt.Errorf("mint inflation basement (%s) should be positive", minter.InflationBase.String())
	}
	return nil
}

// get the provisions for a block based on the annual provisions rate
func (m Minter) NextAnnualProvisions(params Params) (provisions sdk.Dec) {
	return params.Inflation.MulInt(m.InflationBase)
}

// get the provisions for a block based on the annual provisions rate
func (m Minter) BlockProvision(params Params, annualProvisions sdk.Dec) sdk.Coin {
	blockInflationAmount := annualProvisions.QuoInt(sdk.NewInt(blocksPerYear))
	return sdk.NewCoin(m.MintDenom, blockInflationAmount.TruncateInt())
}
